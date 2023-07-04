package dao

import (
	"context"
	"sync"
	"time"

	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"

	"github.com/ProtoconNet/mitum-currency/v3/common"
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
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot vote, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao contract account not found, %q: %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("fee currency doesn't exist, %q: %w", fact.Currency(), err), nil
	}

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao design state not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao value not found from state, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal state not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	if p.Status() != types.PreSnapped {
		return nil, base.NewBaseOperationProcessReasonError("proposal not in pre-snapped status, %s-%s-%s, %q", fact.Contract(), fact.DAOID(), fact.ProposalID(), p.Status()), nil
	}

	blockMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(design.Policy(), p.Proposal(), blockMap)

	if period != types.Voting {
		return nil, base.NewBaseOperationProcessReasonError("not voting period, vote period start : %s, end : %s", time.Unix(start, 0), time.Unix(end, 0)), nil
	}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("sender is not registered as voter, %s, %s-%s-%s: %w", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	default:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		for i, voter := range voters {
			if voter.Account().Equal(fact.Sender()) {
				break
			}

			if i == len(voters)-1 {
				return nil, base.NewBaseOperationProcessReasonError("sender is not registered as voter, %s, %s-%s-%s", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
			}
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		vpb, err := state.StateVotingPowerBoxValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		powers := vpb.VotingPowers()
		if _, found := powers[fact.Sender()]; found {
			return nil, base.NewBaseOperationProcessReasonError("sender already voted, %s, %s-%s-%s", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
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

	sts := make([]base.StateMergeValue, 2)

	var delegators []base.Address
	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voter state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		vs, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voter value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		for i, vt := range vs {
			if vt.Account().Equal(fact.Sender()) {
				delegators = vt.Delegators()
				break
			}

			if i == len(vs)-1 {
				return nil, base.NewBaseOperationProcessReasonError("sender is not voter, %s, %s-%s-%s: %w", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
			}
		}
	default:
		return nil, base.NewBaseOperationProcessReasonError("sender is not voter, %s, %s-%s-%s: %w", fact.Sender(), fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	var votingPowerBox types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		vpb, err := state.StateVotingPowerBoxValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}
		votingPowerBox = vpb
	default:
		votingPowerBox = types.NewVotingPowerBox(common.ZeroBig, map[base.Address]common.Big{})
	}

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao design state not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao value not found from state, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	token := design.Policy().Token()

	addTotal := common.ZeroBig
	for _, delegator := range delegators {
		st, err = currencystate.ExistsState(currency.StateKeyBalance(delegator, token), "key of delegator balance", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("delegator balance not found, %q: %w", delegator, err), nil
		}

		b, err := currency.StateBalanceValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", currency.StateKeyBalance(delegator, token), err), nil
		}

		addTotal = addTotal.Add(b.Big())
	}

	total := votingPowerBox.Total().Add(addTotal)

	powers := votingPowerBox.VotingPowers()
	powers[fact.Sender()] = addTotal

	result := votingPowerBox.Result()
	if _, found := result[fact.Vote()]; found {
		result[fact.Vote()] = result[fact.Vote()].Add(addTotal)
	}

	votingPowerBox.SetTotal(total)
	votingPowerBox.SetVotingPowers(powers)
	votingPowerBox.SetResult(result)

	sts[0] = currencystate.NewStateMergeValue(
		state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()),
		state.NewVotingPowerBoxStateValue(votingPowerBox),
	)

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err = currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}
	sb := currencystate.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := currency.StateBalanceValue(st); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", currency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(currency.BalanceStateValue)
	if !ok {
		return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[1] = currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))))

	return sts, nil, nil
}

func (opp *VoteProcessor) Close() error {
	voteProcessorPool.Put(opp)

	return nil
}
