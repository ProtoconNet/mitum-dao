package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CancelProposalFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner      base.Address             `json:"sender"`
	Contract   base.Address             `json:"contract"`
	DAOID      currencytypes.ContractID `json:"dao_id"`
	ProposalID string                   `json:"proposal_id"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (fact CancelProposalFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CancelProposalFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		DAOID:                 fact.daoID,
		ProposalID:            fact.proposalID,
		Currency:              fact.currency,
	})
}

type CancelProposalFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner      string `json:"sender"`
	Contract   string `json:"contract"`
	DAOID      string `json:"dao_id"`
	ProposalID string `json:"proposal_id"`
	Currency   string `json:"currency"`
}

func (fact *CancelProposalFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CancelProposalFact")

	var uf CancelProposalFactJSONUnMarshaler
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

type CancelProposalJSONMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CancelProposal) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CancelProposalJSONMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CancelProposal) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CancelProposal")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
