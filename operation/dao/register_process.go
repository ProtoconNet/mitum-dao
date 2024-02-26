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
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender cannot be a contract account, %s: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckExistsState(extensioncurrency.StateKeyContractAccount(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao contract account not found, %s: %w", fact.Contract(), err), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("fee currency doesn't exist, %q: %w", fact.Currency(), err), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyDesign(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao design not found, %s: %w", fact.Contract(), err), nil
	}

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal state not found, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s, %q", fact.Contract(), fact.ProposalID()), nil
	}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
		}

		var voter types.VoterInfo = types.VoterInfo{}

		for _, v := range voters {
			if !fact.Delegated().Equal(v.Account()) {
				continue
			}
			for _, d := range voter.Delegators() {
				if fact.Sender().Equal(d) {
					return nil, base.NewBaseOperationProcessReasonError(
						"sender already delegates the account, %s delegated by %s",
						fact.Delegated(),
						fact.Sender(),
					), nil
				}
			}
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyDelegators(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find delegators state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		delegators, err := state.StateDelegatorsValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find delegators value from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
		}

		for _, delegator := range delegators {
			if delegator.Account().Equal(fact.Sender()) {
				return nil, base.NewBaseOperationProcessReasonError("sender %s already delegates, %s, %q: %w", fact.Sender(), fact.Contract(), fact.ProposalID()), nil
			}
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

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal state not found, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	}

	blockMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(p.Policy(), p.Proposal(), types.Registration, blockMap)
	if period != types.Registration {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within the Registration period, Registration period; start(%d), end(%d), but now(%d)", start, end, blockMap.Manifest().ProposedAt().Unix()), nil
	}

	sts := []base.StateMergeValue{}

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

	var voters []types.VoterInfo
	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		vs, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
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

	sts = append(sts,
		common.NewBaseStateMergeValue(
			state.StateKeyVoters(fact.Contract(), fact.ProposalID()),
			state.NewVotersStateValue(voters),
			func(height base.Height, st base.State) base.StateValueMerger {
				return state.NewVotersStateValueMerger(height, state.StateKeyVoters(fact.Contract(), fact.ProposalID()), st)
			},
		),
	)

	switch st, found, err := getStateFunc(state.StateKeyDelegators(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find delegators state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		delegators, err := state.StateDelegatorsValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find delegators value from state, %s,%q: %w", fact.Contract(), fact.ProposalID(), err), nil
		}

		delegators = append(delegators, types.NewDelegatorInfo(fact.Sender(), fact.Delegated()))

		sts = append(sts,
			common.NewBaseStateMergeValue(
				state.StateKeyDelegators(fact.Contract(), fact.ProposalID()),
				state.NewDelegatorsStateValue(delegators),
				func(height base.Height, st base.State) base.StateValueMerger {
					return state.NewDelegatorsStateValueMerger(height, state.StateKeyDelegators(fact.Contract(), fact.ProposalID()), st)
				},
			),
		)
	default:
		sts = append(sts,
			common.NewBaseStateMergeValue(
				state.StateKeyDelegators(fact.Contract(), fact.ProposalID()),
				state.NewDelegatorsStateValue([]types.DelegatorInfo{types.NewDelegatorInfo(fact.Sender(), fact.Delegated())}),
				func(height base.Height, st base.State) base.StateValueMerger {
					return state.NewDelegatorsStateValueMerger(height, state.StateKeyDelegators(fact.Contract(), fact.ProposalID()), st)
				},
			),
		)
	}

	return sts, nil, nil
}

func (opp *RegisterProcessor) Close() error {
	registerProcessorPool.Put(opp)

	return nil
}
