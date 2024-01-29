package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type PreSnapFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner      base.Address             `json:"sender"`
	Contract   base.Address             `json:"contract"`
	ProposalID string                   `json:"proposal_id"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (fact PreSnapFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PreSnapFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		ProposalID:            fact.proposalID,
		Currency:              fact.currency,
	})
}

type PreSnapFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner      string `json:"sender"`
	Contract   string `json:"contract"`
	ProposalID string `json:"proposal_id"`
	Currency   string `json:"currency"`
}

func (fact *PreSnapFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of PreSnapFact")

	var uf PreSnapFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.ProposalID,
		uf.Currency,
	)
}

type PreSnapJSONMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op PreSnap) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PreSnapJSONMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *PreSnap) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of PreSnap")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
