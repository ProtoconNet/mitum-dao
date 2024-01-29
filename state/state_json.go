package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
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

func (de *DesignStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	var design types.Design

	if err := design.DecodeJSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	} else if err = design.IsValid(nil); err != nil {
		return e.Wrap(err)
	} else {
		de.design = design
	}

	return nil
}

type ProposalStateValueJSONMarshaler struct {
	hint.BaseHinter
	Status   types.ProposalStatus `json:"status"`
	Proposal types.Proposal       `json:"proposal"`
	Policy   types.Policy         `json:"policy"`
}

func (p ProposalStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposalStateValueJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Status:     p.Status(),
		Proposal:   p.proposal,
		Policy:     p.policy,
	})
}

type ProposalStateValueJSONUnmarshaler struct {
	Status   uint8           `json:"status"`
	Proposal json.RawMessage `json:"proposal"`
	Policy   json.RawMessage `json:"policy"`
}

func (p *ProposalStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
	} else if err := pr.IsValid(nil); err != nil {
		return e.Wrap(err)
	} else {
		p.proposal = pr
	}

	if hinter, err := enc.Decode(u.Policy); err != nil {
		return e.Wrap(err)
	} else if po, ok := hinter.(types.Policy); !ok {
		return e.Wrap(errors.Errorf("expected Policy, not %T", hinter))
	} else if err := po.IsValid(nil); err != nil {
		return e.Wrap(err)
	} else {
		p.policy = po
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
	Delegators json.RawMessage `json:"delegators"`
}

func (dg *DelegatorsStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DelegatorsStateValue")

	var u DelegatorsStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	hr, err := enc.DecodeSlice(u.Delegators)
	if err != nil {
		return err
	}

	dgs := make([]types.DelegatorInfo, len(hr))
	for i, hinter := range hr {
		if v, ok := hinter.(types.DelegatorInfo); !ok {
			return e.Wrap(errors.Errorf("expected types.DelegatorInfo, not %T", hinter))
		} else if err := v.IsValid(nil); err != nil {
			return e.Wrap(err)
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

func (vt *VotersStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
		if rg, ok := hinter.(types.VoterInfo); !ok {
			return e.Wrap(errors.Errorf("expected types.VoterInfo, not %T", hinter))
		} else if err := rg.IsValid(nil); err != nil {
			return e.Wrap(err)
		} else {
			infos[i] = rg
		}
	}
	vt.voters = infos

	return nil
}

type VotingPowerBoxStateValueJSONMarshaler struct {
	hint.BaseHinter
	VotingPowerBox types.VotingPowerBox `json:"voting_power_box"`
}

func (vb VotingPowerBoxStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VotingPowerBoxStateValueJSONMarshaler{
		BaseHinter:     vb.BaseHinter,
		VotingPowerBox: vb.votingPowerBox,
	})
}

type VotingPowerBoxStateValueJSONUnmarshaler struct {
	VotingPowerBox json.RawMessage `json:"voting_power_box"`
}

func (vb *VotingPowerBoxStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of VotingPowerBoxStateValue")

	var u VotingPowerBoxStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	var vpb types.VotingPowerBox
	if err := vpb.DecodeJSON(u.VotingPowerBox, enc); err != nil {
		return e.Wrap(err)
	} else if err = vpb.IsValid(nil); err != nil {
		return e.Wrap(err)
	} else {
		vb.votingPowerBox = vpb
	}

	return nil
}
