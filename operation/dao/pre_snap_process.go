package dao

import (
	"context"
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
	"sync"
)

var preSnapProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(PreSnapProcessor)
	},
}

func (PreSnap) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type PreSnapProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc types.GetLastBlockFunc
}

func NewPreSnapProcessor(getLastBlockFunc types.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new PreSnapProcessor")

		nopp := preSnapProcessorPool.Get()
		opp, ok := nopp.(*PreSnapProcessor)
		if !ok {
			return nil, errors.Errorf("expected PreSnapProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e(err, "")
		}

		opp.BaseOperationProcessor = b
		opp.getLastBlockFunc = getLastBlockFunc

		return opp, nil
	}
}

func (opp *PreSnapProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess PreSnap")

	fact, ok := op.Fact().(PreSnapFact)
	if !ok {
		return ctx, nil, e(nil, "not PreSnapFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot preSnap, %q: %w", fact.Sender(), err), nil
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
		return nil, base.NewBaseOperationProcessReasonError("dao design value not found from state, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	} else if p.Status() == types.PreSnapped {
		return nil, base.NewBaseOperationProcessReasonError("already preSnapped, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	blocMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(design.Policy(), p.Proposal(), blocMap)
	if period != types.PreSnapshot {
		return nil, base.NewBaseOperationProcessReasonError("currency time is not within the PreSnapshotPeriod, PreSnapshotPeriod start : %s, end %s", start, end), nil
	}

	_, err = currencystate.ExistsState(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of voters", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voters state not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	_, err = currencystate.ExistsState(state.StateKeyDelegators(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of delegators", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("delegators state not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	_, err = currencystate.NotExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of voting power box", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voting power box state already created, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *PreSnapProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process PreSnap")

	fact, ok := op.Fact().(PreSnapFact)
	if !ok {
		return nil, nil, e(nil, "expected PreSnapFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 1)

	//st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("dao not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	//}
	//
	//design, err := state.StateDesignValue(st)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("dao design value not found from state, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	//}
	//
	//st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//}
	//
	//p, err := state.StateProposalValue(st)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//}

	//var votingPowerBox types.VotingPowerBox
	//switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	//case err != nil:
	//	return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//case found:
	//	if vb, err := state.StateVotingPowerBoxValue(st); err != nil {
	//		return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//	} else {
	//		votingPowerBox = vb
	//	}
	//default:
	//	votingPowerBox = types.NewVotingPowerBox(common.ZeroBig, []types.VotingPower{})
	//}

	//st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//}
	//
	//votingPowerToken := design.Policy().Token()
	//
	//switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	//case err != nil:
	//	return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//case found:
	//	voters, err := state.StateVotersValue(st)
	//	if err != nil {
	//		return nil, base.NewBaseOperationProcessReasonError("failed to find voters value, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//	}
	//
	//	total := common.ZeroBig
	//	for _, info := range voters {
	//		votingPower := common.ZeroBig
	//		for _, delegator := range info.Delegators() {
	//			st, err = currencystate.ExistsState(currency.StateKeyBalance(delegator, votingPowerToken), "key of balance", getStateFunc)
	//			if err != nil {
	//				continue
	//			}
	//
	//			b, err := currency.StateBalanceValue(st)
	//			if err != nil {
	//				return nil, base.NewBaseOperationProcessReasonError("failed to find balance value of the delegator from state, %q, %q: %w", delegator, votingPowerToken, err), nil
	//			}
	//
	//			votingPower = votingPower.Add(b.Big())
	//		}
	//		votingPowerBox[info.Account()] {
	//
	//		}
	//		votingPowers = append(votingPowers, types.NewVotingPower(info.Account(), total))
	//	}
	//}
	//
	//snaps = append(snaps, types.NewSnapHistory(uint64(time.Now().UnixMilli()), votingPowers))
	//
	//sts[0] = currencystate.NewStateMergeValue(
	//	state.StateKeySnapHistories(fact.Contract(), fact.DAOID(), fact.ProposalID()),
	//	state.NewSnapHistoriesStateValue(snaps),
	//)
	//
	//st, err = currencystate.ExistsState(currency.StateKeyCurrencyDesign(votingPowerToken), "key of currency design", getStateFunc)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("currency design not found, %q: %w", votingPowerToken, err), nil
	//}
	//
	//vpDesign, err := currency.StateCurrencyDesignValue(st)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("currency design value not found, %q: %w", votingPowerToken, err), nil
	//}
	//
	//aggregate := vpDesign.Aggregate()
	//turnout := design.Policy().Turnout()
	//tq := turnout.Quorum(aggregate)
	//
	//if blocktime < votingstart {
	//	registerTotal := common.ZeroBig
	//	for _, vp := range votingPowers {
	//		registerTotal = registerTotal.Add(vp.Amount())
	//	}
	//
	//	if registerTotal.Compare(tq) < 0 {
	//		sts = append(sts,
	//			currencystate.NewStateMergeValue(
	//				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
	//				state.NewProposalStateValue(false, proposal),
	//			),
	//		)
	//
	//		sts = append(sts,
	//			currencystate.NewStateMergeValue(
	//				state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()),
	//				state.NewVotingPowerBoxStateValue(false, 0, nil),
	//			),
	//		)
	//	} else {
	//		vps := make([]types.VotingPowerBox, proposal.Options())
	//		for i := range vps {
	//			vps[i] = types.NewVotingPowerBox(common.ZeroBig, []types.VotingPower{})
	//		}
	//
	//		sts = append(sts,
	//			currencystate.NewStateMergeValue(
	//				state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()),
	//				state.NewVotingPowerBoxStateValue(true, 0, vps),
	//			),
	//		)
	//	}
	//} else if votingend <= blocktime {
	//	st, err = currencystate.ExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of votes", getStateFunc)
	//	if err != nil {
	//		return nil, base.NewBaseOperationProcessReasonError("votes not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//	}
	//
	//	votes, err := state.StateVotingPowerBoxValue(st)
	//	if err != nil {
	//		return nil, base.NewBaseOperationProcessReasonError("votes value not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	//	}
	//
	//	votingTotal := common.ZeroBig
	//	for _, vps := range votes.votingPowerBox() {
	//		votingTotal = votingTotal.Add(vps.Total())
	//	}
	//
	//	if votingTotal.Compare(tq) < 0 {
	//		sts = append(sts,
	//			currencystate.NewStateMergeValue(
	//				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
	//				state.NewProposalStateValue(false, proposal),
	//			),
	//		)
	//
	//		sts = append(sts,
	//			currencystate.NewStateMergeValue(
	//				state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()),
	//				state.NewVotingPowerBoxStateValue(false, 0, votes.votingPowerBox()),
	//			),
	//		)
	//	}
	//}

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
	sts = append(sts, currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))))

	return sts, nil, nil
}

func (opp *PreSnapProcessor) Close() error {
	preSnapProcessorPool.Put(opp)

	return nil
}