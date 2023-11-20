package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (vp VotingPower) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":        vp.Hint().String(),
			"account":      vp.account,
			"voted":        vp.voted,
			"vote_for":     vp.voteFor,
			"voting_power": vp.amount,
		},
	)
}

type VotingPowerBSONUnmarshaler struct {
	Hint        string `bson:"_hint"`
	Account     string `bson:"account"`
	Voted       bool   `bson:"voted"`
	VoteFor     uint8  `bson:"vote_for"`
	VotingPower string `bson:"voting_power"`
}

func (vp *VotingPower) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of VotingPower")

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
	vp.voteFor = u.VoteFor

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
	Hint         string           `bson:"_hint"`
	Total        string           `bson:"total"`
	VotingPowers bson.Raw         `bson:"voting_powers"`
	Result       map[uint8]string `bson:"result"`
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

	vp.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(u.Total)
	if err != nil {
		return e.Wrap(err)
	}
	vp.total = big

	votingPowers := make(map[string]VotingPower)
	m, err := enc.DecodeMap(u.VotingPowers)
	if err != nil {
		return e.Wrap(err)
	}
	for k, v := range m {
		vp, ok := v.(VotingPower)
		if !ok {
			return e.Wrap(errors.Errorf("expected VotingPower, not %T", v))
		}

		if _, err := base.DecodeAddress(k, enc); err != nil {
			return e.Wrap(err)
		}

		votingPowers[k] = vp
	}
	vp.votingPowers = votingPowers

	result := make(map[uint8]common.Big)
	for k, v := range u.Result {

		big, err := common.NewBigFromString(v)
		if err != nil {
			return e.Wrap(err)
		}

		result[k] = big
	}
	vp.result = result

	return nil
}
