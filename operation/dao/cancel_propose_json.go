package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CancelProposeFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner      base.Address             `json:"sender"`
	Contract   base.Address             `json:"contract"`
	DAOID      currencytypes.ContractID `json:"dao_id"`
	ProposalID string                   `json:"proposal_id"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (fact CancelProposeFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CancelProposeFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		DAOID:                 fact.daoID,
		ProposalID:            fact.proposalID,
		Currency:              fact.currency,
	})
}

type CancelProposeFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner      string `json:"sender"`
	Contract   string `json:"contract"`
	DAOID      string `json:"dao_id"`
	ProposalID string `json:"proposal_id"`
	Currency   string `json:"currency"`
}

func (fact *CancelProposeFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CancelProposeFact")

	var uf CancelProposeFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.DAOID,
		uf.ProposalID,
		uf.Currency,
	)
}

type CancelProposeJSONMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CancelPropose) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CancelProposeJSONMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CancelPropose) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CancelPropose")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
