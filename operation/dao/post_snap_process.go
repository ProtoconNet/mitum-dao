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

var postSnapProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(PostSnapProcessor)
	},
}

func (PostSnap) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type PostSnapProcessor struct {
	*base.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewPostSnapProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new PostSnapProcessor")

		nopp := postSnapProcessorPool.Get()
		opp, ok := nopp.(*PostSnapProcessor)
		if !ok {
			return nil, errors.Errorf("expected PostSnapProcessor, not %T", nopp)
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

func (opp *PostSnapProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(PostSnapFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", PostSnapFact{}, op.Fact())), nil
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
	} else if p.Status() == types.PostSnapped {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMValueInvalid).
				Errorf("already post snapped, %s, %q", fact.Contract(), fact.ProposalID())), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVoters(fact.Contract(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("voters, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMStateNF).
				Errorf("voting power box, %s, %v: %v", fact.Contract(), fact.ProposalID(), err)), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	return ctx, nil, nil
}

func (opp *PostSnapProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process PostSnap")

	fact, ok := op.Fact().(PostSnapFact)
	if !ok {
		return nil, nil, e.Errorf("expected PostSnapFact, not %T", op.Fact())
	}

	st, err := currencystate.ExistsState(state.StateKeyProposal(fact.Contract(), fact.ProposalID()), "proposal", getStateFunc)
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

	period, start, end := types.GetPeriodOfCurrentTime(p.Policy(), p.Proposal(), types.PostSnapshot, blockMap)
	if period != types.PostSnapshot {
		return nil, base.NewBaseOperationProcessReasonError("current time is not within the PostSnapshotPeriod, PostSnapshotPeriod; start(%d), end(%d), but now(%d)", start, end, blockMap.Manifest().ProposedAt().Unix()), nil
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

	if p.Status() != types.PreSnapped {
		sts = append(sts,
			currencystate.NewStateMergeValue(
				st.Key(),
				state.NewProposalStateValue(types.Canceled, p.Proposal(), p.Policy()),
			),
		)

		return sts, nil, nil
	}

	var ovpb types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		if vb, err := state.StateVotingPowerBoxValue(st); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
		} else {
			ovpb = vb
		}
	default:
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s, %q", fact.Contract(), fact.ProposalID()), nil
	}

	votingPowerToken := p.Policy().Token()

	var nvpb = types.NewVotingPowerBox(common.ZeroBig, map[string]types.VotingPower{})

	nvps := map[string]types.VotingPower{}
	nvt := common.ZeroBig

	votedTotal := common.ZeroBig
	votingResult := map[uint8]common.Big{}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
	case found:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value, %s, %q: %w", fact.Contract(), fact.ProposalID(), err), nil
		}

		for _, info := range voters {
			a := info.Account().String()

			if !ovpb.VotingPowers()[a].Voted() {
				nvps[a] = ovpb.VotingPowers()[a]
				continue
			}

			vp := common.ZeroBig
			for _, delegator := range info.Delegators() {
				st, err = currencystate.ExistsState(currency.StateKeyBalance(delegator, votingPowerToken), "key of balance", getStateFunc)
				if err != nil {
					continue
				}

				b, err := currency.StateBalanceValue(st)
				if err != nil {
					return nil, base.NewBaseOperationProcessReasonError("failed to find balance value of the delegator from state, %s, %q: %w", delegator, votingPowerToken, err), nil
				}

				vp = vp.Add(b.Big())
			}

			ovp := ovpb.VotingPowers()[a]
			if ovp.Amount().Compare(vp) < 0 {
				nvps[a] = ovp
			} else {
				nvp := types.NewVotingPower(info.Account(), vp)
				nvp.SetVoted(ovp.Voted())
				nvp.SetVoteFor(ovp.VoteFor())

				nvps[a] = nvp
			}

			nvt = nvt.Add(nvps[a].Amount())

			if nvps[a].Voted() {
				if _, found := votingResult[nvps[a].VoteFor()]; !found {
					votingResult[nvps[a].VoteFor()] = common.ZeroBig
				}
				votingResult[nvps[a].VoteFor()] = votingResult[nvps[a].VoteFor()].Add(vp)
				votedTotal = votedTotal.Add(nvps[a].Amount())
			}
		}

		nvpb.SetVotingPowers(nvps)
		nvpb.SetTotal(nvt)
		nvpb.SetResult(votingResult)
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyVotingPowerBox(fact.Contract(), fact.ProposalID()),
		state.NewVotingPowerBoxStateValue(nvpb),
	))

	st, err = currencystate.ExistsState(currency.StateKeyCurrencyDesign(votingPowerToken), "key of currency design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency state, %q: %w", votingPowerToken, err), nil
	}

	currencyDesign, err := currency.StateCurrencyDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency design value from state, %q: %w", votingPowerToken, err), nil
	}

	actualTurnoutCount := p.Policy().Turnout().Quorum(currencyDesign.Aggregate())
	actualQuorumCount := p.Policy().Quorum().Quorum(votedTotal)

	r := types.Rejected

	switch {
	case nvpb.Total().Compare(actualTurnoutCount) < 0:
		r = types.Canceled
	case votedTotal.Compare(actualQuorumCount) < 0:
	case p.Proposal().Option() == types.ProposalCrypto:
		vr0, found0 := votingResult[0]
		vr1, found1 := votingResult[1]
		if !found0 {
			r = types.Rejected
			break
		} else {
			if found1 {
				if (0 < vr0.Compare(actualQuorumCount)) && (0 < vr0.Compare(vr1)) {
					r = types.Completed
					break
				}
			} else {
				if 0 < vr0.Compare(actualQuorumCount) {
					r = types.Completed
					break
				}
			}
		}
	case p.Proposal().Option() == types.ProposalBiz:
		options := p.Proposal().VoteOptionsCount() - 1

		var count = 0
		var mvp = common.ZeroBig
		var i uint8 = 0

		for ; i < options; i++ {
			if votingResult[i].Compare(actualQuorumCount) >= 0 {
				if mvp.Compare(votingResult[i]) < 0 {
					count = 1
					mvp = votingResult[i]
				} else if mvp.Equal(votingResult[i]) {
					count += 1
				}
			}
		}

		if count == 1 {
			r = types.Completed
		}
	}

	//if nvpb.Total().Compare(actualTurnoutCount) < 0 {
	//	r = types.Canceled
	//} else if votedTotal.Compare(actualQuorumCount) < 0 {
	//	r = types.Rejected
	//} else if p.Proposal().Option() == types.ProposalCrypto {
	//	vr0, found0 := votingResult[0]
	//	vr1, found1 := votingResult[1]
	//
	//	if !(found0 && 0 < vr0.Compare(actualQuorumCount) && (!found1 || (found1 && 0 < vr0.Compare(vr1)))) {
	//		r = types.Rejected
	//	}
	//} else if p.Proposal().Option() == types.ProposalBiz {
	//	options := p.Proposal().VoteOptionsCount() - 1
	//
	//	var count = 0
	//	var mvp = common.ZeroBig
	//	var i uint8 = 0
	//
	//	for ; i < options; i++ {
	//		if votingResult[i].Compare(actualQuorumCount) >= 0 {
	//			if mvp.Compare(votingResult[i]) < 0 {
	//				count = 1
	//				mvp = votingResult[i]
	//			} else if mvp.Equal(votingResult[i]) {
	//				count += 1
	//			}
	//		}
	//	}
	//
	//	if count != 1 {
	//		r = types.Rejected
	//	}
	//}

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyProposal(fact.Contract(), fact.ProposalID()),
		state.NewProposalStateValue(r, p.Proposal(), p.Policy()),
	))

	return sts, nil, nil
}

func (opp *PostSnapProcessor) Close() error {
	postSnapProcessorPool.Put(opp)

	return nil
}
