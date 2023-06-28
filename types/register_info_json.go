package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type RegisterInfoJSONMarshaler struct {
	hint.BaseHinter
	Account    base.Address   `json:"account"`
	ApprovedBy []base.Address `json:"approved_by"`
}

func (r RegisterInfo) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterInfoJSONMarshaler{
		BaseHinter: r.BaseHinter,
		Account:    r.account,
		ApprovedBy: r.approvedBy,
	})
}

type RegisterInfoJSONUnmarshaler struct {
	Account    string   `json:"account"`
	ApprovedBy []string `json:"approved_by"`
}

func (r *RegisterInfo) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RegisterInfo")

	var u RegisterInfoJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e(err, "")
	default:
		r.account = a
	}

	acc := make([]base.Address, len(u.ApprovedBy))
	for i, ba := range u.ApprovedBy {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e(err, "")
		}
		acc[i] = ac

	}
	r.approvedBy = acc

	return nil
}
