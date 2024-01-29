package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type VoteFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner      base.Address             `json:"sender"`
	Contract   base.Address             `json:"contract"`
	ProposalID string                   `json:"proposal_id"`
	Vote       uint8                    `json:"vote"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (fact VoteFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VoteFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		ProposalID:            fact.proposalID,
		Vote:                  fact.vote,
		Currency:              fact.currency,
	})
}

type VoteFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner      string `json:"sender"`
	Contract   string `json:"contract"`
	ProposalID string `json:"proposal_id"`
	Vote       uint8  `json:"vote"`
	Currency   string `json:"currency"`
}

func (fact *VoteFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of VoteFact")

	var uf VoteFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.ProposalID,
		uf.Vote,
		uf.Currency,
	)
}

type VoteMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Vote) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(VoteMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Vote) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Vote")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
