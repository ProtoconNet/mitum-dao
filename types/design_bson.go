package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  de.Hint().String(),
			"option": de.option,
			"dao_id": de.daoID,
			"policy": de.policy,
		},
	)
}

type DesignBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Option string   `bson:"option"`
	DAO    string   `bson:"dao_id"`
	Policy bson.Raw `bson:"policy"`
}

func (de *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Design")

	var ud DesignBSONUnmarshaler
	if err := enc.Unmarshal(b, &ud); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(ud.Hint)
	if err != nil {
		return e(err, "")
	}

	return de.unpack(enc, ht, ud.Option, ud.DAO, ud.Policy)
}
