package dao

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var preSnapProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(PreSnapProcessor)
	},
}

func (PreSnap) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type PreSnapProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewPreSnapProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new PreSnapProcessor")

		nopp := preSnapProcessorPool.Get()
		opp, ok := nopp.(*PreSnapProcessor)
		if !ok {
			return nil, errors.Errorf("expected PreSnapProcessor, not %T", nopp)
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

func (opp *PreSnapProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(PreSnapFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", PreSnapFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id, %v", fact.Currency())), nil
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

	if st, err := currencystate.ExistsState(state.StateKeyDesign(fact.Contract()), "design", getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("dao design, %v",
				fact.Contract(),
			)), nil
	} else if _, err := state.StateDesignValue(st); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("dao design, %v",
				fact.Contract(),
			)), nil
	}

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateNF).Errorf("proposal, %s,%v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	}

	p, err := state.StateProposalValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateValInvalid).Errorf("proposal, %s,%v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	}

	if p.Status() == types.Canceled {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("already canceled proposal, %s, %q", fact.Contract(), fact.ProposalID())), nil
	} else if p.Status() == types.PreSnapped {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("already preSnapped, %s, %q", fact.Contract(), fact.ProposalID())), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVoters(fact.Contract(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("voters, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyDelegators(fact.Contract(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("delegators state not found, %s, %q: %w", fact.Contract(), fact.ProposalID(), err)), nil
	}

	if found, err := currencystate.CheckNotExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID()), getStateFunc); found {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateE).
				Errorf("voting power box state already created, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *PreSnapProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process PreSnap")

	fact, ok := op.Fact().(PreSnapFact)
	if !ok {
		return nil, nil, e.Errorf("expected PreSnapFact, not %T", op.Fact())
	}

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "key of proposal", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
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

	period, start, end := types.GetPeriodOfCurrentTime(p.Policy(), p.Proposal(), types.PreSnapshot, blockMap)
	if period != types.PreSnapshot {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within the PreSnapshotPeriod, PreSnapshotPeriod; start(%d), end(%d), but now(%d)", start, end, blockMap.Manifest().ProposedAt().Unix()), nil
	}

	var sts []base.StateMergeValue

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

	//st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "key of proposal", getStateFunc)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("proposal not found, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	//}
	//
	//p, err := state.StateProposalValue(st)
	//if err != nil {
	//	return nil, base.NewBaseOperationProcessReasonError("proposal value not found from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	//}

	var votingPowerBox types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		if vb, err := state.StateVotingPowerBoxValue(st); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
		} else {
			votingPowerBox = vb
		}
	default:
		votingPowerBox = types.NewVotingPowerBox(common.ZeroBig, map[string]types.VotingPower{})
	}

	votingPowerToken := p.Policy().Token()

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
		}

		total := common.ZeroBig
		votingPowers := map[string]types.VotingPower{}
		for _, info := range voters {
			votingPower := common.ZeroBig

			for _, delegator := range info.Delegators() {
				st, err = currencystate.ExistsState(currency.StateKeyBalance(delegator, votingPowerToken), "key of balance", getStateFunc)
				if err != nil {
					continue
				}

				b, err := currency.StateBalanceValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to find balance value of the delegator from state, %s, %q: %w", delegator, votingPowerToken, err), nil
				}

				votingPower = votingPower.Add(b.Big())
			}

			v, found := votingPowers[info.Account().String()]
			if found {
				votingPower = v.Amount().Add(votingPower)
			}

			votingPowers[info.Account().String()] = types.NewVotingPower(info.Account(), votingPower)
		}

		for _, v := range votingPowers {
			total = total.Add(v.Amount())
		}
		votingPowerBox.SetVotingPowers(votingPowers)
		votingPowerBox.SetTotal(total)
	}

	st, err = currencystate.ExistsState(currency.StateKeyCurrencyDesign(votingPowerToken), "key of currency design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency design, %q: %w", votingPowerToken, err), nil
	}

	currencyDesign, err := currency.StateCurrencyDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency design value from state, %q: %w", votingPowerToken, err), nil
	}

	actualTurnoutCount := p.Policy().Turnout().Quorum(currencyDesign.Aggregate())
	if votingPowerBox.Total().Compare(actualTurnoutCount) < 0 {
		sts = append(sts, currencystate.NewStateMergeValue(
			state.StateKeyProposal(fact.Contract(), fact.ProposalID()),
			state.NewProposalStateValue(types.Canceled, p.Proposal(), p.Policy()),
		))
	} else {
		sts = append(sts,
			currencystate.NewStateMergeValue(
				state.StateKeyProposal(fact.Contract(), fact.ProposalID()),
				state.NewProposalStateValue(types.PreSnapped, p.Proposal(), p.Policy()),
			),
			currencystate.NewStateMergeValue(
				state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID()),
				state.NewVotingPowerBoxStateValue(votingPowerBox),
			),
		)
	}

	return sts, nil, nil
}

func (opp *PreSnapProcessor) Close() error {
	preSnapProcessorPool.Put(opp)

	return nil
}
