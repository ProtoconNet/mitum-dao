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

var cancelProposalProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CancelProposalProcessor)
	},
}

func (CancelProposal) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CancelProposalProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewCancelProposalProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new CancelProposalProcessor")

		nopp := cancelProposalProcessorPool.Get()
		opp, ok := nopp.(*CancelProposalProcessor)
		if !ok {
			return nil, errors.Errorf("expected CancelProposalProcessor, not %T", nopp)
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

func (opp *CancelProposalProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess CancelProposal")

	fact, ok := op.Fact().(CancelProposalFact)
	if !ok {
		return ctx, nil, e.Errorf("not CancelProposalFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot cancel-proposal, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao contract account not found, %s: %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("fee currency doesn't exist, %q: %w", fact.Currency(), err), nil
	}

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao design state not found, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao design value not found from state, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if !fact.Sender().Equal(p.Proposal().Proposer()) {
		return nil, base.NewBaseOperationProcessReasonError("sender is not proposer of the proposal, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s, %q, %q", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	} else if p.Status() != types.Proposed {
		return nil, base.NewBaseOperationProcessReasonError("cancel-proposal is unavailable, %s, %q, %q", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	blockMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, _, _ := types.GetPeriodOfCurrentTime(design.Policy(), p.Proposal(), blockMap)
	if !(period == types.PreLifeCycle || period == types.ProposalReview || period == types.Registration) {
		start, _ := types.GetPeriodTime(types.Voting, design.Policy(), p.Proposal())
		return nil, base.NewBaseOperationProcessReasonError("cancellable period has passed; voting-started(%d), now(%d)", start, blockMap.Manifest().ProposedAt().Unix()), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *CancelProposalProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process CancelProposal")

	fact, ok := op.Fact().(CancelProposalFact)
	if !ok {
		return nil, nil, e.Errorf("expected CancelProposalFact, not %T", op.Fact())
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

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
		state.NewProposalStateValue(types.Canceled, p.Proposal()),
	))

	return sts, nil, nil
}

func (opp *CancelProposalProcessor) Close() error {
	cancelProposalProcessorPool.Put(opp)

	return nil
}
