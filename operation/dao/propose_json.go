package dao

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type ProposeFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner      base.Address             `json:"sender"`
	Contract   base.Address             `json:"contract"`
	ProposalID string                   `json:"proposal_id"`
	Proposal   types.Proposal           `json:"proposal"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (fact ProposeFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposeFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		ProposalID:            fact.proposalID,
		Proposal:              fact.proposal,
		Currency:              fact.currency,
	})
}

type ProposeFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner      string          `json:"sender"`
	Contract   string          `json:"contract"`
	ProposalID string          `json:"proposal_id"`
	Proposal   json.RawMessage `json:"proposal"`
	Currency   string          `json:"currency"`
}

func (fact *ProposeFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf ProposeFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)
	if err := fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.ProposalID,
		uf.Proposal,
		uf.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type ProposeMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Propose) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposeMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Propose) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
