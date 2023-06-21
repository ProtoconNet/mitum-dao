package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type ApproveFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner     base.Address             `json:"sender"`
	Contract  base.Address             `json:"contract"`
	DAOID     currencytypes.ContractID `json:"daoid"`
	ProposeID string                   `json:"proposeid"`
	Target    base.Address             `json:"target"`
	Quorum    types.PercentRatio       `json:"quorum"`

	Currency currencytypes.CurrencyID `json:"currency"`
}

func (fact ApproveFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		DAOID:                 fact.daoID,
		ProposeID:             fact.proposeID,
		Target:                fact.target,
		Currency:              fact.currency,
	})
}

type ApproveFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner     string `json:"sender"`
	Contract  string `json:"contract"`
	DAOID     string `json:"daoid"`
	ProposeID string `json:"proposeid"`
	Target    string `json:"target"`
	Currency  string `json:"currency"`
}

func (fact *ApproveFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of ApproveFact")

	var uf ApproveFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.DAOID,
		uf.ProposeID,
		uf.Target,
		uf.Currency,
	)
}

type ApproveMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Approve) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Approve) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Approve")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
