package dao

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var voteProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(VoteProcessor)
	},
}

func (Vote) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type VoteProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewVoteProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new VoteProcessor")

		nopp := voteProcessorPool.Get()
		opp, ok := nopp.(*VoteProcessor)
		if !ok {
			return nil, errors.Errorf("expected VoteProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b
		opp.getLastBlockFunc = getLastBlockFunc

		return opp, nil
	}
}

func (opp *VoteProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Vote")

	fact, ok := op.Fact().(VoteFact)
	if !ok {
		return ctx, nil, e.Errorf("not VoteFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot vote, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao contract account not found, %s: %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("fee currency doesn't exist, %q: %w", fact.Currency(), err), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao design not found, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal state not found, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s, %q, %q", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	if p.Status() != types.PreSnapped {
		return nil, base.NewBaseOperationProcessReasonError("proposal not in pre-snapped status, %s, %q, %q, %q", fact.Contract(), fact.DAOID(), fact.ProposalID(), p.Status()), nil
	}

	blockMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(p.Policy(), p.Proposal(), types.Voting, blockMap)
	if period != types.Voting {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within Voting period, Voting period; start(%d), end(%d), but now(%d)", start, end, blockMap.Manifest().ProposedAt().Unix()), nil
	}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	default:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		for i, v := range voters {
			if v.Account().Equal(fact.Sender()) {
				break
			}

			if i == len(voters)-1 {
				return nil, base.NewBaseOperationProcessReasonError("sender is not registered as voter, sender(%s), %s, %q, %q", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
			}
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		vpb, err := state.StateVotingPowerBoxValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		vp, found := vpb.VotingPowers()[fact.Sender().String()]
		if !found {
			return nil, base.NewBaseOperationProcessReasonError("sender voting power not found, sender(%s), %s, %q, %q", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
		}

		if vp.Voted() {
			return nil, base.NewBaseOperationProcessReasonError("sender already voted, sender(%s), %s, %q, %q", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *VoteProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Vote")

	fact, ok := op.Fact().(VoteFact)
	if !ok {
		return nil, nil, e.Errorf("expected VoteFact, not %T", op.Fact())
	}

	sts := []base.StateMergeValue{}

	{ // caculate operation fee
		policy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find currency policy, %q: %w", fact.Currency(), err), nil
		}

		fee, err := policy.Feeer().Fee(common.ZeroBig)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
		}

		st, err := currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %s, %q: %w", fact.Sender(), fact.Currency(), err), nil
		}
		sb := currencystate.NewStateMergeValue(st.Key(), st.Value())

		switch b, err := currency.StateBalanceValue(st); {
		case err != nil:
			return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %s, %q: %w", fact.Sender(), fact.Currency(), err), nil
		case b.Big().Compare(fee) < 0:
			return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %s, %q", fact.Sender(), fact.Currency()), nil
		}

		v, ok := sb.Value().(currency.BalanceStateValue)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
		}
		sts = append(sts, currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))))
	}

	var votingPowerBox types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	default:
		vpb, err := state.StateVotingPowerBoxValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}
		votingPowerBox = vpb
	}

	vp, found := votingPowerBox.VotingPowers()[fact.Sender().String()]
	if !found {
		return nil, base.NewBaseOperationProcessReasonError("sender voting power not found, sender(%s), %s, %q, %q", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}
	vp.SetVoted(true)
	vp.SetVoteFor(fact.Vote())

	vpb := votingPowerBox.VotingPowers()
	vpb[fact.Sender().String()] = vp
	votingPowerBox.SetVotingPowers(vpb)

	result := votingPowerBox.Result()
	if _, found := result[fact.Vote()]; found {
		result[fact.Vote()] = result[fact.Vote()].Add(vp.Amount())
	} else {
		result[fact.Vote()] = common.ZeroBig.Add(vp.Amount())
	}
	votingPowerBox.SetResult(result)

	sts = append(sts,
		currencystate.NewStateMergeValue(
			state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()),
			state.NewVotingPowerBoxStateValue(votingPowerBox),
		),
	)

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err := currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %s, %q: %w", fact.Sender(), fact.Currency(), err), nil
	}
	sb := currencystate.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %s, %q: %w", fact.Sender(), fact.Currency(), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %s, %q", fact.Sender(), fact.Currency()), nil
	}

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}

	sts = append(sts,
		currencystate.NewStateMergeValue(
			sb.Key(),
			currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
		),
	)

	return sts, nil, nil
}

func (opp *VoteProcessor) Close() error {
	voteProcessorPool.Put(opp)

	return nil
}
