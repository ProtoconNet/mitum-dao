package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (cd TransferCalldata) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    cd.Hint().String(),
			"sender":   cd.sender,
			"receiver": cd.receiver,
			"amount":   cd.amount,
		},
	)
}

type TransferCalldataBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Sender   string   `bson:"sender"`
	Receiver string   `bson:"receiver"`
	Amount   bson.Raw `bson:"amount"`
}

func (cd *TransferCalldata) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of TransferCalldata")

	var uc TransferCalldataBSONUnmarshaler
	if err := enc.Unmarshal(b, &uc); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uc.Hint)
	if err != nil {
		return e(err, "")
	}

	return cd.unpack(enc, ht, uc.Sender, uc.Receiver, uc.Amount)
}

func (cd GovernanceCalldata) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  cd.Hint().String(),
			"policy": cd.policy,
		},
	)
}

type GovernanceCalldataBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Policy bson.Raw `bson:"policy"`
}

func (cd *GovernanceCalldata) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of GovernanceCalldata")

	var uc GovernanceCalldataBSONUnmarshaler
	if err := enc.Unmarshal(b, &uc); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uc.Hint)
	if err != nil {
		return e(err, "")
	}

	return cd.unpack(enc, ht, uc.Policy)
}
