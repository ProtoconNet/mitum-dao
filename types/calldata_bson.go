package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (cd TransferCallData) MarshalBSON() ([]byte, error) {
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

func (cd *TransferCallData) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of TransferCallData")

	var uc TransferCalldataBSONUnmarshaler
	if err := enc.Unmarshal(b, &uc); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uc.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return cd.unpack(enc, ht, uc.Sender, uc.Receiver, uc.Amount)
}

func (cd GovernanceCallData) MarshalBSON() ([]byte, error) {
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

func (cd *GovernanceCallData) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of GovernanceCallData")

	var uc GovernanceCalldataBSONUnmarshaler
	if err := enc.Unmarshal(b, &uc); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uc.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return cd.unpack(enc, ht, uc.Policy)
}
