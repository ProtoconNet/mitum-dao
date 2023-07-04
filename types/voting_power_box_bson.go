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
			"voted":        vp.voted,
			"voting_power": vp.amount,
		},
	)
}

type VotingPowerBSONUnmarshaler struct {
	Hint        string `bson:"_hint"`
	Account     string `bson:"account"`
	Voted       bool   `bson:"voted"`
	VotingPower string `bson:"voting_power"`
}

func (vp *VotingPower) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Amount")

	var u VotingPowerBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	vp.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		vp.account = a
	}

	big, err := common.NewBigFromString(u.VotingPower)
	if err != nil {
		return e.Wrap(err)
	}
	vp.amount = big
	vp.voted = u.Voted

	return nil
}

func (vp VotingPowerBox) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":         vp.Hint().String(),
			"total":         vp.total.String(),
			"voting_powers": vp.votingPowers,
			"result":        vp.result,
		},
	)
}

type VotingPowerBoxBSONUnmarshaler struct {
	Hint         string   `bson:"_hint"`
	Total        string   `bson:"total"`
	VotingPowers bson.Raw `bson:"voting_powers"`
	Result       bson.Raw `bson:"result"`
}

func (vp *VotingPowerBox) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of VotingPowerBox")

	var u VotingPowerBoxBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return vp.unpack(enc, ht, u.Total, u.VotingPowers, u.Result)
}
