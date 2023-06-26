package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (de DesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": de.Hint().String(),
			"dao":   de.Design,
		},
	)
}

type DesignStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	DAO  bson.Raw `bson:"dao"`
}

func (de *DesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of DesignStateValue")

	var u DesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	de.BaseHinter = hint.NewBaseHinter(ht)

	var design types.Design
	if err := design.DecodeBSON(u.DAO, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}

func (p ProposalStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    p.Hint().String(),
			"proposal": p.Proposal,
		},
	)
}

type ProposalStateValueBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Proposal bson.Raw `bson:"proposal"`
}

func (p *ProposalStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of ProposalStateValue")

	var u ProposalStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	p.BaseHinter = hint.NewBaseHinter(ht)

	if hinter, err := enc.Decode(u.Proposal); err != nil {
		return e(err, "")
	} else if pr, ok := hinter.(types.Proposal); !ok {
		return e(util.ErrWrongType.Errorf("expected Proposal, not %T", hinter), "")
	} else {
		p.Proposal = pr
	}

	return nil
}

func (r RegisterInfo) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       r.Hint().String(),
			"account":     r.Account,
			"approved_by": r.ApprovedBy,
		},
	)
}

func (ap ApprovingListStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    ap.Hint().String(),
			"accounts": ap.Accounts,
		},
	)
}

type ApprovingListStateValueBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Accounts []string `bson:"accounts"`
}

func (ap *ApprovingListStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of ApprovingStateValue")

	var u ApprovingListStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	ap.BaseHinter = hint.NewBaseHinter(ht)

	acc := make([]base.Address, len(u.Accounts))
	for i, ba := range u.Accounts {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e(err, "")
		}
		acc[i] = ac

	}
	ap.Accounts = acc

	return nil
}

type RegisterInfoBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Account    string   `bson:"account"`
	ApprovedBy []string `bson:"approved_by"`
}

func (r *RegisterInfo) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of RegisterInfo")

	var u RegisterInfoBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	r.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e(err, "")
	default:
		r.Account = a
	}

	acc := make([]base.Address, len(u.ApprovedBy))
	for i, ba := range u.ApprovedBy {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e(err, "")
		}
		acc[i] = ac

	}
	r.ApprovedBy = acc

	return nil
}

func (r RegisterListStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     r.Hint().String(),
			"registers": r.Registers,
		},
	)
}

type RegisterListStateValueBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	Registers bson.Raw `bson:"registers"`
}

func (r *RegisterListStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of RegisterListStateValue")

	var u RegisterListStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	r.BaseHinter = hint.NewBaseHinter(ht)

	hit, err := enc.DecodeSlice(u.Registers)
	if err != nil {
		return e(err, "")
	}

	rs := make([]RegisterInfo, len(hit))
	for i, hinter := range hit {
		rg, ok := hinter.(RegisterInfo)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected RegisterInfo, not %T", hinter), "")
		}

		rs[i] = rg
	}
	r.Registers = rs

	return nil
}
