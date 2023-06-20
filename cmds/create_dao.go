package cmds

import (
	"context"

	"github.com/pkg/errors"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type CreateDAOCommand struct {
	baseCommand
	OperationFlags
	Sender           AddressFlag        `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract         AddressFlag        `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	DAO              ContractIDFlag     `arg:"" name:"dao-id" help:"credential id" required:"true"`
	Option           string             `arg:"" name:"dao-option" help:"dao option" required:"true"`
	VotingPowerToken CurrencyIDFlag     `arg:"" name:"voting-power-token" help:"voting power token" required:"true"`
	Threshold        CurrencyAmountFlag `arg:"" name:"threshold" help:"threshold to propose" required:"true"`
	Fee              CurrencyAmountFlag `arg:"" name:"fee" help:"fee to propose" required:"true"`
	Delaytime        uint64             `arg:"" name:"delaytime" help:"delaytime" required:"true"`
	Snaptime         uint64             `arg:"" name:"snaptime" help:"snaptime" required:"true"`
	Timelock         uint64             `arg:"" name:"timelock" help:"timelock" required:"true"`
	Turnout          uint               `arg:"" name:"turnout" help:"turnout" required:"true"`
	Quorum           uint               `arg:"" name:"quorum" help:"quorum" required:"true"`
	Whitelist        AddressFlag        `name:"whitelist" help:"whitelist account"`
	Currency         CurrencyIDFlag     `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender           base.Address
	contract         base.Address
	whitelist        types.Whitelist
	threshold        currencytypes.Amount
	fee              currencytypes.Amount
}

func NewCreateDAOCommand() CreateDAOCommand {
	cmd := NewbaseCommand()
	return CreateDAOCommand{
		baseCommand: *cmd,
	}
}

func (cmd *CreateDAOCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *CreateDAOCommand) parseFlags() error {
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

	if 0 < len(cmd.Whitelist.s) {
		whitelist, err := cmd.Whitelist.Encode(enc)
		if err != nil {
			return errors.Wrapf(err, "invalid whitelist account format, %q", cmd.Whitelist.String())
		}
		cmd.whitelist = types.NewWhitelist(true, []base.Address{whitelist})
	} else {
		cmd.whitelist = types.NewWhitelist(false, []base.Address{})
	}

	cmd.threshold = currencytypes.NewAmount(cmd.Threshold.Big, cmd.Threshold.CID)
	cmd.fee = currencytypes.NewAmount(cmd.Fee.Big, cmd.Fee.CID)

	return nil
}

func (cmd *CreateDAOCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringErrorFunc("failed to create create-credential-service operation")

	fact := dao.NewCreateDAOFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.DAO.ID,
		types.DAOOption(cmd.Option),
		cmd.VotingPowerToken.CID,
		cmd.threshold,
		cmd.fee,
		cmd.whitelist,
		cmd.Delaytime,
		cmd.Snaptime,
		cmd.Timelock,
		types.PercentRatio(cmd.Turnout),
		types.PercentRatio(cmd.Quorum),
		cmd.Currency.CID,
	)

	op, err := dao.NewCreateDAO(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
