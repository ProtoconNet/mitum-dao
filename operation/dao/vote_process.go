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

var voteProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(VoteProcessor)
	},
}

func (Vote) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type VoteProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewVoteProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new VoteProcessor")

		nopp := voteProcessorPool.Get()
		opp, ok := nopp.(*VoteProcessor)
		if !ok {
			return nil, errors.Errorf("expected VoteProcessor, not %T", nopp)
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

func (opp *VoteProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(VoteFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", VoteFact{}, op.Fact())), nil
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

	if err := currencystate.CheckExistsState(currency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMCurrencyNF).Errorf("fee currency id, %v", fact.Currency())), nil
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
				Errorf("already canceled proposal, %s, %v", fact.Contract(), fact.ProposalID())), nil
	}

	if p.Status() != types.PreSnapped {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("proposal not in pre-snapped status, %s, %v, %v", fact.Contract(), fact.ProposalID(), p.Status())), nil
	}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("voters, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	case !found:
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("voters, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	default:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
					Errorf("voters, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
		}

		for i, v := range voters {
			if v.Account().Equal(fact.Sender()) {
				break
			}

			if i == len(voters)-1 {
				return nil, base.NewBaseOperationProcessReasonError(
					common.ErrMPreProcess.Wrap(common.ErrMAccountNAth).
						Errorf("sender, %v is not registered as voter, %v, %v", fact.Sender(), fact.Contract(), fact.ProposalID())), nil
			}
		}
	}

	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("voting power box, %v, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	case found:
		vpb, err := state.StateVotingPowerBoxValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMStateValInvalid).
					Errorf("voting power box, %v, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
		}

		vp, found := vpb.VotingPowers()[fact.Sender().String()]
		if !found {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
					Errorf("sender, $v has no voting power, %v, %v", fact.Contract(), fact.ProposalID())), nil
		}

		if vp.Voted() {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
					Errorf("sender, $v already voted, %v, %v", fact.Contract(), fact.ProposalID())), nil
		}
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *VoteProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Vote")

	fact, ok := op.Fact().(VoteFact)
	if !ok {
		return nil, nil, e.Errorf("expected VoteFact, not %T", op.Fact())
	}

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "proposal", getStateFunc)
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

	period, start, end := types.GetPeriodOfCurrentTime(p.Policy(), p.Proposal(), types.Voting, blockMap)
	if period != types.Voting {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within Voting period, Voting period; start(%d), end(%d), but now(%d)", start, end, blockMap.Manifest().ProposedAt().Unix()), nil
	}

	var sts []base.StateMergeValue

	var votingPowerBox types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case !found:
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	default:
		vpb, err := state.StateVotingPowerBoxValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
		}
		votingPowerBox = vpb
	}

	vp, found := votingPowerBox.VotingPowers()[fact.Sender().String()]
	if !found {
		return nil, base.NewBaseOperationProcessReasonError("sender voting power not found, sender(%s), %s, %q", fact.Sender(), fact.Contract(), fact.ProposalID()), nil
	}
	vp.SetVoted(true)
	vp.SetVoteFor(fact.Vote())

	vpb := votingPowerBox.VotingPowers()
	vpb[fact.Sender().String()] = vp
	votingPowerBox.SetVotingPowers(vpb)

	result := votingPowerBox.Result()
	if _, found := result[fact.Vote()]; found {
		result[fact.Vote()] = result[fact.Vote()].Add(vp.Amount())
	} else {
		result[fact.Vote()] = common.ZeroBig.Add(vp.Amount())
	}
	votingPowerBox.SetResult(result)

	sts = append(sts,
		currencystate.NewStateMergeValue(
			state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID()),
			state.NewVotingPowerBoxStateValue(votingPowerBox),
		),
	)

	{ // caculate operation fee
		currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
		}

		if currencyPolicy.Feeer().Receiver() == nil {
			return sts, nil, nil
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

	return sts, nil, nil
}

func (opp *VoteProcessor) Close() error {
	voteProcessorPool.Put(opp)

	return nil
}
