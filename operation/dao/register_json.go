package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type RegisterFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner     base.Address             `json:"sender"`
	Contract  base.Address             `json:"contract"`
	DAOID     currencytypes.ContractID `json:"daoid"`
	ProposeID string                   `json:"proposeid"`
	Approved  base.Address             `json:"approved"`
	Quorum    types.PercentRatio       `json:"quorum"`

	Currency currencytypes.CurrencyID `json:"currency"`
}

func (fact RegisterFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		DAOID:                 fact.daoID,
		ProposeID:             fact.proposeID,
		Approved:              fact.approved,
		Currency:              fact.currency,
	})
}

type RegisterFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner     string `json:"sender"`
	Contract  string `json:"contract"`
	DAOID     string `json:"daoid"`
	ProposeID string `json:"proposeid"`
	Approved  string `json:"approved"`
	Currency  string `json:"currency"`
}

func (fact *RegisterFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of RegisterFact")

	var uf RegisterFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.DAOID,
		uf.ProposeID,
		uf.Approved,
		uf.Currency,
	)
}

type RegisterMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Register) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(RegisterMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Register) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of Register")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}