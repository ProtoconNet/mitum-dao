package dao

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"sync"

	"github.com/ProtoconNet/mitum-dao/types"

	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var createDAOProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateDAOProcessor)
	},
}

func (CreateDAO) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CreateDAOProcessor struct {
	*base.BaseOperationProcessor
}

func NewCreateDAOProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new CreateDAOProcessor")

		nopp := createDAOProcessorPool.Get()
		opp, ok := nopp.(*CreateDAOProcessor)
		if !ok {
			return nil, errors.Errorf("expected CreateDAOProcessor, not %T", nopp)
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

func (opp *CreateDAOProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess CreateDAO")

	fact, ok := op.Fact().(CreateDAOFact)
	if !ok {
		return ctx, nil, e.Errorf("not CreateDAOFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender cannot be a contract account, %s: %w", fact.Sender(), err), nil
	}

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account not found, %s: %w", fact.Contract(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account value not found, %s: %w", fact.Contract(), err), nil
	}

	if !(ca.Owner().Equal(fact.sender) || ca.IsOperator(fact.Sender())) {
		return nil, base.NewBaseOperationProcessReasonError("sender is neither the owner nor the operator of the target contract account, %q", fact.Sender()), nil
	}

	if ca.IsActive() {
		return nil, base.NewBaseOperationProcessReasonError("a contract account is already used, %q", fact.Contract().String()), nil
	}

	if err := currencystate.CheckNotExistsState(state.StateKeyDesign(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao already exists, %s: %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency doesn't exist, %q: %w", fact.Currency(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.VotingPowerToken()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voting power token design not found, %q: %w", fact.VotingPowerToken(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *CreateDAOProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process CreateDAO")

	fact, ok := op.Fact().(CreateDAOFact)
	if !ok {
		return nil, nil, e.Errorf("expected CreateDAOFact, not %T", op.Fact())
	}

	policy := types.NewPolicy(
		fact.votingPowerToken, fact.threshold, fact.fee, fact.whitelist,
		fact.proposalReviewPeriod, fact.registrationPeriod, fact.preSnapshotPeriod, fact.votingPeriod,
		fact.postSnapshotPeriod, fact.executionDelayPeriod, fact.turnout, fact.quorum,
	)
	if err := policy.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid dao policy, %s: %w", fact.Contract(), err), nil
	}

	design := types.NewDesign(fact.option, policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid dao design, %s: %w", fact.Contract(), err), nil
	}

	var sts []base.StateMergeValue

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyDesign(fact.Contract()),
		state.NewDesignStateValue(design),
	))

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("target contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q; %w", fact.Contract(), err), nil
	}
	nca := ca.SetIsActive(true)

	sts = append(sts, currencystate.NewStateMergeValue(
		stateextension.StateKeyContractAccount(fact.Contract()),
		stateextension.NewContractAccountStateValue(nca),
	))

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

	return sts, nil, nil
}

func (opp *CreateDAOProcessor) Close() error {
	createDAOProcessorPool.Put(opp)

	return nil
}
