package dao

import (
	"context"
	"sync"

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

var registerProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RegisterProcessor)
	},
}

func (Register) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RegisterProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc types.GetLastBlockFunc
}

func NewRegisterProcessor(getLastBlockFunc types.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringErrorFunc("failed to create new RegisterProcessor")

		nopp := registerProcessorPool.Get()
		opp, ok := nopp.(*RegisterProcessor)
		if !ok {
			return nil, errors.Errorf("expected RegisterProcessor, not %T", nopp)
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

func (opp *RegisterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringErrorFunc("failed to preprocess Register")

	fact, ok := op.Fact().(RegisterFact)
	if !ok {
		return ctx, nil, e(nil, "not RegisterFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e(err, "")
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot register and approve, %q: %w", fact.Sender(), err), nil
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

	starttime := proposal.StartTime() + delaytime
	endtime := starttime + delaytime + registerperiod

	blockmap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	blocktime := uint64(blockmap.Manifest().ProposedAt().Unix())

	if blocktime < starttime || endtime <= blocktime {
		return nil, base.NewBaseOperationProcessReasonError("not registration period, must in %d <= block(%d) < %d", starttime, blocktime, endtime), nil
	}

	switch st, found, err := getStateFunc(state.StateKeyRegisterList(fact.Contract(), fact.DAOID(), fact.ProposeID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find register list, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	case found:
		registers, err := state.StateRegisterListValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find register list value, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
		}

		var target base.Address
		if fact.Approved() != nil {
			target = fact.Approved()
		} else {
			target = fact.Sender()
		}

		for _, info := range registers {
			if info.Account().Equal(target) {
				if fact.Approved() != nil {
					for _, acc := range info.ApprovedBy() {
						if acc.Equal(fact.Sender()) {
							return nil, base.NewBaseOperationProcessReasonError("sender already approve the account, %q approved by %q", fact.Approved(), fact.Sender()), nil
						}
					}
				} else {
					for _, acc := range info.ApprovedBy() {
						if acc.Equal(fact.Sender()) {
							return nil, base.NewBaseOperationProcessReasonError("already registered account, %q", fact.Sender()), nil
						}
					}
				}
			}
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyApprovingList(fact.Contract(), fact.DAOID(), fact.ProposeID(), fact.Sender())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find approving list, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), fact.Sender(), err), nil
	case found:
		approving, err := state.StateApprovingListValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find approving list value, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), fact.Sender(), err), nil
		}

		var target base.Address
		if fact.Approved() != nil {
			target = fact.Approved()
		} else {
			target = fact.Sender()
		}

		for _, acc := range approving {
			if acc.Equal(target) {
				return nil, base.NewBaseOperationProcessReasonError("already approved account, %q approved by %q", target, fact.Sender()), nil
			}
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *RegisterProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringErrorFunc("failed to process Register")

	fact, ok := op.Fact().(RegisterFact)
	if !ok {
		return nil, nil, e(nil, "expected RegisterFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 3)

	switch st, found, err := getStateFunc(state.StateKeyRegisterList(fact.Contract(), fact.DAOID(), fact.ProposeID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find register list, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
	case found:
		registers, err := state.StateRegisterListValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find register list value, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), err), nil
		}

		var target base.Address
		if fact.Approved() != nil {
			target = fact.Approved()
		} else {
			target = fact.Sender()
		}

		for i, info := range registers {
			if info.Account().Equal(target) {
				accs := info.ApprovedBy()
				accs = append(accs, fact.Sender())

				registers[i] = state.NewRegisterInfo(target, accs)

				break
			}

			if i == len(registers)-1 {
				registers = append(registers, state.NewRegisterInfo(target, []base.Address{fact.Sender()}))
			}
		}

		sts[0] = currencystate.NewStateMergeValue(
			st.Key(),
			state.NewRegisterListStateValue(registers),
		)
	default:
		registers := make([]state.RegisterInfo, 1)
		if fact.Approved() != nil {
			registers[0] = state.NewRegisterInfo(fact.Approved(), []base.Address{fact.Sender()})
		} else {
			registers[0] = state.NewRegisterInfo(fact.Sender(), []base.Address{fact.Sender()})
		}
		sts[0] = currencystate.NewStateMergeValue(
			st.Key(),
			state.NewRegisterListStateValue(registers),
		)
	}

	switch st, found, err := getStateFunc(state.StateKeyApprovingList(fact.Contract(), fact.DAOID(), fact.ProposeID(), fact.Sender())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find approving list, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), fact.Sender()), nil
	case found:
		approving, err := state.StateApprovingListValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find approving list value, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposeID(), fact.Sender(), err), nil
		}

		if fact.Approved() != nil {
			approving = append(approving, fact.Approved())
		} else {
			approving = append(approving, fact.Sender())
		}

		sts[1] = currencystate.NewStateMergeValue(
			st.Key(),
			state.NewApprovingListStateValue(approving),
		)
	default:
		if fact.Approved() != nil {
			sts[1] = currencystate.NewStateMergeValue(
				st.Key(),
				state.NewApprovingListStateValue([]base.Address{fact.Approved()}),
			)
		} else {
			sts[1] = currencystate.NewStateMergeValue(
				st.Key(),
				state.NewApprovingListStateValue([]base.Address{fact.Sender()}),
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
	sts[2] = currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))))

	return sts, nil, nil
}

func (opp *RegisterProcessor) Close() error {
	registerProcessorPool.Put(opp)

	return nil
}
