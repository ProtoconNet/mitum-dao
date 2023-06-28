package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type VotingPowerJSONMarshaler struct {
	hint.BaseHinter
	Account     base.Address `json:"account"`
	VotingPower string       `json:"voting_power"`
}

func (vp VotingPower) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotingPowerJSONMarshaler{
		BaseHinter:  vp.BaseHinter,
		Account:     vp.account,
		VotingPower: vp.votingPower.String(),
	})
}

type VotingPowerJSONUnmarshaler struct {
	Account     string `json:"account"`
	VotingPower string `json:"voting_power"`
}

func (vp *VotingPower) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of VotingPower")

	var u VotingPowerJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e(err, "")
	default:
		vp.account = a
	}

	big, err := common.NewBigFromString(u.VotingPower)
	if err != nil {
		return e(err, "")
	}
	vp.votingPower = big

	return nil
}

type VotingPowersJSONMarshaler struct {
	hint.BaseHinter
	Total        string        `json:"total"`
	VotingPowers []VotingPower `json:"voting_powers"`
}

func (v VotingPowers) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotingPowersJSONMarshaler{
		BaseHinter:   v.BaseHinter,
		Total:        v.total.String(),
		VotingPowers: v.votingPowers,
	})
}

type VotingPowersJSONUnmarshaler struct {
	Total        string          `json:"total"`
	VotingPowers json.RawMessage `json:"voting_powers"`
}

func (v *VotingPowers) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("faileod to decde json of VotingPowers")

	var u VotingPowersJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	big, err := common.NewBigFromString(u.Total)
	if err != nil {
		return e(err, "")
	}
	v.total = big

	hv, err := enc.DecodeSlice(u.VotingPowers)
	if err != nil {
		return e(err, "")
	}

	vps := make([]VotingPower, len(hv))
	for i, hinter := range hv {
		vp, ok := hinter.(VotingPower)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected VotingPower, not %T", hinter), "")
		}

		vps[i] = vp
	}
	v.votingPowers = vps

	return nil
}
