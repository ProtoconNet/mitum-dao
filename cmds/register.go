package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type RegisterCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender     currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	DAO        currencycmds.ContractIDFlag `arg:"" name:"dao-id" help:"dao id" required:"true"`
	ProposalID string                      `arg:"" name:"proposal-id" help:"proposal id" required:"true"`
	Delegated  currencycmds.AddressFlag    `arg:"" name:"delegated" help:"target address to be delegated" required:"true"`
	Currency   currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender     base.Address
	contract   base.Address
	delegated  base.Address
}

func (cmd *RegisterCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.Encoders
	enc = cmd.Encoder

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	currencycmds.PrettyPrint(cmd.Out, op)

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

	delegated, err := cmd.Delegated.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid delegated account format, %q", cmd.Delegated.String())
	}
	cmd.delegated = delegated

	return nil
}

func (cmd *RegisterCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError("failed to create register operation")

	fact := dao.NewRegisterFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.DAO.ID,
		cmd.ProposalID,
		cmd.delegated,
		cmd.Currency.CID,
	)

	op, err := dao.NewRegister(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
