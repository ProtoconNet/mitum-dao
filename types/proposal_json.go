package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CryptoProposalJSONMarshaler struct {
	hint.BaseHinter
	StartTime uint64   `json:"start_time"`
	CallData  CallData `json:"call_data"`
}

func (p CryptoProposal) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CryptoProposalJSONMarshaler{
		BaseHinter: p.BaseHinter,
		CallData:   p.callData,
		StartTime:  p.startTime,
	})
}

type CryptoProposalJSONUnmarshaler struct {
	Hint      hint.Hint       `json:"_hint"`
	StartTime uint64          `json:"start_time"`
	CallData  json.RawMessage `json:"call_data"`
}

func (p *CryptoProposal) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CryptoProposal")

	var up CryptoProposalJSONUnmarshaler
	if err := enc.Unmarshal(b, &up); err != nil {
		return e(err, "")
	}

	return p.unpack(enc, up.Hint, up.StartTime, up.CallData)
}

type BizProposalJSONMarshaler struct {
	hint.BaseHinter
	StartTime uint64 `json:"start_time"`
	Url       URL    `json:"url"`
	Hash      string `json:"hash"`
}

func (p BizProposal) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(BizProposalJSONMarshaler{
		BaseHinter: p.BaseHinter,
		StartTime:  p.startTime,
		Url:        p.url,
		Hash:       p.hash,
	})
}

type BizProposalJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	StartTime uint64    `json:"start_time"`
	Url       string    `json:"url"`
	Hash      string    `json:"hash"`
}

func (p *BizProposal) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of BizProposal")

	var up BizProposalJSONUnmarshaler
	if err := enc.Unmarshal(b, &up); err != nil {
		return e(err, "")
	}

	return p.unpack(enc, up.Hint, up.StartTime, up.Url, up.Hash)
}
