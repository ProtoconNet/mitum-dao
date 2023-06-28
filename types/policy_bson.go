package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (wl Whitelist) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    wl.Hint().String(),
			"active":   wl.active,
			"accounts": wl.accounts,
		},
	)
}

type WhitelistBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Active   bool     `bson:"active"`
	Accounts bson.Raw `bson:"accounts"`
}

func (wl *Whitelist) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Whitelist")

	var uw WhitelistBSONUnmarshaler
	if err := enc.Unmarshal(b, &uw); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uw.Hint)
	if err != nil {
		return e(err, "")
	}

	return wl.unpack(enc, ht, uw.Active, uw.Accounts)
}

func (po Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      po.Hint().String(),
			"token":      po.token,
			"threshold":  po.threshold,
			"fee":        po.fee,
			"whitelist":  po.whitelist,
			"delaytime":  po.delaytime,
			"snaptime":   po.snaptime,
			"voteperiod": po.voteperiod,
			"timelock":   po.timelock,
			"turnout":    po.turnout,
			"quorum":     po.quorum,
		},
	)
}

type PolicyBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Token      string   `bson:"token"`
	Threshold  bson.Raw `bson:"threshold"`
	Fee        bson.Raw `bson:"fee"`
	Whitelist  bson.Raw `bson:"whitelist"`
	Delaytime  uint64   `bson:"delaytime"`
	Snaptime   uint64   `bson:"snaptime"`
	VotePeriod uint64   `bson:"voteperiod"`
	Timelock   uint64   `bson:"timelock"`
	Turnout    uint     `bson:"turnout"`
	Quorum     uint     `bson:"quorum"`
}

func (po *Policy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Policy")

	var upo PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(upo.Hint)
	if err != nil {
		return e(err, "")
	}

	return po.unpack(enc, ht,
		upo.Token,
		upo.Threshold,
		upo.Fee,
		upo.Whitelist,
		upo.Delaytime,
		upo.Snaptime,
		upo.VotePeriod,
		upo.Timelock,
		upo.Turnout,
		upo.Quorum,
	)
}
