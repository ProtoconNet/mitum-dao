package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type RegisterCommand struct {
	baseCommand
	OperationFlags
	Sender    AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract  AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	DAO       ContractIDFlag `arg:"" name:"dao-id" help:"dao id" required:"true"`
	ProposeID string         `arg:"" name:"propose-id" help:"propose id" required:"true"`
	Currency  CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	Approved  AddressFlag    `name:"approved" help:"target address to approve"`
	sender    base.Address
	contract  base.Address
	approved  base.Address
}

func NewRegisterCommand() RegisterCommand {
	cmd := NewbaseCommand()
	return RegisterCommand{
		baseCommand: *cmd,
	}
}

func (cmd *RegisterCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *RegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	if cmd.Approved.s != "" {
		ap, err := cmd.Approved.Encode(enc)
		if err != nil {
			return errors.Wrapf(err, "invalid approved account format, %q", cmd.Approved.String())
		}
		cmd.approved = ap
	} else {
		cmd.approved = nil
	}

	return nil
}

func (cmd *RegisterCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringErrorFunc("failed to create register operation")

	fact := dao.NewRegisterFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.DAO.ID,
		cmd.ProposeID,
		cmd.approved,
		cmd.Currency.CID,
	)

	op, err := dao.NewRegister(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
