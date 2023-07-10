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
	e := util.StringError("failed to preprocess PostSnap")

	fact, ok := op.Fact().(PostSnapFact)
	if !ok {
		return ctx, nil, e.Errorf("not PostSnapFact, %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(currency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	if err := currencystate.CheckNotExistsState(extensioncurrency.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("contract account cannot preSnap, %q: %w", fact.Sender(), err), nil
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
	} else if p.Status() == types.PostSnapped {
		return nil, base.NewBaseOperationProcessReasonError("already post snapped, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	blocMap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("get LastBlock failed: %w", err), nil
	} else if !found {
		return nil, base.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	period, start, end := types.GetPeriodOfCurrentTime(design.Policy(), p.Proposal(), blocMap)
	if period != types.PostSnapshot {
		return nil, base.NewBaseOperationProcessReasonError("currency time is not within the PostSnapshotPeriod, PostSnapshotPeriod start : %d, end %d", start, end), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voters state not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if err := currencystate.CheckExistsState(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
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

	if p.Status() != types.PreSnapped {
		sts = append(sts,
			currencystate.NewStateMergeValue(
				st.Key(),
				state.NewProposalStateValue(types.Canceled, p.Proposal()),
			),
		)

		return sts, nil, nil
	}

	var ovpb types.VotingPowerBox
	switch st, found, err := getStateFunc(state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		if vb, err := state.StateVotingPowerBoxValue(st); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voting power box value from state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		} else {
			ovpb = vb
		}
	default:
		return nil, base.NewBaseOperationProcessReasonError("voting power box state not found, %s-%s-%s", fact.Contract(), fact.DAOID(), fact.ProposalID()), nil
	}

	votingPowerToken := design.Policy().Token()

	var nvpb = types.NewVotingPowerBox(common.ZeroBig, map[base.Address]types.VotingPower{})

	nvps := map[base.Address]types.VotingPower{}
	nvt := common.ZeroBig

	votedTotal := common.ZeroBig
	votingResult := map[uint8]common.Big{}

	switch st, found, err := getStateFunc(state.StateKeyVoters(fact.Contract(), fact.DAOID(), fact.ProposalID())); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("failed to find voters state, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
	case found:
		voters, err := state.StateVotersValue(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError("failed to find voters value, %s-%s-%s: %w", fact.Contract(), fact.DAOID(), fact.ProposalID(), err), nil
		}

		for _, info := range voters {
			if !ovpb.VotingPowers()[info.Account()].Voted() {
				nvps[info.Account()] = ovpb.VotingPowers()[info.Account()]
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
					return nil, base.NewBaseOperationProcessReasonError("failed to find balance value of the delegator from state, %q, %q: %w", delegator, votingPowerToken, err), nil
				}

				vp = vp.Add(b.Big())
			}

			ovp := ovpb.VotingPowers()[info.Account()]
			if ovp.Amount().Compare(vp) < 0 {
				nvps[info.Account()] = ovp
			} else {
				nvp := types.NewVotingPower(info.Account(), vp)
				nvp.SetVoted(ovp.Voted())
				nvp.SetVoteFor(ovp.VoteFor())

				nvps[info.Account()] = nvp
			}

			nvt = nvt.Add(nvps[info.Account()].Amount())

			if nvps[info.Account()].Voted() {
				votedTotal = votedTotal.Add(nvps[info.Account()].Amount())

				if _, found := votingResult[nvps[info.Account()].VoteFor()]; !found {
					votingResult[nvps[info.Account()].VoteFor()] = common.ZeroBig
				}

				votingResult[nvps[info.Account()].VoteFor()] = votingResult[nvps[info.Account()].VoteFor()].Add(vp)
			}
		}

		nvpb.SetVotingPowers(nvps)
		nvpb.SetTotal(nvt)
		nvpb.SetResult(votingResult)
	}

	sts = append(sts, currencystate.NewStateMergeValue(
		state.StateKeyVotingPowerBox(fact.Contract(), fact.DAOID(), fact.ProposalID()),
		state.NewVotingPowerBoxStateValue(nvpb),
	))

	st, err = currencystate.ExistsState(currency.StateKeyCurrencyDesign(votingPowerToken), "key of currency design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency state, %s: %w", votingPowerToken, err), nil
	}

	currencyDesign, err := currency.StateCurrencyDesignValue(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to find voting power token currency design value from state, %s: %w", votingPowerToken, err), nil
	}

	actualTurnoutCount := design.Policy().Turnout().Quorum(currencyDesign.Aggregate())
	actualTurnoutQuorum := design.Policy().Quorum().Quorum(votedTotal)

	if nvpb.Total().Compare(actualTurnoutCount) < 0 {
		sts = append(sts, currencystate.NewStateMergeValue(
			state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
			state.NewProposalStateValue(types.Canceled, p.Proposal()),
		))
	} else if votedTotal.Compare(actualTurnoutQuorum) < 0 {
		sts = append(sts, currencystate.NewStateMergeValue(
			state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
			state.NewProposalStateValue(types.Rejected, p.Proposal()),
		))
	} else if p.Proposal().Type() == types.ProposalCrypto {
		if votingResult[0].Compare(actualTurnoutQuorum) >= 0 {
			sts = append(sts, currencystate.NewStateMergeValue(
				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
				state.NewProposalStateValue(types.Completed, p.Proposal()),
			))
		} else {
			sts = append(sts, currencystate.NewStateMergeValue(
				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
				state.NewProposalStateValue(types.Rejected, p.Proposal()),
			))
		}
	} else if p.Proposal().Type() == types.ProposalBiz {
		options := p.Proposal().Options()

		overQuorum := map[string][]uint8{}
		var maxVotingPower = common.ZeroBig
		var i uint8 = 0

		for {
			if i == options {
				break
			}

			if votingResult[i].Compare(actualTurnoutQuorum) >= 0 {
				if len(overQuorum) == 0 {
					overQuorum[votingResult[i].String()] = []uint8{i}
					maxVotingPower = votingResult[i]
					i++
					continue
				}

				overQuorum[votingResult[i].String()] = append(overQuorum[votingResult[i].String()], i)

				if votingResult[i].Compare(maxVotingPower) > 0 {
					maxVotingPower = votingResult[i]
				}
			}
			i++
		}

		if len(overQuorum[maxVotingPower.String()]) != 1 {
			sts = append(sts, currencystate.NewStateMergeValue(
				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
				state.NewProposalStateValue(types.Rejected, p.Proposal()),
			))
		} else {
			sts = append(sts, currencystate.NewStateMergeValue(
				state.StateKeyProposal(fact.Contract(), fact.DAOID(), fact.ProposalID()),
				state.NewProposalStateValue(types.Completed, p.Proposal()),
			))
		}
	}

	return sts, nil, nil
}

func (opp *PostSnapProcessor) Close() error {
	postSnapProcessorPool.Put(opp)

	return nil
}
