package dao

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact SnapFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     fact.Hint().String(),
			"sender":    fact.sender,
			"contract":  fact.contract,
			"daoid":     fact.daoID,
			"proposeid": fact.proposeID,
			"currency":  fact.currency,
			"hash":      fact.BaseFact.Hash().String(),
			"token":     fact.BaseFact.Token(),
		},
	)
}

type SnapFactBSONUnmarshaler struct {
	Hint      string `bson:"_hint"`
	Sender    string `bson:"sender"`
	Contract  string `bson:"contract"`
	DAOID     string `bson:"daoid"`
	ProposeID string `bson:"proposeid"`
	Currency  string `bson:"currency"`
}

func (fact *SnapFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of SnapFact")

	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf SnapFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e(err, "")
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc,
		uf.Sender,
		uf.Contract,
		uf.DAOID,
		uf.ProposeID,
		uf.Currency,
	)
}

func (op Snap) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *Snap) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Snap")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
