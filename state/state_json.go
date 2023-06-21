package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	DAO types.Design `json:"dao"`
}

func (de DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter: de.BaseHinter,
		DAO:        de.Design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	DAO json.RawMessage `json:"dao"`
}

func (de *DesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var design types.Design

	if err := design.DecodeJSON(u.DAO, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}

type ProposalStateValueJSONMarshaler struct {
	hint.BaseHinter
	Proposal types.Proposal `json:"proposal"`
}

func (p ProposalStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposalStateValueJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Proposal:   p.Proposal,
	})
}

type ProposalStateValueJSONUnmarshaler struct {
	Proposal json.RawMessage `json:"proposal"`
}

func (p *ProposalStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ProposalStateValue")

	var u ProposalStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	if hinter, err := enc.Decode(u.Proposal); err != nil {
		return e(err, "")
	} else if pr, ok := hinter.(types.Proposal); !ok {
		return e(util.ErrWrongType.Errorf("expected Proposal, not %T", hinter), "")
	} else {
		p.Proposal = pr
	}

	return nil
}

type ApprovingListStateValueJSONMarshaler struct {
	hint.BaseHinter
	Accounts []base.Address `json:"accounts"`
}

func (ap ApprovingListStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApprovingListStateValueJSONMarshaler{
		BaseHinter: ap.BaseHinter,
		Accounts:   ap.Accounts,
	})
}

type ApprovingListStateValueJSONUnmarshaler struct {
	Accounts []string `json:"accounts"`
}

func (ap *ApprovingListStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ApprovingListStateValue")

	var u ApprovingListStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	acc := make([]base.Address, len(u.Accounts))
	for i, ba := range u.Accounts {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e(err, "")
		}
		acc[i] = ac

	}
	ap.Accounts = acc

	return nil
}

type RegisterInfoJSONMarshaler struct {
	hint.BaseHinter
	Account    base.Address   `json:"account"`
	ApprovedBy []base.Address `json:"approved_by"`
}

func (r RegisterInfo) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterInfoJSONMarshaler{
		BaseHinter: r.BaseHinter,
		Account:    r.Account,
		ApprovedBy: r.ApprovedBy,
	})
}

type RegisterInfoJSONUnmarshaler struct {
	Account    string   `json:"account"`
	ApprovedBy []string `json:"approved_by:"`
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
		r.Account = a
	}

	acc := make([]base.Address, len(u.ApprovedBy))
	for i, ba := range u.ApprovedBy {
		ac, err := base.DecodeAddress(ba, enc)
		if err != nil {
			return e(err, "")
		}
		acc[i] = ac

	}
	r.ApprovedBy = acc

	return nil
}

type RegisterListStateValueJSONMarshaler struct {
	hint.BaseHinter
	Registers []RegisterInfo `json:"registers"`
}

func (r RegisterListStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterListStateValueJSONMarshaler{
		BaseHinter: r.BaseHinter,
		Registers:  r.Registers,
	})
}

type RegisterListStateValueJSONUnmarshaler struct {
	Registers json.RawMessage `json:"account"`
}

func (r *RegisterListStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RegisterListStateValue")

	var u RegisterListStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	hit, err := enc.DecodeSlice(u.Registers)
	if err != nil {
		return e(err, "")
	}

	rs := make([]RegisterInfo, len(hit))
	for i, hinter := range hit {
		rg, ok := hinter.(RegisterInfo)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected RegisterInfo, not %T", hinter), "")
		}

		rs[i] = rg
	}
	r.Registers = rs

	return nil
}
