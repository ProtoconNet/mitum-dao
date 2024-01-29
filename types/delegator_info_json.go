package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

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

func (r *DelegatorInfo) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
