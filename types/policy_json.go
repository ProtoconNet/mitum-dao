package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
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
	Hint     hint.Hint `json:"_hint"`
	Active   bool      `json:"active"`
	Accounts []string  `json:"accounts"`
}

func (wl *Whitelist) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Whitelist")

	var uw WhitelistJSONUnmarshaler
	if err := enc.Unmarshal(b, &uw); err != nil {
		return e.Wrap(err)
	}

	return wl.unpack(enc, uw.Hint, uw.Active, uw.Accounts)
}

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	Token                currencytypes.CurrencyID `json:"token"`
	Threshold            common.Big               `json:"threshold"`
	Fee                  currencytypes.Amount     `json:"fee"`
	Whitelist            Whitelist                `json:"whitelist"`
	ProposalReviewPeriod uint64                   `json:"proposal_review_period"`
	RegistrationPeriod   uint64                   `json:"registration_period"`
	PreSnapshotPeriod    uint64                   `json:"pre_snapshot_period"`
	VotingPeriod         uint64                   `json:"voting_period"`
	PostSnapshotPeriod   uint64                   `json:"post_snapshot_period"`
	ExecutionDelayPeriod uint64                   `json:"execution_delay_period"`
	Turnout              PercentRatio             `json:"turnout"`
	Quorum               PercentRatio             `json:"quorum"`
}

func (po Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter:           po.BaseHinter,
		Token:                po.token,
		Threshold:            po.threshold,
		Fee:                  po.fee,
		Whitelist:            po.whitelist,
		ProposalReviewPeriod: po.proposalReviewPeriod,
		RegistrationPeriod:   po.registrationPeriod,
		PreSnapshotPeriod:    po.preSnapshotPeriod,
		VotingPeriod:         po.votingPeriod,
		PostSnapshotPeriod:   po.postSnapshotPeriod,
		ExecutionDelayPeriod: po.executionDelayPeriod,
		Turnout:              po.turnout,
		Quorum:               po.quorum,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint                 hint.Hint       `json:"_hint"`
	Token                string          `json:"token"`
	Threshold            string          `json:"threshold"`
	Fee                  json.RawMessage `json:"fee"`
	Whitelist            json.RawMessage `json:"whitelist"`
	ProposalReviewPeriod uint64          `json:"proposal_review_period"`
	RegistrationPeriod   uint64          `json:"registration_period"`
	PreSnapshotPeriod    uint64          `json:"pre_snapshot_period"`
	VotingPeriod         uint64          `json:"voting_period"`
	PostSnapshotPeriod   uint64          `json:"post_snapshot_period"`
	ExecutionDelayPeriod uint64          `json:"execution_delay_period"`
	Turnout              uint            `json:"turnout"`
	Quorum               uint            `json:"quorum"`
}

func (po *Policy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Policy")

	var upo PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &upo); err != nil {
		return e.Wrap(err)
	}

	return po.unpack(enc, upo.Hint,
		upo.Token,
		upo.Threshold,
		upo.Fee,
		upo.Whitelist,
		upo.ProposalReviewPeriod,
		upo.RegistrationPeriod,
		upo.PreSnapshotPeriod,
		upo.VotingPeriod,
		upo.PostSnapshotPeriod,
		upo.ExecutionDelayPeriod,
		upo.Turnout,
		upo.Quorum,
	)
}
