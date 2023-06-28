package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (sh SnapHistory) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":     sh.Hint().String(),
			"timestamp": sh.timestamp,
			"snaps":     sh.snaps,
		},
	)
}

type SnapHistoryBSONUnmarshaler struct {
	Hint      string   `bson:"_hint"`
	TimeStamp uint64   `bson:"timestamp"`
	Snaps     bson.Raw `bson:"snaps"`
}

func (sh *SnapHistory) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode bson of SnapHistory")

	var u SnapHistoryBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e(err, "")
	}

	sh.BaseHinter = hint.NewBaseHinter(ht)
	sh.timestamp = u.TimeStamp

	hs, err := enc.DecodeSlice(u.Snaps)
	if err != nil {
		return e(err, "")
	}

	snaps := make([]VotingPower, len(hs))
	for i := range hs {
		s, ok := hs[i].(VotingPower)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected VotingPower, not %T", hs[i]), "")
		}

		snaps[i] = s
	}
	sh.snaps = snaps

	return nil
}
