package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type PostSnapFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner      base.Address             `json:"sender"`
	Contract   base.Address             `json:"contract"`
	ProposalID string                   `json:"proposal_id"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (fact PostSnapFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PostSnapFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		ProposalID:            fact.proposalID,
		Currency:              fact.currency,
	})
}

type PostSnapFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner      string `json:"sender"`
	Contract   string `json:"contract"`
	ProposalID string `json:"proposal_id"`
	Currency   string `json:"currency"`
}

func (fact *PostSnapFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of PostSnapFact")

	var uf PostSnapFactJSONUnMarshaler
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

type PostSnapJSONMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op PostSnap) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PostSnapJSONMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *PostSnap) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of PostSnap")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
