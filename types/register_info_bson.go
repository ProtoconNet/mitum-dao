package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

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
