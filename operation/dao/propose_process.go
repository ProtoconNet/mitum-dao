package dao

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-dao/types"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextionsion "github.com/ProtoconNet/mitum-currency/v3/state/extension"
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
	e := util.StringError("failed to preprocess Propose")

	fact, ok := op.Fact().(ProposeFact)
	if !ok {
		return ctx, nil, e.Errorf("not ProposeFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextionsion.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot propose proposal, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(stateextionsion.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account not found, %s: %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckNotExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal already exists, %s, %q, %q: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	required := map[string]common.Big{}

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	required[fact.currency.String()] = fee

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao not found, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao value not found, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
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
			return nil, base.NewBaseOperationProcessReasonError("sender balance not found, %s, %q: %w", fact.Sender(), fact.Currency(), err), nil
		}

		switch b, err := currency.StateBalanceValue(st); {
		case err != nil:
			return nil, base.NewBaseOperationProcessReasonError("failed to get balance value, %s, %q: %w", fact.Sender(), fact.Currency(), err), nil
		case b.Big().Compare(v) < 0:
			return nil, base.NewBaseOperationProcessReasonError("not enough balance of sender, %s, %q", fact.Sender(), fact.Currency()), nil
		}
	}

	if whitelist.Active() && !whitelist.IsExist(fact.Sender()) {
		return nil, base.NewBaseOperationProcessReasonError("sender not in whitelist, %s", fact.Sender()), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
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

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao not found, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao value not found, %s, %q: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	proposeFee := design.Policy().Fee()

	sts = append(sts,
		currencystate.NewStateMergeValue(
			state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
			state.NewProposalStateValue(types.Proposed, fact.Proposal(), design.Policy()),
		),
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
	balance, err := currency.StateBalanceValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance value not found, %s, %q: %w", fact.Sender(), fact.Currency(), err), nil
	}
	sb := currency.NewBalanceStateValue(balance)

	sts = append(sts,
		currencystate.NewStateMergeValue(st.Key(), currency.NewBalanceStateValue(sb.Amount.WithBig(sb.Amount.Big().Sub(fee)))),
	)

	st, err = currencystate.ExistsState(currency.StateKeyBalance(fact.Sender(), proposeFee.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance for propose fee not found, %s, %q: %w", fact.Sender(), proposeFee.Currency(), err), nil
	}
	balance, err = currency.StateBalanceValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender balance value for propose fee not found, %s, %q: %w", fact.Sender(), proposeFee.Currency(), err), nil
	}
	fb := currency.NewBalanceStateValue(balance)

	sts = append(sts,
		currencystate.NewStateMergeValue(st.Key(), currency.NewBalanceStateValue(fb.Amount.WithBig(fb.Amount.Big().Sub(proposeFee.Big())))),
	)

	return sts, nil, nil
}

func (opp *ProposeProcessor) Close() error {
	proposeProcessorPool.Put(opp)

	return nil
}
