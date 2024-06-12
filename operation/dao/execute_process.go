package dao

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	crcystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var executeProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ExecuteProcessor)
	},
}

func (Execute) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ExecuteProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewExecuteProcessor(getLastBlockFunc processor.GetLastBlockFunc) crcytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new ExecuteProcessor")

		nopp := executeProcessorPool.Get()
		opp, ok := nopp.(*ExecuteProcessor)
		if !ok {
			return nil, errors.Errorf("expected ExecuteProcessor, not %T", nopp)
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

func (opp *ExecuteProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ExecuteFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", ExecuteFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := crcystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id %q", fact.Currency())), nil
	}

	if _, _, aErr, cErr := crcystate.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v", cErr)), nil
	}

	_, _, aErr, cErr := crcystate.ExistsCAccount(fact.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", cErr)), nil
	}

	if st, err := crcystate.ExistsState(state.StateKeyDesign(fact.Contract()), "design", getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("dao design state for contract account %v",
				fact.Contract(),
			)), nil
	} else if _, err := state.StateDesignValue(st); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("dao design state value for contract account %v",
				fact.Contract(),
			)), nil
	}

	st, err := crcystate.ExistsState(
		state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateNF).Errorf(
				"proposal state %q for contract account %v", fact.ProposalID(), fact.Contract())), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateValInvalid).Errorf(
				"proposal state value %q for contract account %v", fact.ProposalID(), fact.Contract())), nil
	}

	if p.Status() == types.Canceled {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("already canceled proposal %q for contract account %v",
					fact.ProposalID(), fact.Contract())), nil
	} else if p.Status() == types.Rejected {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("already rejected proposal %q for contract account %v",
					fact.ProposalID(), fact.Contract())), nil
	} else if p.Status() == types.Executed {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("already executed proposal %q for contract account %v",
					fact.ProposalID(), fact.Contract())), nil
	}

	if err := crcystate.CheckExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("voting power box of proposal %q in contract account %v", fact.ProposalID(), fact.Contract())), nil
	}

	if err := crcystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *ExecuteProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Execute")

	fact, ok := op.Fact().(ExecuteFact)
	if !ok {
		return nil, nil, e.Errorf("expected ExecuteFact, not %T", op.Fact())
	}

	st, err := crcystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"proposal not found, %s, %q: %w", fact.Contract(), fact.ProposalID(), err,
		), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"proposal value not found from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err,
		), nil
	}

	blockMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(p.Policy(), p.Proposal(), types.Execute, blockMap)
	if period != types.Execute {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within the Execution, Execution period; start(%d), end(%d), but now(%d)", start, end, blockMap.Manifest().ProposedAt().Unix()), nil
	}

	var sts []base.StateMergeValue

	{ // caculate operation fee
		currencyPolicy, err := crcystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
		}

		fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"failed to check fee of currency, %q; %w",
				fact.Currency(),
				err,
			), nil
		}

		senderBalSt, err := crcystate.ExistsState(
			currency.StateKeyBalance(fact.Sender(), fact.Currency()),
			"key of sender balance",
			getStateFunc,
		)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"sender balance not found, %q; %w",
				fact.Sender(),
				err,
			), nil
		}

		switch senderBal, err := currency.StateBalanceValue(senderBalSt); {
		case err != nil:
			return nil, base.NewBaseOperationProcessReasonError(
				"failed to get balance value, %q; %w",
				currency.StateKeyBalance(fact.Sender(), fact.Currency()),
				err,
			), nil
		case senderBal.Big().Compare(fee) < 0:
			return nil, base.NewBaseOperationProcessReasonError(
				"not enough balance of sender, %q",
				fact.Sender(),
			), nil
		}

		v, ok := senderBalSt.Value().(currency.BalanceStateValue)
		if !ok {
			return nil, base.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", senderBalSt.Value()), nil
		}

		if currencyPolicy.Feeer().Receiver() != nil {
			if err := crcystate.CheckExistsState(currency.StateKeyAccount(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
				return nil, nil, err
			} else if feeRcvrSt, found, err := getStateFunc(currency.StateKeyBalance(currencyPolicy.Feeer().Receiver(), fact.currency)); err != nil {
				return nil, nil, err
			} else if !found {
				return nil, nil, errors.Errorf("feeer receiver %s not found", currencyPolicy.Feeer().Receiver())
			} else if feeRcvrSt.Key() != senderBalSt.Key() {
				r, ok := feeRcvrSt.Value().(currency.BalanceStateValue)
				if !ok {
					return nil, nil, errors.Errorf("expected %T, not %T", currency.BalanceStateValue{}, feeRcvrSt.Value())
				}
				sts = append(sts, common.NewBaseStateMergeValue(
					feeRcvrSt.Key(),
					currency.NewAddBalanceStateValue(r.Amount.WithBig(fee)),
					func(height base.Height, st base.State) base.StateValueMerger {
						return currency.NewBalanceStateValueMerger(height, feeRcvrSt.Key(), fact.currency, st)
					},
				))

				sts = append(sts, common.NewBaseStateMergeValue(
					senderBalSt.Key(),
					currency.NewDeductBalanceStateValue(v.Amount.WithBig(fee)),
					func(height base.Height, st base.State) base.StateValueMerger {
						return currency.NewBalanceStateValueMerger(height, senderBalSt.Key(), fact.currency, st)
					},
				))
			}
		}
	}

	if p.Status() != types.Completed {
		sts = append(sts,
			crcystate.NewStateMergeValue(
				st.Key(),
				state.NewProposalStateValue(types.Canceled, p.Proposal(), p.Policy()),
			),
		)

		return sts, nil, nil
	}

	sts = append(sts, crcystate.NewStateMergeValue(
		state.StateKeyProposal(fact.Contract(), fact.ProposalID()),
		state.NewProposalStateValue(types.Executed, p.Proposal(), p.Policy()),
	))

	if p.Proposal().Option() == types.ProposalCrypto {
		cp, _ := p.Proposal().(types.CryptoProposal)

		switch cp.CallData().Type() {
		case types.CalldataTransfer:
			cd, ok := cp.CallData().(types.TransferCallData)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected TransferCalldata, not %T", cp.CallData()), nil
			}

			if err := crcystate.CheckExistsState(currency.StateKeyAccount(cd.Sender()), getStateFunc); err != nil {
				return nil, base.NewBaseOperationProcessReasonError("calldata sender not found, %s: %w", cd.Sender(), err), nil
			}

			if err := crcystate.CheckExistsState(currency.StateKeyAccount(cd.Receiver()), getStateFunc); err != nil {
				return nil, base.NewBaseOperationProcessReasonError("calldata receiver not found, %s: %w", cd.Receiver(), err), nil
			}

			st, err = crcystate.ExistsState(currency.StateKeyBalance(cd.Sender(), cd.Amount().Currency()), "key of balance", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("failed to find calldata sender balance, %s, %q: %w", cd.Sender(), cd.Amount().Currency(), err), nil
			}

			sb, err := currency.StateBalanceValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("failed to find calldata sender balance value, %s, %q: %w", cd.Sender(), cd.Amount().Currency(), err), nil
			}

			if sb.Big().Compare(cd.Amount().Big()) >= 0 {
				sts = append(sts, crcystate.NewStateMergeValue(
					st.Key(),
					currency.NewBalanceStateValue(
						crcytypes.NewAmount(sb.Big().Sub(cd.Amount().Big()), cd.Amount().Currency()),
					),
				))

				switch st, found, err := getStateFunc(currency.StateKeyBalance(cd.Receiver(), cd.Amount().Currency())); {
				case err != nil:
					return nil, base.NewBaseOperationProcessReasonError("failed to find calldata receiver balance, %s, %q: %w", cd.Receiver(), cd.Amount().Currency(), err), nil
				case found:
					rb, err := currency.StateBalanceValue(st)
					if err != nil {
						return nil, base.NewBaseOperationProcessReasonError("failed to find calldata receiver balance value, %s, %q: %w", cd.Receiver(), cd.Amount().Currency(), err), nil
					}

					sts = append(sts, crcystate.NewStateMergeValue(
						st.Key(),
						currency.NewBalanceStateValue(
							crcytypes.NewAmount(rb.Big().Add(cd.Amount().Big()), cd.Amount().Currency()),
						),
					))
				default:
					sts = append(sts, crcystate.NewStateMergeValue(
						st.Key(),
						currency.NewBalanceStateValue(
							crcytypes.NewAmount(common.ZeroBig.Add(cd.Amount().Big()), cd.Amount().Currency()),
						),
					))
				}
			}
		case types.CalldataGovernance:
			cd, ok := cp.CallData().(types.GovernanceCallData)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected GovernanceCalldata, not %T", cp.CallData()), nil
			}

			st, err := crcystate.ExistsState(state.StateKeyDesign(fact.Contract()), "key of design", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("dao design not found, %s: %w", fact.Contract(), err), nil
			}

			design, err := state.StateDesignValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("dao design value not found, %s: %w", fact.Contract(), err), nil
			}

			nd := types.NewDesign(design.Option(), cd.Policy())

			if err := nd.IsValid(nil); err != nil {
				sts = append(sts, crcystate.NewStateMergeValue(
					state.StateKeyDesign(fact.Contract()),
					state.NewDesignStateValue(
						nd,
					),
				))
			}
		default:
			return nil, base.NewBaseOperationProcessReasonError("invalid calldata, %s, %q", fact.Contract(), fact.ProposalID()), nil
		}
	}

	return sts, nil, nil
}

func (opp *ExecuteProcessor) Close() error {
	executeProcessorPool.Put(opp)

	return nil
}
