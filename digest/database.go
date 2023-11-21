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
	defaultColNameAccount         = "digest_ac"
	defaultColNameContractAccount = "digest_ca"
	defaultColNameBalance         = "digest_bl"
	defaultColNameCurrency        = "digest_cr"
	defaultColNameOperation       = "digest_op"
	defaultColNameBlock           = "digest_bm"
	defaultColNameDAO             = "digest_dao_design"
	defaultColNameProposal        = "digest_dao_proposal"
	defaultColNameDelegators      = "digest_dao_delegators"
	defaultColNameVoters          = "digest_dao_voters"
	defaultColNameVotingPowerBox  = "digest_dao_voting_power_box"
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

func DelegatorInfo(st *currencydigest.Database, contract, proposalID, delegator string) (*types.DelegatorInfo, error) {
	var (
		delegators    []types.DelegatorInfo
		sta           mitumbase.State
		delegatorInfo *types.DelegatorInfo
		err           error
	)

	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	if err = st.DatabaseClient().GetByFilter(
		defaultColNameDelegators,
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

func Voters(st *currencydigest.Database, contract, proposalID string) ([]types.VoterInfo, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	var voters []types.VoterInfo
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameVoters,
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

func Proposal(st *currencydigest.Database, contract, proposalID string) (*state.ProposalStateValue, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	var proposal state.ProposalStateValue
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameProposal,
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

func VotingPowerBox(st *currencydigest.Database, contract, proposalID string) (*types.VotingPowerBox, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("proposal_id", proposalID)

	var votingPowerBox types.VotingPowerBox
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameVotingPowerBox,
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

//func CredentialsByServiceAndTemplate(
//	st *currencydigest.Database,
//	contract,
//	serviceID, templateID string,
//	reverse bool,
//	offset string,
//	limit int64,
//	callback func(nft types.Credential, st mitumbase.State) (bool, error),
//) error {
//	filter, err := buildCredentialFilterByService(contract, serviceID, templateID, offset, reverse)
//	if err != nil {
//		return err
//	}
//
//	sr := 1
//	if reverse {
//		sr = -1
//	}
//
//	opt := options.Find().SetSort(
//		util.NewBSONFilter("height", sr).D(),
//	)
//
//	switch {
//	case limit <= 0: // no limit
//	case limit > maxLimit:
//		opt = opt.SetLimit(maxLimit)
//	default:
//		opt = opt.SetLimit(limit)
//	}
//
//	return st.DatabaseClient().Find(
//		context.Background(),
//		defaultColNameDIDCredential,
//		filter,
//		func(cursor *mongo.Cursor) (bool, error) {
//			st, err := currencydigest.LoadState(cursor.Decode, st.DatabaseEncoders())
//			if err != nil {
//				return false, err
//			}
//			credential, err := state.StateCredentialValue(st)
//			if err != nil {
//				return false, err
//			}
//			return callback(credential, st)
//		},
//		opt,
//	)
//}
//
//func buildCredentialFilterByService(contract, col, templateID string, offset string, reverse bool) (bson.D, error) {
//	filterA := bson.A{}
//
//	// filter fot matching collection
//	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
//	filterSymbol := bson.D{{"service", bson.D{{"$in", []string{col}}}}}
//	filterTemplate := bson.D{{"template", bson.D{{"$in", []string{templateID}}}}}
//	filterA = append(filterA, filterSymbol)
//	filterA = append(filterA, filterContract)
//	filterA = append(filterA, filterTemplate)
//
//	// if offset exist, apply offset
//	if len(offset) > 0 {
//		if !reverse {
//			filterOffset := bson.D{
//				{"credential_id", bson.D{{"$gt", offset}}},
//			}
//			filterA = append(filterA, filterOffset)
//			// if reverse true, lesser then offset height
//		} else {
//			filterHeight := bson.D{
//				{"credential_id", bson.D{{"$lt", offset}}},
//			}
//			filterA = append(filterA, filterHeight)
//		}
//	}
//
//	filter := bson.D{}
//	if len(filterA) > 0 {
//		filter = bson.D{
//			{"$and", filterA},
//		}
//	}
//
//	return filter, nil
//}
