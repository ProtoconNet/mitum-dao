package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type CryptoProposalJSONMarshaler struct {
	hint.BaseHinter
	StartTime uint64   `json:"starttime"`
	Calldata  Calldata `json:"calldata"`
}

func (p CryptoProposal) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CryptoProposalJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Calldata:   p.calldata,
		StartTime:  p.starttime,
	})
}

type CryptoProposalJSONUnmarshaler struct {
	Hint      hint.Hint       `json:"_hint"`
	StartTime uint64          `json:"starttime"`
	Calldata  json.RawMessage `json:"calldata"`
}

func (p *CryptoProposal) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CryptoProposal")

	var up CryptoProposalJSONUnmarshaler
	if err := enc.Unmarshal(b, &up); err != nil {
		return e(err, "")
	}

	return p.unpack(enc, up.Hint, up.StartTime, up.Calldata)
}

type BizProposalJSONMarshaler struct {
	hint.BaseHinter
	StartTime uint64 `json:"starttime"`
	Url       URL    `json:"url"`
	Hash      string `json:"hash"`
}

func (p BizProposal) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(BizProposalJSONMarshaler{
		BaseHinter: p.BaseHinter,
		StartTime:  p.starttime,
		Url:        p.url,
		Hash:       p.hash,
	})
}

type BizProposalJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	StartTime uint64    `json:"starttime"`
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
