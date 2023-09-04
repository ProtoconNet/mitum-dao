package dao

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-dao/types"

	"github.com/ProtoconNet/mitum-currency/v3/common"
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
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot create dao, %s: %w", fact.Sender(), err), nil
	}

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account not found, %s: %w", fact.Contract(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account value not found, %s: %w", fact.Contract(), err), nil
	}

	if !ca.Owner().Equal(fact.Sender()) {
		return nil, base.NewBaseOperationProcessReasonError("not contract account owner, %s", fact.Sender()), nil
	}

	if ca.IsActive() {
		return nil, base.NewBaseOperationProcessReasonError("a design is already registered, %q", fact.Contract().String()), nil
	}

	if err := currencystate.CheckNotExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao already exists, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
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
		return nil, base.NewBaseOperationProcessReasonError("invalid dao policy, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design := types.NewDesign(fact.option, fact.DAOID(), policy)
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid dao design, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	sts := make([]base.StateMergeValue, 3)

	sts[0] = currencystate.NewStateMergeValue(
		state.StateKeyDesign(fact.Contract(), fact.DAOID()),
		state.NewDesignStateValue(design),
	)

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("target contract account not found, %q; %w", fact.Contract(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q; %w", fact.Contract(), err), nil
	}
	ca.SetIsActive(true)

	sts[1] = currencystate.NewStateMergeValue(
		stateextension.StateKeyContractAccount(fact.Contract()),
		stateextension.NewContractAccountStateValue(ca),
	)

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
	sts[2] = currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))))

	return sts, nil, nil
}

func (opp *CreateDAOProcessor) Close() error {
	createDAOProcessorPool.Put(opp)

	return nil
}
