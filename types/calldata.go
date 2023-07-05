package types

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

const (
	CalldataTransfer   = "transfer"
	CalldataGovernance = "governance"
)

var (
	TransferCalldataHint   = hint.MustNewHint("mitum-dao-transfer-calldata-v0.0.1")
	GovernanceCalldataHint = hint.MustNewHint("mitum-dao-governance-calldata-v0.0.1")
)

type CallData interface {
	util.IsValider
	hint.Hinter
	Type() string
	Bytes() []byte
	Addresses() []base.Address
}

type TransferCallData struct {
	hint.BaseHinter
	sender   base.Address
	receiver base.Address
	amount   currencytypes.Amount
}

func NewTransferCallData(sender base.Address, receiver base.Address, amount currencytypes.Amount) TransferCallData {
	return TransferCallData{
		BaseHinter: hint.NewBaseHinter(TransferCalldataHint),
		sender:     sender,
		receiver:   receiver,
		amount:     amount,
	}
}

func (TransferCallData) Type() string {
	return CalldataTransfer
}

func (cd TransferCallData) Bytes() []byte {
	return util.ConcatBytesSlice(cd.sender.Bytes(), cd.receiver.Bytes(), cd.amount.Bytes())
}

func (cd TransferCallData) Sender() base.Address {
	return cd.sender
}

func (cd TransferCallData) Receiver() base.Address {
	return cd.receiver
}

func (cd TransferCallData) Amount() currencytypes.Amount {
	return cd.amount
}

func (cd TransferCallData) IsValid([]byte) error {
	if err := cd.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, cd.sender, cd.receiver, cd.amount); err != nil {
		return util.ErrInvalid.Errorf("invalid transfer calldata: %v", err)
	}

	if !cd.amount.Big().OverZero() {
		return util.ErrInvalid.Errorf("transfer calldata - amount under zero")
	}

	if cd.sender.Equal(cd.receiver) {
		return util.ErrInvalid.Errorf("transfer calldata - sender == receiver, %s", cd.sender)
	}

	return nil
}

func (cd TransferCallData) Addresses() []base.Address {
	as := make([]base.Address, 2)

	as[0] = cd.sender
	as[1] = cd.receiver

	return as
}

type GovernanceCallData struct {
	hint.BaseHinter
	policy Policy
}

func NewGovernanceCallData(policy Policy) GovernanceCallData {
	return GovernanceCallData{
		BaseHinter: hint.NewBaseHinter(GovernanceCalldataHint),
		policy:     policy,
	}
}

func (GovernanceCallData) Type() string {
	return CalldataGovernance
}

func (cd GovernanceCallData) Bytes() []byte {
	return cd.policy.Bytes()
}

func (cd GovernanceCallData) Policy() Policy {
	return cd.policy
}

func (cd GovernanceCallData) IsValid([]byte) error {
	if err := cd.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := cd.policy.IsValid(nil); err != nil {
		return util.ErrInvalid.Errorf("governance calldata - invalid policy: %v", err)
	}

	return nil
}

func (cd GovernanceCallData) Addresses() []base.Address {
	return cd.policy.whitelist.accounts
}
