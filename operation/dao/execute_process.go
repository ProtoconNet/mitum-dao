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

func NewExecuteProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
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
	e := util.StringError("failed to preprocess Execute")

	fact, ok := op.Fact().(ExecuteFact)
	if !ok {
		return ctx, nil, e.Errorf("not ExecuteFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot execute, %q: %w", fact.Sender(), err), nil
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
		return nil, base.NewBaseOperationProcessReasonError("dao design value not found from state, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() == types.Canceled {
		return nil, base.NewBaseOperationProcessReasonError("already canceled proposal, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	} else if p.Status() == types.Rejected {
		return nil, base.NewBaseOperationProcessReasonError("rejected proposal, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	} else if p.Status() == types.Executed {
		return nil, base.NewBaseOperationProcessReasonError("already executed, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	blocMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(design.Policy(), p.Proposal(), blocMap)
	if period != types.Execute {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within the Execution, Execution period start : %d, end %d", start, end), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
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

	var sts []base.StateMergeValue

	{ // caculate operation fee
		policy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find currency policy, %s: %w", fact.Currency(), err), nil
		}

		fee, err := policy.Feeer().Fee(common.ZeroBig)
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
		sts = append(sts, currencystate.NewStateMergeValue(sb.Key(), currency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee)))))
	}

	st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract(), fact.DAOID()), "key of design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao not found, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("dao design value not found from state, %s-%s: %w", fact.Contract(), fact.DAOID(), err), nil
	}

	st, err = currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if p.Status() != types.Completed {
		sts = append(sts,
			currencystate.NewStateMergeValue(
				st.Key(),
				state.NewProposalStateValue(types.Canceled, p.Proposal()),
			),
		)

		return sts, nil, nil
	}

	var vpb types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		vpb, err = state.StateVotingPowerBoxValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}
	default:
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	actualTurnoutQuorum := design.Policy().Quorum().Quorum(vpb.Total())
	execute := false

	if p.Proposal().Type() == types.ProposalCrypto {
		agree, reject := vpb.Result()[0], vpb.Result()[1]

		if agree.Compare(actualTurnoutQuorum) >= 0 && agree.Compare(reject) >= 0 {
			sts = append(sts, currencystate.NewStateMergeValue(
				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
				state.NewProposalStateValue(types.Executed, p.Proposal()),
			))
			execute = true
		} else {
			sts = append(sts, currencystate.NewStateMergeValue(
				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
				state.NewProposalStateValue(types.Rejected, p.Proposal()),
			))
		}
	} else if p.Proposal().Type() == types.ProposalBiz {
		options := p.Proposal().Options() - 1
		result := vpb.Result()

		var ok = false
		var selected uint8 = 0
		var i uint8 = 0

		for ; i < options; i++ {
			if result[i].Compare(actualTurnoutQuorum) >= 0 {
				if !ok {
					ok, selected = true, i
					continue
				}

				if result[selected].Compare(result[i]) < 0 {
					selected = i
				}
			}
		}

		sts = append(sts, currencystate.NewStateMergeValue(
			state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
			state.NewProposalStateValue(types.Executed, p.Proposal()),
		))
	}

	if execute {
		cp, _ := p.Proposal().(types.CryptoProposal)

		switch cp.CallData().Type() {
		case types.CalldataTransfer:
			cd, ok := cp.CallData().(types.TransferCallData)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected TransferCalldata, not %T", cp.CallData()), nil
			}

			if err := currencystate.CheckExistsState(currency.StateKeyAccount(cd.Sender()), getStateFunc); err != nil {
				return nil, base.NewBaseOperationProcessReasonError("calldata sender not found, %s: %w", cd.Sender(), err), nil
			}

			if err := currencystate.CheckExistsState(currency.StateKeyAccount(cd.Receiver()), getStateFunc); err != nil {
				return nil, base.NewBaseOperationProcessReasonError("calldata receiver not found, %s: %w", cd.Receiver(), err), nil
			}

			st, err = currencystate.ExistsState(currency.StateKeyBalance(cd.Sender(), cd.Amount().Currency()), "key of balance", getStateFunc)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("failed to find calldata sender balance, %s, %s: %w", cd.Sender(), cd.Amount().Currency(), err), nil
			}

			sb, err := currency.StateBalanceValue(st)
			if err != nil {
				return nil, base.NewBaseOperationProcessReasonError("failed to find calldata sender balance value, %s, %s: %w", cd.Sender(), cd.Amount().Currency(), err), nil
			}

			if sb.Big().Compare(cd.Amount().Big()) >= 0 {
				sts = append(sts, currencystate.NewStateMergeValue(
					st.Key(),
					currency.NewBalanceStateValue(
						currencytypes.NewAmount(sb.Big().Sub(cd.Amount().Big()), cd.Amount().Currency()),
					),
				))

				switch st, found, err := getStateFunc(currency.StateKeyBalance(cd.Receiver(), cd.Amount().Currency())); {
				case err != nil:
					return nil, base.NewBaseOperationProcessReasonError("failed to find calldata receiver balance, %s, %s: %w", cd.Receiver(), cd.Amount().Currency(), err), nil
				case found:
					rb, err := currency.StateBalanceValue(st)
					if err != nil {
						return nil, base.NewBaseOperationProcessReasonError("failed to find calldata receiver balance value, %s, %s: %w", cd.Receiver(), cd.Amount().Currency(), err), nil
					}

					sts = append(sts, currencystate.NewStateMergeValue(
						st.Key(),
						currency.NewBalanceStateValue(
							currencytypes.NewAmount(rb.Big().Add(cd.Amount().Big()), cd.Amount().Currency()),
						),
					))
				default:
					sts = append(sts, currencystate.NewStateMergeValue(
						st.Key(),
						currency.NewBalanceStateValue(
							currencytypes.NewAmount(common.ZeroBig.Add(cd.Amount().Big()), cd.Amount().Currency()),
						),
					))
				}
			}
		case types.CalldataGovernance:
			cd, ok := cp.CallData().(types.GovernanceCallData)
			if !ok {
				return nil, base.NewBaseOperationProcessReasonError("expected GovernanceCalldata, not %T", cp.CallData()), nil
			}

			nd := types.NewDesign(design.Option(), design.DAOID(), cd.Policy())

			err := nd.IsValid(nil)
			if err == nil {
				sts = append(sts, currencystate.NewStateMergeValue(
					state.StateKeyDesign(fact.Contract(), fact.DAOID()),
					state.NewDesignStateValue(
						nd,
					),
				))
			}
		default:
			return nil, base.NewBaseOperationProcessReasonError("invalid calldata, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
		}
	}

	return sts, nil, nil
}

func (opp *ExecuteProcessor) Close() error {
	executeProcessorPool.Put(opp)

	return nil
}
