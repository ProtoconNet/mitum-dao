package dao

import (
	"context"
	"sync"
	"time"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var snapProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SnapProcessor)
	},
}

func (Snap) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type SnapProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc types.GetLastBlockFunc
}

func NewSnapProcessor(getLastBlockFunc types.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new SnapProcessor")

		nopp := snapProcessorPool.Get()
		opp, ok := nopp.(*SnapProcessor)
		if !ok {
			return nil, errors.Errorf("expected SnapProcessor, not %T", nopp)
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

func (opp *SnapProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess Snap")

	fact, ok := op.Fact().(SnapFact)
	if !ok {
		return ctx, nil, e(nil, "not SnapFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot snap, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account not found, %q: %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency doesn't exist, %q: %w", fact.Currency(), err), nil
	}

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao value not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	delaytime := design.Policy().DelayTime()
	registerperiod := design.Policy().RegsiterPeriod()
	snaptime := design.Policy().SnapTime()
	voteperiod := design.Policy().VotePeriod()

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposeID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	}

	if !p.Active {
		return nil, base.NewBaseOperationProcessReasonError("already closed proposal, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposeID()), nil
	}
	proposal := p.Proposal

	starttime := proposal.StartTime()

	blockmap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	blocktime := uint64(blockmap.Manifest().ProposedAt().Unix())

	if blocktime < starttime+delaytime+registerperiod {
		return nil, base.NewBaseOperationProcessReasonError("registration is still in progress, must in %d <= block(%d)", starttime+delaytime, blocktime), nil
	}

	votingstart := starttime + delaytime + registerperiod + snaptime
	votingend := starttime + delaytime + registerperiod + snaptime + voteperiod

	if votingstart <= blocktime && blocktime < votingend {
		return nil, base.NewBaseOperationProcessReasonError("voting is still in progress, now voting start <= block(%d) < voting end", blocktime), nil
	}

	switch st, found, err := getStateFunc(state.StateKeySnapHistories(fact.Contract(), fact.DAOID(), fact.ProposeID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find snap histories, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	case found:
		snaps, err := state.StateSnapHistoriesValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find snap histories value, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
		}

		if 0 < len(snaps) {
			lastSnapped := snaps[len(snaps)-1].TimeStamp()

			if (lastSnapped < votingstart && blocktime < votingstart) || votingend <= lastSnapped {
				return nil, base.NewBaseOperationProcessReasonError("already snapped proposal, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
			}
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *SnapProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process Snap")

	fact, ok := op.Fact().(SnapFact)
	if !ok {
		return nil, nil, e(nil, "expected SnapFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 1)

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao value not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	delaytime := design.Policy().DelayTime()
	registerperiod := design.Policy().RegsiterPeriod()
	snaptime := design.Policy().SnapTime()
	voteperiod := design.Policy().VotePeriod()

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposeID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	}

	if !p.Active {
		return nil, base.NewBaseOperationProcessReasonError("already closed proposal, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposeID()), nil
	}
	proposal := p.Proposal

	starttime := proposal.StartTime()

	blockmap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	blocktime := uint64(blockmap.Manifest().ProposedAt().Unix())

	votingstart := starttime + delaytime + registerperiod + snaptime
	votingend := starttime + delaytime + registerperiod + snaptime + voteperiod

	var snaps []types.SnapHistory

	switch st, found, err := getStateFunc(state.StateKeySnapHistories(fact.Contract(), fact.DAOID(), fact.ProposeID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find snap histories, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	case found:
		sn, err := state.StateSnapHistoriesValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find snap histories value, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
		}
		snaps = sn
	default:
		snaps = []types.SnapHistory{}
	}

	votingPowers := []types.VotingPower{}
	votingPowerToken := design.Policy().Token()

	switch st, found, err := getStateFunc(state.StateKeyRegisterList(fact.Contract(), fact.DAOID(), fact.ProposeID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find register list, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	case found:
		registers, err := state.StateRegisterListValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find register list value, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
		}

		for _, info := range registers {
			total := common.ZeroBig
			for _, approved := range info.ApprovedBy() {
				st, err = currencystate.ExistsState(currency.StateKeyBalance(approved, votingPowerToken), "key of balance", getStateFunc)
				if err != nil {
					continue
				}

				b, err := currency.StateBalanceValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to find balance value, %q, %q: %w", approved, votingPowerToken, err), nil
				}

				total = total.Add(b.Big())
			}
			votingPowers = append(votingPowers, types.NewVotingPower(info.Account(), total))
		}
	}

	snaps = append(snaps, types.NewSnapHistory(uint64(time.Now().UnixMilli()), votingPowers))

	sts[0] = currencystate.NewStateMergeValue(
		state.StateKeySnapHistories(fact.Contract(), fact.DAOID(), fact.ProposeID()),
		state.NewSnapHistoriesStateValue(snaps),
	)

	st, err = currencystate.ExistsState(currency.StateKeyCurrencyDesign(votingPowerToken), "key of currency design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency design not found, %q: %w", votingPowerToken, err), nil
	}

	vpDesign, err := currency.StateCurrencyDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency design value not found, %q: %w", votingPowerToken, err), nil
	}

	aggregate := vpDesign.Aggregate()
	turnout := design.Policy().Turnout()
	tq := turnout.Quorum(aggregate)

	if blocktime < votingstart {
		registerTotal := common.ZeroBig
		for _, vp := range votingPowers {
			registerTotal = registerTotal.Add(vp.VotingPower())
		}

		if registerTotal.Compare(tq) < 0 {
			sts = append(sts,
				currencystate.NewStateMergeValue(
					state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposeID()),
					state.NewProposalStateValue(false, proposal),
				),
			)

			sts = append(sts,
				currencystate.NewStateMergeValue(
					state.StateKeyVotes(fact.Contract(), fact.DAOID(), fact.ProposeID()),
					state.NewVotesStateValue(false, 0, nil),
				),
			)
		} else {
			vps := make([]types.VotingPowers, proposal.Options())
			for i := range vps {
				vps[i] = types.NewVotingPowers(common.ZeroBig, []types.VotingPower{})
			}

			sts = append(sts,
				currencystate.NewStateMergeValue(
					state.StateKeyVotes(fact.Contract(), fact.DAOID(), fact.ProposeID()),
					state.NewVotesStateValue(true, 0, vps),
				),
			)
		}
	} else if votingend <= blocktime {
		st, err = currencystate.ExistsState(state.StateKeyVotes(fact.Contract(), fact.DAOID(), fact.ProposeID()), "key of votes", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("votes not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
		}

		votes, err := state.StateVotesValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("votes value not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
		}

		votingTotal := common.ZeroBig
		for _, vps := range votes.Votes {
			votingTotal = votingTotal.Add(vps.Total())
		}

		if votingTotal.Compare(tq) < 0 {
			sts = append(sts,
				currencystate.NewStateMergeValue(
					state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposeID()),
					state.NewProposalStateValue(false, proposal),
				),
			)

			sts = append(sts,
				currencystate.NewStateMergeValue(
					state.StateKeyVotes(fact.Contract(), fact.DAOID(), fact.ProposeID()),
					state.NewVotesStateValue(false, 0, votes.Votes),
				),
			)
		}
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
	sts = append(sts, currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))))

	return sts, nil, nil
}

func (opp *SnapProcessor) Close() error {
	snapProcessorPool.Put(opp)

	return nil
}
