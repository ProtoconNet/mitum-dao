package types

import (
	"encoding/json"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TransferCalldataJSONMarshaler struct {
	hint.BaseHinter
	Sender   base.Address         `json:"sender"`
	Receiver base.Address         `json:"receiver"`
	Amount   currencytypes.Amount `json:"amount"`
}

func (cd TransferCallData) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferCalldataJSONMarshaler{
		BaseHinter: cd.BaseHinter,
		Sender:     cd.sender,
		Receiver:   cd.receiver,
		Amount:     cd.amount,
	})
}

type TransferCalldataJSONUnmarshaler struct {
	Hint     hint.Hint       `json:"_hint"`
	Sender   string          `json:"sender"`
	Receiver string          `json:"receiver"`
	Amount   json.RawMessage `json:"amount"`
}

func (cd *TransferCallData) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of TransferCallData")

	var uc TransferCalldataJSONUnmarshaler
	if err := enc.Unmarshal(b, &uc); err != nil {
		return e.Wrap(err)
	}

	return cd.unpack(enc, uc.Hint, uc.Sender, uc.Receiver, uc.Amount)
}

type GovernanceCalldataJSONMarshaler struct {
	hint.BaseHinter
	Policy Policy `json:"policy"`
}

func (cd GovernanceCallData) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(GovernanceCalldataJSONMarshaler{
		BaseHinter: cd.BaseHinter,
		Policy:     cd.policy,
	})
}

type GovernanceCalldataJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	Policy json.RawMessage `json:"policy"`
}

func (cd *GovernanceCallData) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of GovernanceCallData")

	var uc GovernanceCalldataJSONUnmarshaler
	if err := enc.Unmarshal(b, &uc); err != nil {
		return e.Wrap(err)
	}

	return cd.unpack(enc, uc.Hint, uc.Policy)
}
