package state

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	Design types.Design `json:"design"`
}

func (de DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter: de.BaseHinter,
		Design:     de.design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	Design json.RawMessage `json:"design"`
}

func (de *DesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	var design types.Design

	if err := design.DecodeJSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}

	de.design = design

	return nil
}

type ProposalStateValueJSONMarshaler struct {
	hint.BaseHinter
	Status   types.ProposalStatus `json:"status"`
	Proposal types.Proposal       `json:"proposal"`
}

func (p ProposalStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposalStateValueJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Status:     p.Status(),
		Proposal:   p.proposal,
	})
}

type ProposalStateValueJSONUnmarshaler struct {
	Status   uint8           `json:"status"`
	Proposal json.RawMessage `json:"proposal"`
}

func (p *ProposalStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of ProposalStateValue")

	var u ProposalStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	p.status = types.ProposalStatus(u.Status)

	if hinter, err := enc.Decode(u.Proposal); err != nil {
		return e.Wrap(err)
	} else if pr, ok := hinter.(types.Proposal); !ok {
		return e.Wrap(errors.Errorf("expected Proposal, not %T", hinter))
	} else {
		p.proposal = pr
	}

	return nil
}

type DelegatorsStateValueJSONMarshaler struct {
	hint.BaseHinter
	Delegators []types.DelegatorInfo `json:"delegators"`
}

func (dg DelegatorsStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DelegatorsStateValueJSONMarshaler{
		BaseHinter: dg.BaseHinter,
		Delegators: dg.delegators,
	})
}

type DelegatorsStateValueJSONUnmarshaler struct {
	Delegators bson.Raw `json:"delegators"`
}

func (dg *DelegatorsStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of DelegatorsStateValue")

	var u DelegatorsStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	hr, err := enc.DecodeSlice(u.Delegators)
	if err != nil {
		return err
	}

	dgs := make([]types.DelegatorInfo, len(u.Delegators))
	for i, hinter := range hr {
		if v, ok := hinter.(types.DelegatorInfo); !ok {
			return e.Wrap(errors.Errorf("expected types.DelegatorInfo, not %T", hinter))
		} else {
			dgs[i] = v
		}
	}
	dg.delegators = dgs

	return nil
}

type VotersStateValueJSONMarshaler struct {
	hint.BaseHinter
	Voters []types.VoterInfo `json:"voters"`
}

func (vt VotersStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotersStateValueJSONMarshaler{
		BaseHinter: vt.BaseHinter,
		Voters:     vt.voters,
	})
}

type VotersStateValueJSONUnmarshaler struct {
	Voters json.RawMessage `json:"voters"`
}

func (vt *VotersStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of VotersStateValue")

	var u VotersStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	hr, err := enc.DecodeSlice(u.Voters)
	if err != nil {
		return e.Wrap(err)
	}

	infos := make([]types.VoterInfo, len(hr))
	for i, hinter := range hr {
		rg, ok := hinter.(types.VoterInfo)
		if !ok {
			return e.Wrap(errors.Errorf("expected types.VoterInfo, not %T", hinter))
		}

		infos[i] = rg
	}
	vt.voters = infos

	return nil
}

//
//type SnapHistoriesStateValueJSONMarshaler struct {
//	hint.BaseHinter
//	Histories []types.SnapHistory `json:"histories"`
//}
//
//func (sh SnapHistoriesStateValue) MarshalJSON() ([]byte, error) {
//	return util.MarshalJSON(SnapHistoriesStateValueJSONMarshaler{
//		BaseHinter: sh.BaseHinter,
//		Histories:  sh.Histories,
//	})
//}
//
//type SnapHistoriesStateValueJSONUnmarshaler struct {
//	Histories json.RawMessage `json:"histories"`
//}
//
//func (sh *SnapHistoriesStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
//	e := util.StringError("failed to decode json of SnapHistoriesStateValue")
//
//	var u SnapHistoriesStateValueJSONUnmarshaler
//	if err := enc.Unmarshal(b, &u); err != nil {
//		return e.Wrap(err)
//	}
//
//	hs, err := enc.DecodeSlice(u.Histories)
//	if err != nil {
//		return e.Wrap(err)
//	}
//
//	histories := make([]types.SnapHistory, len(hs))
//	for i, hinter := range hs {
//		h, ok := hinter.(types.SnapHistory)
//		if !ok {
//			return e.Wrap(errors.Errorf("expected types.SnapHistory, not %T", hinter))
//		}
//
//		histories[i] = h
//	}
//	sh.Histories = histories
//
//	return nil
//}

type VotesStateValueJSONMarshaler struct {
	hint.BaseHinter
	VotingPowerBox types.VotingPowerBox `json:"voting_power_box"`
}

func (vb VotingPowerBoxStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotesStateValueJSONMarshaler{
		BaseHinter:     vb.BaseHinter,
		VotingPowerBox: vb.votingPowerBox,
	})
}

type VotesStateValueJSONUnmarshaler struct {
	VotingPowrBox json.RawMessage `json:"voting_power_box"`
}

func (vb *VotingPowerBoxStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of VotingPowerBoxStateValue")

	var u VotesStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	if hinter, err := enc.Decode(u.VotingPowrBox); err != nil {
		return e.Wrap(err)
	} else if v, ok := hinter.(types.VotingPowerBox); !ok {
		return e.Wrap(errors.Errorf("expected VotingPowerBox, not %T", hinter))
	} else {
		vb.votingPowerBox = v
	}

	return nil
}
