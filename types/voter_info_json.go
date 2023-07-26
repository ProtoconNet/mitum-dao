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
	Delegators []base.Address `json:"delegators"`
}

func (r VoterInfo) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterInfoJSONMarshaler{
		BaseHinter: r.BaseHinter,
		Account:    r.account,
		Delegators: r.delegators,
	})
}

type RegisterInfoJSONUnmarshaler struct {
	Account    string   `json:"account"`
	Delegators []string `json:"delegators"`
}

func (r *VoterInfo) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of VoterInfo")

	var u RegisterInfoJSONUnmarshaler
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

type DelegatorInfoJSONMarshaler struct {
	hint.BaseHinter
	Account   base.Address `json:"account"`
	Delegatee base.Address `json:"delegatee"`
}

func (r DelegatorInfo) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DelegatorInfoJSONMarshaler{
		BaseHinter: r.BaseHinter,
		Account:    r.account,
		Delegatee:  r.delegatee,
	})
}

type DelegatorInfoJSONUnmarshaler struct {
	Account   string `json:"account"`
	Delegatee string `json:"delegatee"`
}

func (r *DelegatorInfo) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of DelegatorInfo")

	var u DelegatorInfoJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	switch a, err := base.DecodeAddress(u.Account, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		r.account = a
	}

	switch a, err := base.DecodeAddress(u.Delegatee, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		r.delegatee = a
	}

	return nil
}
