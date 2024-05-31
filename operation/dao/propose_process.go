package dao

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-dao/types"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var proposeProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ProposeProcessor)
	},
}

func (Propose) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ProposeProcessor struct {
	*base.BaseOperationProcessor
}

func NewProposeProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new ProposeProcessor")

		nopp := proposeProcessorPool.Get()
		opp, ok := nopp.(*ProposeProcessor)
		if !ok {
			return nil, errors.Errorf("expected ProposeProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *ProposeProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(ProposeFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", ProposeFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if _, _, aErr, cErr := currencystate.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v", cErr)), nil
	}

	_, _, aErr, cErr := currencystate.ExistsCAccount(fact.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", cErr)), nil
	}

	if found, err := currencystate.CheckNotExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), getStateFunc); found {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateE).
				Errorf("proposal already exists, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	}

	required := map[string]common.Big{}

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id, %v", fact.Currency())), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Errorf("check fee of currency, %q: %v", fact.Currency(), err)), nil
	}

	required[fact.currency.String()] = fee

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract()), "design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("dao design not found, %s: %w", fact.Contract(), err)), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
				Errorf("dao design, %s: %v", fact.Contract(), err)), nil
	}

	if design.Option() != fact.Proposal().Option() {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("dao option != proposal option, dao(%s) != proposal(%s)", design.Option(), fact.Proposal().Option())), nil
	}

	votingPowerToken := design.Policy().Token()
	threshold := design.Policy().Threshold()
	proposeFee := design.Policy().Fee()
	whitelist := design.Policy().Whitelist()

	if _, found := required[votingPowerToken.String()]; !found {
		required[votingPowerToken.String()] = common.ZeroBig
	}

	if _, found := required[proposeFee.Currency().String()]; !found {
		required[proposeFee.Currency().String()] = common.ZeroBig
	}

	required[votingPowerToken.String()] = required[votingPowerToken.String()].Add(threshold)
	required[proposeFee.Currency().String()] = required[proposeFee.Currency().String()].Add(proposeFee.Big())

	for k, v := range required {
		st, err = currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), currencytypes.CurrencyID(k)), "key of sender balance", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMStateNF).
					Errorf("sender balance, %s, %v: %v", fact.Sender(), fact.Currency(), err)), nil
		}

		switch b, err := currency.StateBalanceValue(st); {
		case err != nil:
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
					Errorf("failed to get balance value, %s, %v: %v", fact.Sender(), fact.Currency(), err)), nil
		case b.Big().Compare(v) < 0:
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMValOOR).
					Errorf("not enough balance of sender, %s, %q", fact.Sender(), fact.Currency())), nil
		}
	}

	if whitelist.Active() && !whitelist.IsExist(fact.Sender()) {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMAccountNAth).
				Errorf("sender not in whitelist, %s", fact.Sender())), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *ProposeProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Propose")

	fact, ok := op.Fact().(ProposeFact)
	if !ok {
		return nil, nil, e.Errorf("expected ProposeFact, not %T", op.Fact())
	}

	var sts []base.StateMergeValue

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao not found, %s: %w", fact.Contract(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao value not found, %s: %w", fact.Contract(), err), nil
	}

	proposeFee := design.Policy().Fee()

	sts = append(sts,
		currencystate.NewStateMergeValue(
			state.StateKeyProposal(fact.Contract(), fact.ProposalID()),
			state.NewProposalStateValue(types.Proposed, fact.Proposal(), design.Policy()),
		),
	)

	{ // caculate operation fee
		currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
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

		senderBalSt, err := currencystate.ExistsState(
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
			if err := currencystate.CheckExistsState(currency.StateKeyAccount(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
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

	st, err = currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), proposeFee.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance for propose fee not found, %s, %q: %w", fact.Sender(), proposeFee.Currency(), err), nil
	}
	fBalance, err := currency.StateBalanceValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance value for propose fee not found, %s, %q: %w", fact.Sender(), proposeFee.Currency(), err), nil
	}

	sts = append(sts,
		common.NewBaseStateMergeValue(
			st.Key(),
			currency.NewDeductBalanceStateValue(fBalance.WithBig(proposeFee.Big())),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(height, st.Key(), fact.currency, st)
			},
		),
	)

	var cBalance currencytypes.Amount
	var cBalanceKey string
	switch st, found, err := getStateFunc(currency.StateKeyBalance(fact.Contract(), proposeFee.Currency())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("contract account balance for propose fee not found, %s, %q: %w", fact.Contract(), proposeFee.Currency(), err), nil
	case !found:
		cBalance = currencytypes.NewAmount(common.ZeroBig, proposeFee.Currency())
		cBalanceKey = st.Key()
	default:
		cBalance, err = currency.StateBalanceValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("contract balance value for propose fee not found, %s, %q: %w", fact.Contract(), proposeFee.Currency(), err), nil
		}
		cBalanceKey = st.Key()
	}

	sts = append(sts,
		currencystate.NewStateMergeValue(cBalanceKey, currency.NewBalanceStateValue(cBalance.WithBig(cBalance.Big().Add(proposeFee.Big())))),
	)

	return sts, nil, nil
}

func (opp *ProposeProcessor) Close() error {
	proposeProcessorPool.Put(opp)

	return nil
}
