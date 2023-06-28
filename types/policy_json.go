package types

import (
	"encoding/json"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type WhitelistJSONMarshaler struct {
	hint.BaseHinter
	Active   bool           `json:"active"`
	Accounts []base.Address `json:"accounts"`
}

func (wl Whitelist) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(WhitelistJSONMarshaler{
		BaseHinter: wl.BaseHinter,
		Active:     wl.active,
		Accounts:   wl.accounts,
	})
}

type WhitelistJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Active   bool            `json:"active"`
	Accounts json.RawMessage `json:"accounts"`
}

func (wl *Whitelist) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Whitelist")

	var uw WhitelistJSONUnmarshaler
	if err := enc.Unmarshal(b, &uw); err != nil {
		return e(err, "")
	}

	return wl.unpack(enc, uw.Hint, uw.Active, uw.Accounts)
}

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	Token       currencytypes.CurrencyID `json:"token"`
	Threshold   currencytypes.Amount     `json:"threshold"`
	Fee         currencytypes.Amount     `json:"fee"`
	Whitelist   Whitelist                `json:"whitelist"`
	Delaytime   uint64                   `json:"delaytime"`
	Snaptime    uint64                   `json:"snaptime"`
	Voteoperiod uint64                   `json:"voteperiod"`
	Timelock    uint64                   `json:"timelock"`
	Turnout     PercentRatio             `json:"turnout"`
	Quorum      PercentRatio             `json:"quorum"`
}

func (po Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter:  po.BaseHinter,
		Token:       po.token,
		Threshold:   po.threshold,
		Fee:         po.fee,
		Whitelist:   po.whitelist,
		Delaytime:   po.delaytime,
		Snaptime:    po.snaptime,
		Voteoperiod: po.voteperiod,
		Timelock:    po.timelock,
		Turnout:     po.turnout,
		Quorum:      po.quorum,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint       hint.Hint       `json:"_hint"`
	Token      string          `json:"token"`
	Threshold  json.RawMessage `json:"threshold"`
	Fee        json.RawMessage `json:"fee"`
	Whitelist  json.RawMessage `json:"whitelist"`
	Delaytime  uint64          `json:"delaytime"`
	Snaptime   uint64          `json:"snaptime"`
	Voteperiod uint64          `json:"voteperiod"`
	Timelock   uint64          `json:"timelock"`
	Turnout    uint            `json:"turnout"`
	Quorum     uint            `json:"quorum"`
}

func (po *Policy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Policy")

	var upo PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e(err, "")
	}

	return po.unpack(enc, upo.Hint,
		upo.Token,
		upo.Threshold,
		upo.Fee,
		upo.Whitelist,
		upo.Delaytime,
		upo.Snaptime,
		upo.Voteperiod,
		upo.Timelock,
		upo.Turnout,
		upo.Quorum,
	)
}
