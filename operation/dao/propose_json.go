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

type ProposeFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner     base.Address             `json:"sender"`
	Contract  base.Address             `json:"contract"`
	DAOID     currencytypes.ContractID `json:"daoid"`
	ProposeID string                   `json:"proposeid"`
	StartTime uint64                   `json:"starttime"`
	Proposal  types.Proposal           `json:"proposal"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (fact ProposeFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposeFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		DAOID:                 fact.daoID,
		ProposeID:             fact.proposeID,
		StartTime:             fact.starttime,
		Proposal:              fact.proposal,
		Currency:              fact.currency,
	})
}

type ProposeFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner     string          `json:"sender"`
	Contract  string          `json:"contract"`
	DAOID     string          `json:"daoid"`
	ProposeID string          `json:"proposeid"`
	StartTime uint64          `json:"starttime"`
	Proposal  json.RawMessage `json:"proposal"`
	Currency  string          `json:"currency"`
}

func (fact *ProposeFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ProposeFact")

	var uf ProposeFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.DAOID,
		uf.ProposeID,
		uf.StartTime,
		uf.Proposal,
		uf.Currency,
	)
}

type ProposeMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Propose) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ProposeMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Propose) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Propose")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
