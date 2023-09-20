package digest

import (
	mongodbstorage "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type DesignDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	de types.Design
}

func NewDesignDoc(st base.State, enc encoder.Encoder) (DesignDoc, error) {
	de, err := state.StateDesignValue(st)
	if err != nil {
		return DesignDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DesignDoc{}, err
	}

	return DesignDoc{
		BaseDoc: b,
		st:      st,
		de:      de,
	}, nil
}

func (doc DesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.DAOPrefix)
	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()
	//m["design"] = doc.de

	return bsonenc.Marshal(m)
}

type ProposalDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	pr types.Proposal
	ps types.ProposalStatus
}

func NewProposalDoc(st base.State, enc encoder.Encoder) (ProposalDoc, error) {
	pv, err := state.StateProposalValue(st)
	if err != nil {
		return ProposalDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return ProposalDoc{}, err
	}

	return ProposalDoc{
		BaseDoc: b,
		st:      st,
		pr:      pv.Proposal(),
		ps:      pv.Status(),
	}, nil
}

func (doc ProposalDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.DAOPrefix)
	m["contract"] = parsedKey[1]
	m["proposal_id"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["proposal"] = doc.pr
	m["proposal_status"] = doc.ps

	return bsonenc.Marshal(m)
}

type DelegatorsDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	di []types.DelegatorInfo
}

func NewDelegatorsDoc(st base.State, enc encoder.Encoder) (DelegatorsDoc, error) {
	di, err := state.StateDelegatorsValue(st)
	if err != nil {
		return DelegatorsDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DelegatorsDoc{}, err
	}

	return DelegatorsDoc{
		BaseDoc: b,
		st:      st,
		di:      di,
	}, nil
}

func (doc DelegatorsDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.DAOPrefix)
	m["contract"] = parsedKey[1]
	m["proposal_id"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["delegators"] = doc.di

	return bsonenc.Marshal(m)
}

type VotersDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	vi []types.VoterInfo
}

func NewVotersDoc(st base.State, enc encoder.Encoder) (VotersDoc, error) {
	vi, err := state.StateVotersValue(st)
	if err != nil {
		return VotersDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return VotersDoc{}, err
	}

	return VotersDoc{
		BaseDoc: b,
		st:      st,
		vi:      vi,
	}, nil
}

func (doc VotersDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.DAOPrefix)
	m["contract"] = parsedKey[1]
	m["proposal_id"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["voters"] = doc.vi

	return bsonenc.Marshal(m)
}

type VotingPowerBoxDoc struct {
	mongodbstorage.BaseDoc
	st  base.State
	vpb types.VotingPowerBox
}

func NewVotingPowerBoxDoc(st base.State, enc encoder.Encoder) (VotingPowerBoxDoc, error) {
	vpb, err := state.StateVotingPowerBoxValue(st)
	if err != nil {
		return VotingPowerBoxDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return VotingPowerBoxDoc{}, err
	}

	return VotingPowerBoxDoc{
		BaseDoc: b,
		st:      st,
		vpb:     vpb,
	}, nil
}

func (doc VotingPowerBoxDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.DAOPrefix)
	m["contract"] = parsedKey[1]
	m["proposal_id"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["voting_power_box"] = doc.vpb

	return bsonenc.Marshal(m)
}
