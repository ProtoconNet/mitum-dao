package digest

import (
	"github.com/ProtoconNet/mitum-dao/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareDAO() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var daoDesignModels []mongo.WriteModel
	var daoProposalModels []mongo.WriteModel
	var daoDelegatorsModels []mongo.WriteModel
	var daoVotersModels []mongo.WriteModel
	var daoVotingPowerBoxModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case state.IsStateDesignKey(st.Key()):
			j, err := bs.handleDAODesignState(st)
			if err != nil {
				return err
			}
			daoDesignModels = append(daoDesignModels, j...)
		case state.IsStateProposalKey(st.Key()):
			j, err := bs.handleDAOProposalState(st)
			if err != nil {
				return err
			}
			daoProposalModels = append(daoProposalModels, j...)
		case state.IsStateDelegatorsKey(st.Key()):
			j, err := bs.handleDAODelegatorsState(st)
			if err != nil {
				return err
			}
			daoDelegatorsModels = append(daoDelegatorsModels, j...)
		case state.IsStateVotersKey(st.Key()):
			j, err := bs.handleDAOVotersState(st)
			if err != nil {
				return err
			}
			daoVotersModels = append(daoVotersModels, j...)
		case state.IsStateVotingPowerBoxKey(st.Key()):
			j, err := bs.handleDAOVotingPowerBoxState(st)
			if err != nil {
				return err
			}
			daoVotingPowerBoxModels = append(daoVotingPowerBoxModels, j...)
		default:
			continue
		}
	}

	bs.daoDesignModels = daoDesignModels
	bs.daoProposalModels = daoProposalModels
	bs.daoDelegatorsModels = daoDelegatorsModels
	bs.daoVotersModels = daoVotersModels
	bs.daoVotingPowerBoxModels = daoVotingPowerBoxModels

	return nil
}

func (bs *BlockSession) handleDAODesignState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if designDoc, err := NewDAODesignDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(designDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDAOProposalState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftCollectionDoc, err := NewDAOProposalDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftCollectionDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDAODelegatorsState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if delegatorsDoc, err := NewDAODelegatorsDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(delegatorsDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDAOVotersState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if votersDoc, err := NewDAOVotersDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(votersDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDAOVotingPowerBoxState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftLastIndexDoc, err := NewDAOVotingPowerBoxDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftLastIndexDoc),
		}, nil
	}
}
