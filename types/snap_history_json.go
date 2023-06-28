package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type SnapHistoryJSONMarshaler struct {
	hint.BaseHinter
	TimeStamp uint64        `json:"timestamp"`
	Snaps     []VotingPower `json:"snaps"`
}

func (sh SnapHistory) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SnapHistoryJSONMarshaler{
		BaseHinter: sh.BaseHinter,
		TimeStamp:  sh.timestamp,
		Snaps:      sh.snaps,
	})
}

type SnapHistoryJSONUnmarshaler struct {
	TimeStamp uint64          `json:"timestamp"`
	Snaps     json.RawMessage `json:"snaps"`
}

func (sh *SnapHistory) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of SnapHistory")

	var u SnapHistoryJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

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
