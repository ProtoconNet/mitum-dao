package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type SnapFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner     base.Address             `json:"sender"`
	Contract  base.Address             `json:"contract"`
	DAOID     currencytypes.ContractID `json:"daoid"`
	ProposeID string                   `json:"proposeid"`
	Currency  currencytypes.CurrencyID `json:"currency"`
}

func (fact SnapFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SnapFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		DAOID:                 fact.daoID,
		ProposeID:             fact.proposeID,
		Currency:              fact.currency,
	})
}

type SnapFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner     string `json:"sender"`
	Contract  string `json:"contract"`
	DAOID     string `json:"daoid"`
	ProposeID string `json:"proposeid"`
	Currency  string `json:"currency"`
}

func (fact *SnapFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of SnapFact")

	var uf SnapFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.DAOID,
		uf.ProposeID,
		uf.Currency,
	)
}

type SnapMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Snap) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SnapMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Snap) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Snap")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
