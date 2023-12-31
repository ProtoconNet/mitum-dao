package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (p CryptoProposal) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      p.Hint().String(),
			"proposer":   p.proposer,
			"start_time": p.startTime,
			"call_data":  p.callData,
		},
	)
}

type CryptoProposalBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	Proposer  string   `bson:"proposer"`
	StartTime uint64   `bson:"start_time"`
	CallData  bson.Raw `bson:"call_data"`
}

func (p *CryptoProposal) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of CryptoProposal")

	var up CryptoProposalBSONUnmarshaler
	if err := enc.Unmarshal(b, &up); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(up.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return p.unpack(enc, ht, up.Proposer, up.StartTime, up.CallData)
}

func (p BizProposal) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      p.Hint().String(),
			"proposer":   p.proposer,
			"start_time": p.startTime,
			"url":        p.url,
			"hash":       p.hash,
			"options":    p.options,
		},
	)
}

type BizProposalBSONUnmarshaler struct {
	Hint      string `bson:"_hint"`
	Proposer  string `bson:"proposer"`
	StartTime uint64 `bson:"start_time"`
	Url       string `bson:"url"`
	Hash      string `bson:"hash"`
	Options   uint8  `bson:"options"`
}

func (p *BizProposal) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of BizProposal")

	var up BizProposalBSONUnmarshaler
	if err := enc.Unmarshal(b, &up); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(up.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return p.unpack(enc, ht, up.Proposer, up.StartTime, up.Url, up.Hash, up.Options)
}
