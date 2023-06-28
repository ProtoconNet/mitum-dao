package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
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

	hr, err := enc.DecodeSlice(u.Registers)
	if err != nil {
		return e(err, "")
	}

	infos := make([]RegisterInfo, len(hr))
	for i, hinter := range hr {
		rg, ok := hinter.(RegisterInfo)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected RegisterInfo, not %T", hinter), "")
		}

		infos[i] = rg
	}
	r.Registers = infos

	return nil
}

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

type SnapHistoriesStateValueJSONMarshaler struct {
	hint.BaseHinter
	Histories []SnapHistory `json:"histories"`
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

	histories := make([]SnapHistory, len(hs))
	for i, hinter := range hs {
		h, ok := hinter.(SnapHistory)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected SnapHistory, not %T", hinter), "")
		}

		histories[i] = h
	}
	sh.Histories = histories

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

type VotesStateValueJSONMarshaler struct {
	hint.BaseHinter
	Active bool           `json:"active"`
	Result uint8          `json:"result"`
	Votes  []VotingPowers `json:"votes"`
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

	votes := make([]VotingPowers, len(hvs))
	for i, hinter := range hvs {
		c, ok := hinter.(VotingPowers)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected VotingPowers, not %T", hinter), "")
		}

		votes[i] = c
	}
	v.Votes = votes

	return nil
}
