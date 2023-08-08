package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (de DesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  de.Hint().String(),
			"design": de.design,
		},
	)
}

type DesignStateValueBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Design bson.Raw `bson:"design"`
}

func (de *DesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of DesignStateValue")

	var u DesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	de.BaseHinter = hint.NewBaseHinter(ht)

	var design types.Design
	if err := design.DecodeBSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}

	de.design = design

	return nil
}

func (p ProposalStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    p.Hint().String(),
			"status":   p.status,
			"proposal": p.proposal,
			"policy":   p.policy,
		},
	)
}

type ProposalStateValueBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Status   uint8    `bson:"status"`
	Proposal bson.Raw `bson:"proposal"`
	Policy   bson.Raw `bson:"policy"`
}

func (p *ProposalStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of ProposalStateValue")

	var u ProposalStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	p.BaseHinter = hint.NewBaseHinter(ht)

	if hinter, err := enc.Decode(u.Proposal); err != nil {
		return e.Wrap(err)
	} else if pr, ok := hinter.(types.Proposal); !ok {
		return e.Wrap(errors.Errorf("expected Proposal, not %T", hinter))
	} else {
		p.proposal = pr
	}

	if hinter, err := enc.Decode(u.Policy); err != nil {
		return e.Wrap(err)
	} else if po, ok := hinter.(types.Policy); !ok {
		return e.Wrap(errors.Errorf("expected Policy, not %T", hinter))
	} else {
		p.policy = po
	}

	p.status = types.ProposalStatus(types.Option(u.Status))

	return nil
}

func (dg DelegatorsStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      dg.Hint().String(),
			"delegators": dg.delegators,
		},
	)
}

type DelegatorsStateValueBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Delegators bson.Raw `bson:"delegators"`
}

func (dg *DelegatorsStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of DelegatorsStateValue")

	var u DelegatorsStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	dg.BaseHinter = hint.NewBaseHinter(ht)

	hr, err := enc.DecodeSlice(u.Delegators)
	if err != nil {
		return err
	}

	dgs := make([]types.DelegatorInfo, len(hr))
	for i, hinter := range hr {
		if v, ok := hinter.(types.DelegatorInfo); !ok {
			return e.Wrap(errors.Errorf("expected types.DelegatorInfo, not %T", hinter))
		} else {
			dgs[i] = v
		}
	}
	dg.delegators = dgs

	return nil
}

func (vt VotersStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     vt.Hint().String(),
			"registers": vt.voters,
		},
	)
}

type VotersStateValueBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	Registers bson.Raw `bson:"registers"`
}

func (vt *VotersStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of VotersStateValue")

	var u VotersStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	vt.BaseHinter = hint.NewBaseHinter(ht)

	hr, err := enc.DecodeSlice(u.Registers)
	if err != nil {
		return e.Wrap(err)
	}

	infos := make([]types.VoterInfo, len(hr))
	for i, hinter := range hr {
		rg, ok := hinter.(types.VoterInfo)
		if !ok {
			return e.Wrap(errors.Errorf("expected types.VoterInfo, not %T", hinter))
		}

		infos[i] = rg
	}
	vt.voters = infos

	return nil
}

//func (sh SnapHistoriesStateValue) MarshalBSON() ([]byte, error) {
//	return bsonenc.Marshal(
//		bson.M{
//			"_hint":     sh.Hint().String(),
//			"histories": sh.Histories,
//		},
//	)
//}
//
//type SnapHistoriesStateValueBSONUnmarshaler struct {
//	Hint      string   `bson:"_hint"`
//	Histories bson.Raw `bson:"histories"`
//}
//
//func (sh *SnapHistoriesStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
//	e := util.StringError("failed to decode bson of SnapHistoriesStateValue")
//
//	var u SnapHistoriesStateValueBSONUnmarshaler
//	if err := enc.Unmarshal(b, &u); err != nil {
//		return e.Wrap(err)
//	}
//
//	ht, err := hint.ParseHint(u.Hint)
//	if err != nil {
//		return e.Wrap(err)
//	}
//
//	sh.BaseHinter = hint.NewBaseHinter(ht)
//
//	hs, err := enc.DecodeSlice(u.Histories)
//	if err != nil {
//		return e.Wrap(err)
//	}
//
//	histories := make([]types.SnapHistory, len(hs))
//	for i, hinter := range hs {
//		h, ok := hinter.(types.SnapHistory)
//		if !ok {
//			return e.Wrap(errors.Errorf("expected types.SnapHistory, not %T", hinter))
//		}
//
//		histories[i] = h
//	}
//	sh.Histories = histories
//
//	return nil
//}

func (vb VotingPowerBoxStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":            vb.Hint().String(),
			"voting_power_box": vb.votingPowerBox,
		},
	)
}

type VotesStateValueBSONUnmarshaler struct {
	Hint           string   `bson:"_hint"`
	VotingPowerBox bson.Raw `bson:"voting_power_box"`
}

func (vb *VotingPowerBoxStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of VotingPowerBoxStateValue")

	var u VotesStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	vb.BaseHinter = hint.NewBaseHinter(ht)

	if hinter, err := enc.Decode(u.VotingPowerBox); err != nil {
		return e.Wrap(err)
	} else if v, ok := hinter.(types.VotingPowerBox); !ok {
		return e.Wrap(errors.Errorf("expected VotingPowerBox, not %T", hinter))
	} else {
		vb.votingPowerBox = v
	}

	return nil
}
