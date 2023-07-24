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
			j, err := bs.handleProposalState(st)
			if err != nil {
				return err
			}
			daoProposalModels = append(daoProposalModels, j...)
		case state.IsStateDelegatorsKey(st.Key()):
			j, err := bs.handleDelegatorsState(st)
			if err != nil {
				return err
			}
			daoDelegatorsModels = append(daoDelegatorsModels, j...)
		case state.IsStateVotersKey(st.Key()):
			j, err := bs.handleVotersState(st)
			if err != nil {
				return err
			}
			daoVotersModels = append(daoVotersModels, j...)
		case state.IsStateVotingPowerBoxKey(st.Key()):
			j, err := bs.handleVotingPowerBoxState(st)
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
	if designDoc, err := NewDesignDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(designDoc),
		}, nil
	}
}

func (bs *BlockSession) handleProposalState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftCollectionDoc, err := NewProposalDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftCollectionDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDelegatorsState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftDoc, err := NewDelegatorsDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftDoc),
		}, nil
	}
}

func (bs *BlockSession) handleVotersState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftBoxDoc, err := NewVotersDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftBoxDoc),
		}, nil
	}
}

func (bs *BlockSession) handleVotingPowerBoxState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftLastIndexDoc, err := NewVotingPowerBoxDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftLastIndexDoc),
		}, nil
	}
}
