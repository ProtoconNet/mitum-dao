package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (r VoterInfo) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      r.Hint().String(),
			"account":    r.account,
			"delegators": r.delegators,
		},
	)
}

type VoterInfoBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Account    string   `bson:"account"`
	Delegators []string `bson:"delegators"`
}

func (r *VoterInfo) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of VoterInfo")

	var u VoterInfoBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	r.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		r.account = a
	}

	acc := make([]base.Address, len(u.Delegators))
	for i, ba := range u.Delegators {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e.Wrap(err)
		}
		acc[i] = ac

	}
	r.delegators = acc

	return nil
}
