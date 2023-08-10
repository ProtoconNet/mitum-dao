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

var postSnapProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(PostSnapProcessor)
	},
}

func (PostSnap) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type PostSnapProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewPostSnapProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new PostSnapProcessor")

		nopp := postSnapProcessorPool.Get()
		opp, ok := nopp.(*PostSnapProcessor)
		if !ok {
			return nil, errors.Errorf("expected PostSnapProcessor, not %T", nopp)
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

func (opp *PostSnapProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess PostSnap")

	fact, ok := op.Fact().(PostSnapFact)
	if !ok {
		return ctx, nil, e.Errorf("not PostSnapFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot post-snap, %s: %w", fact.Sender(), err), nil
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
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s, %q, %q", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	} else if p.Status() == types.PostSnapped {
		return nil, base.NewBaseOperationProcessReasonError("already post snapped, %s, %q, %q", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	blockMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(p.Policy(), p.Proposal(), types.PostSnapshot, blockMap)
	if period != types.PostSnapshot {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within the PostSnapshotPeriod, PostSnapshotPeriod; start(%d), end(%d), but now(%d)", start, end, blockMap.Manifest().ProposedAt().Unix()), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voters state not found, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *PostSnapProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process PostSnap")

	fact, ok := op.Fact().(PostSnapFact)
	if !ok {
		return nil, nil, e.Errorf("expected PostSnapFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue

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

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() != types.PreSnapped {
		sts = append(sts,
			currencystate.NewStateMergeValue(
				st.Key(),
				state.NewProposalStateValue(types.Canceled, p.Proposal(), p.Policy()),
			),
		)

		return sts, nil, nil
	}

	var ovpb types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		if vb, err := state.StateVotingPowerBoxValue(st); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		} else {
			ovpb = vb
		}
	default:
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s, %q, %q", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	votingPowerToken := p.Policy().Token()

	var nvpb = types.NewVotingPowerBox(common.ZeroBig, map[string]types.VotingPower{})

	nvps := map[string]types.VotingPower{}
	nvt := common.ZeroBig

	votedTotal := common.ZeroBig
	votingResult := map[uint8]common.Big{}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		for _, info := range voters {
			a := info.Account().String()

			if !ovpb.VotingPowers()[a].Voted() {
				nvps[a] = ovpb.VotingPowers()[a]
				continue
			}

			vp := common.ZeroBig
			for _, delegator := range info.Delegators() {
				st, err = currencystate.ExistsState(currency.StateKeyBalance(delegator, votingPowerToken), "key of balance", getStateFunc)
				if err != nil {
					continue
				}

				b, err := currency.StateBalanceValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to find balance value of the delegator from state, %s, %q: %w", delegator, votingPowerToken, err), nil
				}

				vp = vp.Add(b.Big())
			}

			ovp := ovpb.VotingPowers()[a]
			if ovp.Amount().Compare(vp) < 0 {
				nvps[a] = ovp
			} else {
				nvp := types.NewVotingPower(info.Account(), vp)
				nvp.SetVoted(ovp.Voted())
				nvp.SetVoteFor(ovp.VoteFor())

				nvps[a] = nvp
			}

			nvt = nvt.Add(nvps[a].Amount())

			if nvps[a].Voted() {
				if _, found := votingResult[nvps[a].VoteFor()]; !found {
					votingResult[nvps[a].VoteFor()] = common.ZeroBig
				}
				votingResult[nvps[a].VoteFor()] = votingResult[nvps[a].VoteFor()].Add(vp)
				votedTotal = votedTotal.Add(nvps[a].Amount())
			}
		}

		nvpb.SetVotingPowers(nvps)
		nvpb.SetTotal(nvt)
		nvpb.SetResult(votingResult)
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()),
		state.NewVotingPowerBoxStateValue(nvpb),
	))

	st, err = currencystate.ExistsState(currency.StateKeyCurrencyDesign(votingPowerToken), "key of currency design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency state, %q: %w", votingPowerToken, err), nil
	}

	currencyDesign, err := currency.StateCurrencyDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency design value from state, %q: %w", votingPowerToken, err), nil
	}

	actualTurnoutCount := p.Policy().Turnout().Quorum(currencyDesign.Aggregate())
	actualQuorumCount := p.Policy().Quorum().Quorum(votedTotal)

	r := types.Completed

	if nvpb.Total().Compare(actualTurnoutCount) < 0 {
		r = types.Canceled
	} else if votedTotal.Compare(actualQuorumCount) < 0 {
		r = types.Rejected
	} else if p.Proposal().Type() == types.ProposalCrypto {
		vr0, found0 := votingResult[0]
		vr1, found1 := votingResult[1]

		if !(found0 && 0 < vr0.Compare(actualQuorumCount) && (!found1 || (found1 && 0 < vr0.Compare(vr1)))) {
			r = types.Rejected
		}
	} else if p.Proposal().Type() == types.ProposalBiz {
		options := p.Proposal().Options() - 1

		var count = 0
		var mvp = common.ZeroBig
		var i uint8 = 0

		for ; i < options; i++ {
			if votingResult[i].Compare(actualQuorumCount) >= 0 {
				if mvp.Compare(votingResult[i]) < 0 {
					count = 1
					mvp = votingResult[i]
				} else if mvp.Equal(votingResult[i]) {
					count += 1
				}
			}
		}

		if count != 1 {
			r = types.Rejected
		}
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
		state.NewProposalStateValue(r, p.Proposal(), p.Policy()),
	))

	return sts, nil, nil
}

func (opp *PostSnapProcessor) Close() error {
	postSnapProcessorPool.Put(opp)

	return nil
}
