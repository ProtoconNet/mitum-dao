package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-dao/types"
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
