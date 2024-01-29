package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CryptoProposalJSONMarshaler struct {
	hint.BaseHinter
	Proposer  base.Address `json:"proposer"`
	StartTime uint64       `json:"start_time"`
	CallData  CallData     `json:"call_data"`
}

func (p CryptoProposal) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CryptoProposalJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Proposer:   p.proposer,
		CallData:   p.callData,
		StartTime:  p.startTime,
	})
}

type CryptoProposalJSONUnmarshaler struct {
	Hint      hint.Hint       `json:"_hint"`
	Proposer  string          `json:"proposer"`
	StartTime uint64          `json:"start_time"`
	CallData  json.RawMessage `json:"call_data"`
}

func (p *CryptoProposal) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of CryptoProposal")

	var up CryptoProposalJSONUnmarshaler
	if err := enc.Unmarshal(b, &up); err != nil {
		return e.Wrap(err)
	}

	return p.unpack(enc, up.Hint, up.Proposer, up.StartTime, up.CallData)
}

type BizProposalJSONMarshaler struct {
	hint.BaseHinter
	Proposer  base.Address `json:"proposer"`
	StartTime uint64       `json:"start_time"`
	Url       URL          `json:"url"`
	Hash      string       `json:"hash"`
	Options   uint8        `json:"options"`
}

func (p BizProposal) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(BizProposalJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Proposer:   p.proposer,
		StartTime:  p.startTime,
		Url:        p.url,
		Hash:       p.hash,
		Options:    p.options,
	})
}

type BizProposalJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Proposer  string    `json:"proposer"`
	StartTime uint64    `json:"start_time"`
	Url       string    `json:"url"`
	Hash      string    `json:"hash"`
	Options   uint8     `json:"options"`
}

func (p *BizProposal) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of BizProposal")

	var up BizProposalJSONUnmarshaler
	if err := enc.Unmarshal(b, &up); err != nil {
		return e.Wrap(err)
	}

	return p.unpack(enc, up.Hint, up.Proposer, up.StartTime, up.Url, up.Hash, up.Options)
}
