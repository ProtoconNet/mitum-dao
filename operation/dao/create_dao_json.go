package dao

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CreateDAOFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner                base.Address             `json:"sender"`
	Contract             base.Address             `json:"contract"`
	Option               types.DAOOption          `json:"option"`
	VotingPowerToken     currencytypes.CurrencyID `json:"voting_power_token"`
	Threshold            common.Big               `json:"threshold"`
	Fee                  currencytypes.Amount     `json:"fee"`
	Whitelist            types.Whitelist          `json:"whitelist"`
	ProposalReviewPeriod uint64                   `json:"proposal_review_period"`
	RegistrationPeriod   uint64                   `json:"registration_period"`
	PreSnapshotPeriod    uint64                   `json:"pre_snapshot_period"`
	VotingPeriod         uint64                   `json:"voting_period"`
	PostSnapshotPeriod   uint64                   `json:"post_snapshot_period"`
	ExecutionDelayPeriod uint64                   `json:"execution_delay_period"`
	Turnout              types.PercentRatio       `json:"turnout"`
	Quorum               types.PercentRatio       `json:"quorum"`
	Currency             currencytypes.CurrencyID `json:"currency"`
}

func (fact CreateDAOFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateDAOFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		Option:                fact.option,
		VotingPowerToken:      fact.votingPowerToken,
		Threshold:             fact.threshold,
		Fee:                   fact.fee,
		Whitelist:             fact.whitelist,
		ProposalReviewPeriod:  fact.proposalReviewPeriod,
		RegistrationPeriod:    fact.registrationPeriod,
		PreSnapshotPeriod:     fact.preSnapshotPeriod,
		VotingPeriod:          fact.votingPeriod,
		PostSnapshotPeriod:    fact.postSnapshotPeriod,
		ExecutionDelayPeriod:  fact.executionDelayPeriod,
		Turnout:               fact.turnout,
		Quorum:                fact.quorum,
		Currency:              fact.currency,
	})
}

type CreateDAOFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner                string          `json:"sender"`
	Contract             string          `json:"contract"`
	Option               string          `json:"option"`
	VotingPowerToken     string          `json:"voting_power_token"`
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
	Currency             string          `json:"currency"`
}

func (fact *CreateDAOFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateDAOFact")

	var uf CreateDAOFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.Option,
		uf.VotingPowerToken,
		uf.Threshold,
		uf.Fee,
		uf.Whitelist,
		uf.ProposalReviewPeriod,
		uf.RegistrationPeriod,
		uf.PreSnapshotPeriod,
		uf.VotingPeriod,
		uf.PostSnapshotPeriod,
		uf.ExecutionDelayPeriod,
		uf.Turnout,
		uf.Quorum,
		uf.Currency,
	)
}

type CreateDAOMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateDAO) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateDAOMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateDAO) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateDAO")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
