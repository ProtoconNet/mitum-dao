package dao

import (
	"context"
	"sync"
	"time"

	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	"github.com/ProtoconNet/mitum-dao/utils"

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
)

var registerProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RegisterProcessor)
	},
}

func (Register) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type RegisterProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewRegisterProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new RegisterProcessor")

		nopp := registerProcessorPool.Get()
		opp, ok := nopp.(*RegisterProcessor)
		if !ok {
			return nil, errors.Errorf("expected RegisterProcessor, not %T", nopp)
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

func (opp *RegisterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Register")

	fact, ok := op.Fact().(RegisterFact)
	if !ok {
		return ctx, nil, e.Errorf("not RegisterFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot delegate, %q: %w", fact.Sender(), err), nil
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
		return nil, base.NewBaseOperationProcessReasonError("dao value not found from state, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal state not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	blockMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(design.Policy(), p.Proposal(), blockMap)

	if period != types.Registration {
		return nil, base.NewBaseOperationProcessReasonError("not registration period, registration period start : %s, end : %s", time.Unix(start, 0), time.Unix(end, 0)), nil
	}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		accFound, delegatorFound, err := utils.HasFieldAndSliceValue(
			voters,
			"account",
			"delegators",
			fact.Delegated(),
			fact.Sender(),
		)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to check if delegator exists, %q: %w", fact.Delegated(), err), nil
		}

		if accFound && delegatorFound {
			return nil, base.NewBaseOperationProcessReasonError(
				"sender already delegates the account, %q delegated by %q",
				fact.Delegated(),
				fact.Sender(),
			), nil
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyDelegators(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find delegators state, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), fact.Sender(), err), nil
	case found:
		delegators, err := state.StateDelegatorsValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find delegators value from state, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), fact.Sender(), err), nil
		}

		if found, err := utils.HasFieldValue(delegators, "account", fact.Sender()); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to check account value %s in delegators", fact.Sender()), nil
		} else if found {
			return nil, base.NewBaseOperationProcessReasonError("sender %s already delegates", fact.Sender()), nil
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	return ctx, nil, nil
}

func (opp *RegisterProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Register")

	fact, ok := op.Fact().(RegisterFact)
	if !ok {
		return nil, nil, e.Errorf("expected RegisterFact, not %T", op.Fact())
	}

	sts := make([]base.StateMergeValue, 3)

	var voters []types.VoterInfo
	var key string
	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		vs, err := state.StateVotersValue(st)
		key = st.Key()
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		for i, info := range vs {
			if info.Account().Equal(fact.Delegated()) {
				delegators := info.Delegators()
				delegators = append(delegators, fact.Sender())
				vs[i] = types.NewVoterInfo(fact.Delegated(), delegators)

				break
			}

			if i == len(vs)-1 {
				vs = append(vs, types.NewVoterInfo(fact.Delegated(), []base.Address{fact.Sender()}))
			}
		}
		voters = vs
	default:
		var vs []types.VoterInfo
		vs = append(vs, types.NewVoterInfo(fact.Delegated(), []base.Address{fact.Sender()}))
		voters = vs
	}
	sts[0] = currencystate.NewStateMergeValue(
		key,
		state.NewVotersStateValue(voters),
	)

	switch st, found, err := getStateFunc(state.StateKeyDelegators(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find delegators state, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), fact.Sender(), err), nil
	case found:
		delegators, err := state.StateDelegatorsValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find delegators value from state, %s-%s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), fact.Sender(), err), nil
		}

		delegators = append(delegators, types.NewDelegatorInfo(fact.Sender(), fact.Delegated()))

		sts[1] = currencystate.NewStateMergeValue(
			st.Key(),
			state.NewDelegatorsStateValue(delegators),
		)
	default:
		sts[1] = currencystate.NewStateMergeValue(
			st.Key(),
			state.NewDelegatorsStateValue([]types.DelegatorInfo{types.NewDelegatorInfo(fact.Sender(), fact.Delegated())}))
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
