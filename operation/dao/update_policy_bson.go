package dao

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact UpdatePolicyFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":                  fact.Hint().String(),
			"sender":                 fact.sender,
			"contract":               fact.contract,
			"option":                 fact.option,
			"voting_power_token":     fact.votingPowerToken,
			"threshold":              fact.threshold,
			"fee":                    fact.fee,
			"whitelist":              fact.whitelist,
			"proposal_review_period": fact.proposalReviewPeriod,
			"registration_period":    fact.registrationPeriod,
			"pre_snapshot_period":    fact.preSnapshotPeriod,
			"voting_period":          fact.votingPeriod,
			"post_snapshot_period":   fact.postSnapshotPeriod,
			"execution_delay_period": fact.executionDelayPeriod,
			"turnout":                fact.turnout,
			"quorum":                 fact.quorum,
			"currency":               fact.currency,
			"hash":                   fact.BaseFact.Hash().String(),
			"token":                  fact.BaseFact.Token(),
		},
	)
}

type UpdatePolicyFactBSONUnmarshaler struct {
	Hint                 string   `bson:"_hint"`
	Sender               string   `bson:"sender"`
	Contract             string   `bson:"contract"`
	Option               string   `bson:"option"`
	VotingPowerToken     string   `bson:"voting_power_token"`
	Threshold            string   `bson:"threshold"`
	Fee                  bson.Raw `bson:"fee"`
	Whitelist            bson.Raw `bson:"whitelist"`
	ProposalReviewPeriod uint64   `bson:"proposal_review_period"`
	RegistrationPeriod   uint64   `bson:"registration_period"`
	PreSnapshotPeriod    uint64   `bson:"pre_snapshot_period"`
	VotingPeriod         uint64   `bson:"voting_period"`
	PostSnapshotPeriod   uint64   `bson:"post_snapshot_period"`
	ExecutionDelayPeriod uint64   `bson:"execution_delay_period"`
	Turnout              uint     `bson:"turnout"`
	Quorum               uint     `bson:"quorum"`
	Currency             string   `bson:"currency"`
}

func (fact *UpdatePolicyFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of UpdatePolicyFact")

	var ubf common.BaseFactBSONUnmarshaler

	if err := enc.Unmarshal(b, &ubf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(ubf.Hash))
	fact.BaseFact.SetToken(ubf.Token)

	var uf UpdatePolicyFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unpack(enc,
		uf.Sender,
		uf.Contract,
		uf.Option,
		uf.VotingPowerToken,
		uf.Threshold,
		uf.Fee,
		uf.Whitelist,
		uf.ProposalReviewPeriod,
		uf.RegistrationPeriod,
		uf.PreSnapshotPeriod,
		uf.VotingPeriod,
		uf.PostSnapshotPeriod,
		uf.ExecutionDelayPeriod,
		uf.Turnout,
		uf.Quorum,
		uf.Currency,
	)
}

func (op UpdatePolicy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *UpdatePolicy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of UpdatePolicy")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
