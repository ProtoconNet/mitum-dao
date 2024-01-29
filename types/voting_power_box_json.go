package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type VotingPowerJSONMarshaler struct {
	hint.BaseHinter
	Account     base.Address `json:"account"`
	Voted       bool         `json:"voted"`
	VoteFor     uint8        `json:"vote_for"`
	VotingPower string       `json:"voting_power"`
}

func (vp VotingPower) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotingPowerJSONMarshaler{
		BaseHinter:  vp.BaseHinter,
		Account:     vp.account,
		Voted:       vp.voted,
		VoteFor:     vp.voteFor,
		VotingPower: vp.amount.String(),
	})
}

type VotingPowerJSONUnmarshaler struct {
	Account     string `json:"account"`
	Voted       bool   `json:"voted"`
	VoteFor     uint8  `json:"vote_for"`
	VotingPower string `json:"voting_power"`
}

func (vp *VotingPower) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of VotingPower")

	var u VotingPowerJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		vp.account = a
	}

	big, err := common.NewBigFromString(u.VotingPower)
	if err != nil {
		return e.Wrap(err)
	}
	vp.amount = big
	vp.voted = u.Voted
	vp.voteFor = u.VoteFor

	return nil
}

type VotingPowerBoxJSONMarshaler struct {
	hint.BaseHinter
	Total        string                 `json:"total"`
	VotingPowers map[string]VotingPower `json:"voting_powers"`
	Result       map[uint8]common.Big   `json:"result"`
}

func (vp VotingPowerBox) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotingPowerBoxJSONMarshaler{
		BaseHinter:   vp.BaseHinter,
		Total:        vp.total.String(),
		VotingPowers: vp.votingPowers,
		Result:       vp.result,
	})
}

type VotingPowerBoxJSONUnmarshaler struct {
	Hint         hint.Hint       `json:"_hint"`
	Total        string          `json:"total"`
	VotingPowers json.RawMessage `json:"voting_powers"`
	Result       json.RawMessage `json:"result"`
}

func (vp *VotingPowerBox) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of VotingPowerBox")

	var u VotingPowerBoxJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return vp.unpack(enc, u.Hint, u.Total, u.VotingPowers, u.Result)
}
