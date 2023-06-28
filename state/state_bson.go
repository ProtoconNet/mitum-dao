package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
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
	Active   bool     `bson:"active"`
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
	p.Active = u.Active

	if hinter, err := enc.Decode(u.Proposal); err != nil {
		return e(err, "")
	} else if pr, ok := hinter.(types.Proposal); !ok {
		return e(util.ErrWrongType.Errorf("expected Proposal, not %T", hinter), "")
	} else {
		p.Proposal = pr
	}

	return nil
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

func (r RegisterInfo) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       r.Hint().String(),
			"account":     r.Account,
			"approved_by": r.ApprovedBy,
		},
	)
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
		r.account = a
	}

	acc := make([]base.Address, len(u.ApprovedBy))
	for i, ba := range u.ApprovedBy {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e(err, "")
		}
		acc[i] = ac

	}
	r.approvedBy = acc

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

	hr, err := enc.DecodeSlice(u.Registers)
	if err != nil {
		return e(err, "")
	}

	infos := make([]RegisterInfo, len(hr))
	for i, hinter := range hr {
		rg, ok := hinter.(RegisterInfo)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected RegisterInfo, not %T", hinter), "")
		}

		infos[i] = rg
	}
	r.Registers = infos

	return nil
}

func (vp VotingPower) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        vp.Hint().String(),
			"account":      vp.account,
			"voting_power": vp.votingPower,
		},
	)
}

type VotingPowerBSONUnmarshaler struct {
	Hint        string `bson:"_hint"`
	Account     string `bson:"account"`
	VotingPower string `bson:"voting_power"`
}

func (vp *VotingPower) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of VotingPower")

	var u VotingPowerBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	vp.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e(err, "")
	default:
		vp.account = a
	}

	big, err := common.NewBigFromString(u.VotingPower)
	if err != nil {
		return e(err, "")
	}
	vp.votingPower = big

	return nil
}

func (sh SnapHistory) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     sh.Hint().String(),
			"timestamp": sh.timestamp,
			"snaps":     sh.snaps,
		},
	)
}

type SnapHistoryBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	TimeStamp uint64   `bson:"timestamp"`
	Snaps     bson.Raw `bson:"snaps"`
}

func (sh *SnapHistory) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of SnapHistory")

	var u SnapHistoryBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	sh.BaseHinter = hint.NewBaseHinter(ht)
	sh.timestamp = u.TimeStamp

	hs, err := enc.DecodeSlice(u.Snaps)
	if err != nil {
		return e(err, "")
	}

	snaps := make([]VotingPower, len(hs))
	for i := range hs {
		s, ok := hs[i].(VotingPower)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected VotingPower, not %T", hs[i]), "")
		}

		snaps[i] = s
	}
	sh.snaps = snaps

	return nil
}

func (sh SnapHistoriesStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     sh.Hint().String(),
			"histories": sh.Histories,
		},
	)
}

type SnapHistoriesStateValueBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	Histories bson.Raw `bson:"histories"`
}

func (sh *SnapHistoriesStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of SnapHistoriesStateValue")

	var u SnapHistoriesStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	sh.BaseHinter = hint.NewBaseHinter(ht)

	hs, err := enc.DecodeSlice(u.Histories)
	if err != nil {
		return e(err, "")
	}

	histories := make([]SnapHistory, len(hs))
	for i, hinter := range hs {
		h, ok := hinter.(SnapHistory)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected SnapHistory, not %T", hinter), "")
		}

		histories[i] = h
	}
	sh.Histories = histories

	return nil
}

func (v VotingPowers) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":         v.Hint().String(),
			"total":         v.total.String(),
			"voting_powers": v.votingPowers,
		},
	)
}

type VotingPowersBSONUnmarshaler struct {
	Hint         string   `bson:"_hint"`
	Total        string   `bson:"total"`
	VotingPowers bson.Raw `bson:"voting_powers"`
}

func (v *VotingPowers) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of VotingPowers")

	var u VotingPowersBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	v.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(u.Total)
	if err != nil {
		return e(err, "")
	}
	v.total = big

	hv, err := enc.DecodeSlice(u.VotingPowers)
	if err != nil {
		return e(err, "")
	}

	vps := make([]VotingPower, len(hv))
	for i, hinter := range hv {
		vp, ok := hinter.(VotingPower)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected VotingPower, not %T", hinter), "")
		}

		vps[i] = vp
	}
	v.votingPowers = vps

	return nil
}

func (v VotesStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  v.Hint().String(),
			"active": v.Active,
			"result": v.Result,
			"votes":  v.Votes,
		},
	)
}

type VotesStateValueBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Active bool     `bson:"active"`
	Result uint8    `bson:"results"`
	Votes  bson.Raw `bson:"votes"`
}

func (v *VotesStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of VotesStateValue")

	var u VotesStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	v.BaseHinter = hint.NewBaseHinter(ht)
	v.Active = u.Active
	v.Result = u.Result

	hvs, err := enc.DecodeSlice(u.Votes)
	if err != nil {
		return e(err, "")
	}

	votes := make([]VotingPowers, len(hvs))
	for i, hinter := range hvs {
		c, ok := hinter.(VotingPowers)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected VotingPowers, not %T", hinter), "")
		}

		votes[i] = c
	}
	v.Votes = votes

	return nil
}
