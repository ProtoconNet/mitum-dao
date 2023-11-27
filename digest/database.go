package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	defaultColNameAccount           = "digest_ac"
	defaultColNameContractAccount   = "digest_ca"
	defaultColNameBalance           = "digest_bl"
	defaultColNameCurrency          = "digest_cr"
	defaultColNameOperation         = "digest_op"
	defaultColNameBlock             = "digest_bm"
	defaultColNameDAO               = "digest_dao_de"
	defaultColNameDAOProposal       = "digest_dao_pr"
	defaultColNameDAODelegators     = "digest_dao_dac"
	defaultColNameDAOVoters         = "digest_dao_vac"
	defaultColNameDAOVotingPowerBox = "digest_dao_vpb"
)

func DAOService(st *currencydigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design types.Design
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameDAO,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			design, err = state.StateDesignValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return &design, nil
}

func DAODelegatorInfo(st *currencydigest.Database, contract, proposalID, delegator string) (*types.DelegatorInfo, error) {
	var (
		delegators    []types.DelegatorInfo
		sta           mitumbase.State
		delegatorInfo *types.DelegatorInfo
		err           error
	)

	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	if err = st.DatabaseClient().GetByFilter(
		defaultColNameDAODelegators,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			delegators, err = state.StateDelegatorsValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	for i := range delegators {
		if delegator == delegators[i].Account().String() {
			delegatorInfo = &delegators[i]
			break
		}
	}
	if delegatorInfo == nil {
		return nil, errors.Errorf("delegator not found, %s", delegator)
	}

	return delegatorInfo, nil
}

func DAOVoters(st *currencydigest.Database, contract, proposalID string) ([]types.VoterInfo, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	var voters []types.VoterInfo
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameDAOVoters,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			voters, err = state.StateVotersValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return voters, nil
}

func DAOProposal(st *currencydigest.Database, contract, proposalID string) (*state.ProposalStateValue, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	var proposal state.ProposalStateValue
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameDAOProposal,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			proposal, err = state.StateProposalValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return &proposal, nil
}

func DAOVotingPowerBox(st *currencydigest.Database, contract, proposalID string) (*types.VotingPowerBox, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	var votingPowerBox types.VotingPowerBox
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameDAOVotingPowerBox,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			votingPowerBox, err = state.StateVotingPowerBoxValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return &votingPowerBox, nil
}
