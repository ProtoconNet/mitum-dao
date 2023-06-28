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
	Active   bool           `json:"active"`
	Proposal types.Proposal `json:"proposal"`
}

func (p ProposalStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposalStateValueJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Active:     p.Active,
		Proposal:   p.Proposal,
	})
}

type ProposalStateValueJSONUnmarshaler struct {
	Active   bool            `json:"active"`
	Proposal json.RawMessage `json:"proposal"`
}

func (p *ProposalStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ProposalStateValue")

	var u ProposalStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	p.Active = u.Active

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

type RegisterListStateValueJSONMarshaler struct {
	hint.BaseHinter
	Registers []types.RegisterInfo `json:"registers"`
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

	hr, err := enc.DecodeSlice(u.Registers)
	if err != nil {
		return e(err, "")
	}

	infos := make([]types.RegisterInfo, len(hr))
	for i, hinter := range hr {
		rg, ok := hinter.(types.RegisterInfo)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected types.RegisterInfo, not %T", hinter), "")
		}

		infos[i] = rg
	}
	r.Registers = infos

	return nil
}

type SnapHistoriesStateValueJSONMarshaler struct {
	hint.BaseHinter
	Histories []types.SnapHistory `json:"histories"`
}

func (sh SnapHistoriesStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SnapHistoriesStateValueJSONMarshaler{
		BaseHinter: sh.BaseHinter,
		Histories:  sh.Histories,
	})
}

type SnapHistoriesStateValueJSONUnmarshaler struct {
	Histories json.RawMessage `json:"histories"`
}

func (sh *SnapHistoriesStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of SnapHistoriesStateValue")

	var u SnapHistoriesStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	hs, err := enc.DecodeSlice(u.Histories)
	if err != nil {
		return e(err, "")
	}

	histories := make([]types.SnapHistory, len(hs))
	for i, hinter := range hs {
		h, ok := hinter.(types.SnapHistory)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected types.SnapHistory, not %T", hinter), "")
		}

		histories[i] = h
	}
	sh.Histories = histories

	return nil
}

type VotesStateValueJSONMarshaler struct {
	hint.BaseHinter
	Active bool                 `json:"active"`
	Result uint8                `json:"result"`
	Votes  []types.VotingPowers `json:"votes"`
}

func (v VotesStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotesStateValueJSONMarshaler{
		BaseHinter: v.BaseHinter,
		Active:     v.Active,
		Result:     v.Result,
		Votes:      v.Votes,
	})
}

type VotesStateValueJSONUnmarshaler struct {
	Active bool            `json:"active"`
	Result uint8           `json:"result"`
	Votes  json.RawMessage `json:"votes"`
}

func (v *VotesStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of VotesStateValue")

	var u VotesStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	v.Active = u.Active
	v.Result = u.Result

	hvs, err := enc.DecodeSlice(u.Votes)
	if err != nil {
		return e(err, "")
	}

	votes := make([]types.VotingPowers, len(hvs))
	for i, hinter := range hvs {
		c, ok := hinter.(types.VotingPowers)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected types.VotingPowers, not %T", hinter), "")
		}

		votes[i] = c
	}
	v.Votes = votes

	return nil
}
