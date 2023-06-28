package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

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
