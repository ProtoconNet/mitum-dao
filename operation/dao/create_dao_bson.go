package dao

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact CreateDAOFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":              fact.Hint().String(),
			"sender":             fact.sender,
			"contract":           fact.contract,
			"daoid":              fact.daoID,
			"option":             fact.option,
			"voting_power_token": fact.votingPowerToken,
			"threshold":          fact.threshold,
			"fee":                fact.fee,
			"whitelist":          fact.whitelist,
			"delaytime":          fact.delaytime,
			"snaptime":           fact.snaptime,
			"timelock":           fact.timelock,
			"turnout":            fact.turnout,
			"quorum":             fact.quorum,
			"currency":           fact.currency,
			"hash":               fact.BaseFact.Hash().String(),
			"token":              fact.BaseFact.Token(),
		},
	)
}

type CreateDAOFactBSONUnmarshaler struct {
	Hint             string   `bson:"_hint"`
	Sender           string   `bson:"sender"`
	Contract         string   `bson:"contract"`
	DAOID            string   `bson:"daoid"`
	Option           string   `bson:"option"`
	VotingPowerToken string   `bson:"voting_power_token"`
	Threshold        bson.Raw `bson:"threshold"`
	Fee              bson.Raw `bson:"fee"`
	Whitelist        bson.Raw `bson:"whitelist"`
	Delaytime        uint64   `bson:"delaytime"`
	Snaptime         uint64   `bson:"snaptime"`
	Timelock         uint64   `bson:"timelock"`
	Turnout          uint     `bson:"turnout"`
	Quorum           uint     `bson:"quorum"`
	Currency         string   `bson:"currency"`
}

func (fact *CreateDAOFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateDAOFact")

	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf CreateDAOFactBSONUnmarshaler
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
		uf.Option,
		uf.VotingPowerToken,
		uf.Threshold,
		uf.Fee,
		uf.Whitelist,
		uf.Delaytime,
		uf.Snaptime,
		uf.Timelock,
		uf.Turnout,
		uf.Quorum,
		uf.Currency,
	)
}

func (op CreateDAO) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *CreateDAO) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of CreateDAO")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
