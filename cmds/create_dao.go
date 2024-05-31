package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type CreateDAOCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender               currencycmds.AddressFlag        `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract             currencycmds.AddressFlag        `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	Option               string                          `arg:"" name:"dao-option" help:"dao option" required:"true"`
	VotingPowerToken     currencycmds.CurrencyIDFlag     `arg:"" name:"voting-power-token" help:"voting power token" required:"true"`
	Threshold            currencycmds.BigFlag            `arg:"" name:"threshold" help:"threshold to propose" required:"true"`
	Fee                  currencycmds.CurrencyAmountFlag `arg:"" name:"fee" help:"fee to propose" required:"true"`
	ProposalReviewPeriod uint64                          `arg:"" name:"proposal-review-period" help:"proposal review period" required:"true"`
	RegistrationPeriod   uint64                          `arg:"" name:"registration-period" help:"registration period" required:"true"`
	PreSnapshotPeriod    uint64                          `arg:"" name:"pre-snapshot-period" help:"pre snapshot period" required:"true"`
	VotingPeriod         uint64                          `arg:"" name:"voting-period" help:"voting period" required:"true"`
	PostSnapshotPeriod   uint64                          `arg:"" name:"post-snapshot-period" help:"post snapshot period" required:"true"`
	ExecutionDelayPeriod uint64                          `arg:"" name:"execution-delay-period" help:"execution delay period" required:"true"`
	Turnout              uint                            `arg:"" name:"turnout" help:"turnout" required:"true"`
	Quorum               uint                            `arg:"" name:"quorum" help:"quorum" required:"true"`
	Whitelist            currencycmds.AddressFlag        `name:"whitelist" help:"whitelist account"`
	Currency             currencycmds.CurrencyIDFlag     `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender               base.Address
	contract             base.Address
	whitelist            types.Whitelist
	fee                  currencytypes.Amount
}

func (cmd *CreateDAOCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

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

func (cmd *CreateDAOCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	if 0 < len(cmd.Whitelist.String()) {
		whitelist, err := cmd.Whitelist.Encode(cmd.Encoders.JSON())
		if err != nil {
			return errors.Wrapf(err, "invalid whitelist account format, %q", cmd.Whitelist.String())
		}
		cmd.whitelist = types.NewWhitelist(true, []base.Address{whitelist})
	} else {
		cmd.whitelist = types.NewWhitelist(false, []base.Address{})
	}

	cmd.fee = currencytypes.NewAmount(cmd.Fee.Big, cmd.Fee.CID)

	return nil
}

func (cmd *CreateDAOCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError("failed to create create-dao operation")

	fact := dao.NewCreateDAOFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		types.DAOOption(cmd.Option),
		cmd.VotingPowerToken.CID,
		cmd.Threshold.Big,
		cmd.fee,
		cmd.whitelist,
		cmd.ProposalReviewPeriod,
		cmd.RegistrationPeriod,
		cmd.PreSnapshotPeriod,
		cmd.VotingPeriod,
		cmd.PostSnapshotPeriod,
		cmd.ExecutionDelayPeriod,
		types.PercentRatio(cmd.Turnout),
		types.PercentRatio(cmd.Quorum),
		cmd.Currency.CID,
	)

	op := dao.NewCreateDAO(fact)

	err := op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
