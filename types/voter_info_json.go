package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type VoterInfoJSONMarshaler struct {
	hint.BaseHinter
	Account    base.Address   `json:"account"`
	Delegators []base.Address `json:"delegators"`
}

func (r VoterInfo) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VoterInfoJSONMarshaler{
		BaseHinter: r.BaseHinter,
		Account:    r.account,
		Delegators: r.delegators,
	})
}

type VoterInfoJSONUnmarshaler struct {
	Account    string   `json:"account"`
	Delegators []string `json:"delegators"`
}

func (r *VoterInfo) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of VoterInfo")

	var u VoterInfoJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		r.account = a
	}

	acc := make([]base.Address, len(u.Delegators))
	for i, ba := range u.Delegators {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e.Wrap(err)
		}
		acc[i] = ac

	}
	r.delegators = acc

	return nil
}
