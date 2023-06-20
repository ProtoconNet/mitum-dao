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

type Calldata interface {
	util.IsValider
	hint.Hinter
	Type() string
	Bytes() []byte
}

type TransferCalldata struct {
	hint.BaseHinter
	sender   base.Address
	receiver base.Address
	amount   currencytypes.Amount
}

func NewTransferCalldata(sender base.Address, receiver base.Address, amount currencytypes.Amount) TransferCalldata {
	return TransferCalldata{
		BaseHinter: hint.NewBaseHinter(TransferCalldataHint),
		sender:     sender,
		receiver:   receiver,
		amount:     amount,
	}
}

func (TransferCalldata) Type() string {
	return CalldataTransfer
}

func (cd TransferCalldata) Bytes() []byte {
	return util.ConcatBytesSlice(cd.sender.Bytes(), cd.receiver.Bytes(), cd.amount.Bytes())
}

func (cd TransferCalldata) Sender() base.Address {
	return cd.sender
}

func (cd TransferCalldata) Receiver() base.Address {
	return cd.receiver
}

func (cd TransferCalldata) Amount() currencytypes.Amount {
	return cd.amount
}

func (cd TransferCalldata) IsValid([]byte) error {
	if err := cd.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false, cd.sender, cd.receiver, cd.amount); err != nil {
		return util.ErrInvalid.Errorf("invalid transfer calldata: %w", err)
	}

	if !cd.amount.Big().OverZero() {
		return util.ErrInvalid.Errorf("transfer calldata - amount under zero")
	}

	if cd.sender.Equal(cd.receiver) {
		return util.ErrInvalid.Errorf("transfer calldata - sender == receiver, %s", cd.sender)
	}

	return nil
}

type GovernanceCalldata struct {
	hint.BaseHinter
	policy Policy
}

func NewGovernanceCalldata(policy Policy) GovernanceCalldata {
	return GovernanceCalldata{
		BaseHinter: hint.NewBaseHinter(GovernanceCalldataHint),
		policy:     policy,
	}
}

func (GovernanceCalldata) Type() string {
	return CalldataGovernance
}

func (cd GovernanceCalldata) Bytes() []byte {
	return cd.policy.Bytes()
}

func (cd GovernanceCalldata) Policy() Policy {
	return cd.policy
}

func (cd GovernanceCalldata) IsValid([]byte) error {
	if err := cd.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := cd.policy.IsValid(nil); err != nil {
		return util.ErrInvalid.Errorf("governance calldata - invalid policy: %w", err)
	}

	return nil
}
