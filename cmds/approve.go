package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type ApproveCommand struct {
	baseCommand
	OperationFlags
	Sender    AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract  AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	DAO       ContractIDFlag `arg:"" name:"dao-id" help:"dao id" required:"true"`
	ProposeID string         `arg:"" name:"propose-id" help:"propose id" required:"true"`
	Target    AddressFlag    `arg:"" name:"target" help:"target address to approve" required:"true"`
	Currency  CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender    base.Address
	contract  base.Address
	target    base.Address
}

func NewApproveCommand() ApproveCommand {
	cmd := NewbaseCommand()
	return ApproveCommand{
		baseCommand: *cmd,
	}
}

func (cmd *ApproveCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *ApproveCommand) parseFlags() error {
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

	target, err := cmd.Target.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid target account format, %q", cmd.Target.String())
	}
	cmd.target = target

	return nil
}

func (cmd *ApproveCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringErrorFunc("failed to create approve operation")

	fact := dao.NewApproveFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.DAO.ID,
		cmd.ProposeID,
		cmd.target,
		cmd.Currency.CID,
	)

	op, err := dao.NewApprove(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
