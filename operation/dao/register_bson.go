package dao

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact RegisterFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":       fact.Hint().String(),
			"sender":      fact.sender,
			"contract":    fact.contract,
			"dao_id":      fact.daoID,
			"proposal_id": fact.proposalID,
			"delegated":   fact.delegated,
			"currency":    fact.currency,
			"hash":        fact.BaseFact.Hash().String(),
			"token":       fact.BaseFact.Token(),
		},
	)
}

type RegisterFactBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Sender     string `bson:"sender"`
	Contract   string `bson:"contract"`
	DAOID      string `bson:"dao_id"`
	ProposalID string `bson:"proposal_id"`
	Delegated  string `bson:"delegated"`
	Currency   string `bson:"currency"`
}

func (fact *RegisterFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of RegisterFact")

	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf RegisterFactBSONUnmarshaler
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
		uf.ProposalID,
		uf.Delegated,
		uf.Currency,
	)
}

func (op Register) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *Register) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of Register")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
