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

type TransferCallDataCommand struct {
	From   currencycmds.AddressFlag        `name:"from" help:"call data sender"`
	To     currencycmds.AddressFlag        `name:"to" help:"call data receiver"`
	Amount currencycmds.CurrencyAmountFlag `name:"amount" help:"call data amount"`
}

type GovernanceCallDataCommand struct {
	VotingPowerToken     currencycmds.CurrencyIDFlag     `name:"voting-power-token" help:"voting power token"`
	Threshold            currencycmds.BigFlag            `name:"threshold" help:"threshold to propose"`
	Fee                  currencycmds.CurrencyAmountFlag `name:"fee" help:"fee to propose"`
	ProposalReviewPeriod uint64                          `name:"proposal-review-period" help:"proposal review period"`
	RegistrationPeriod   uint64                          `name:"registration-period" help:"registration period"`
	PreSnapshotPeriod    uint64                          `name:"pre-snapshot-period" help:"pre snapshot period"`
	VotingPeriod         uint64                          `name:"voting-period" help:"voting period"`
	PostSnapshotPeriod   uint64                          `name:"post-snapshot-period" help:"post snapshot period"`
	ExecutionDelayPeriod uint64                          `name:"execution-delay-period" help:"execution delay period"`
	Turnout              uint                            `name:"turnout" help:"turnout"`
	Quorum               uint                            `name:"quorum" help:"quorum"`
	Whitelist            currencycmds.AddressFlag        `name:"whitelist" help:"whitelist account"`
}

type CryptoProposalCommand struct {
	CalldataOption string `name:"calldata-option" help:"calldata option; transfer | governance"`
	TransferCallDataCommand
	GovernanceCallDataCommand
}

type BizProposalCommand struct {
	URL     types.URL `name:"url" help:"proposal url"`
	Hash    string    `name:"hash" help:"proposal hash"`
	Options uint8     `name:"options" help:"number of vote options"`
}

type ProposeCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender     currencycmds.AddressFlag `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   currencycmds.AddressFlag `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	Option     types.DAOOption          `arg:"" name:"option" help:"propose option; crypto | biz" required:"true"`
	ProposalID string                   `arg:"" name:"proposal-id" help:"proposal id" required:"true"`
	StartTime  uint64                   `arg:"" name:"start-time" help:"start time to proposal lifecycle" required:"true"`
	CryptoProposalCommand
	BizProposalCommand
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
	proposal types.Proposal
}

func (cmd *ProposeCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *ProposeCommand) parseFlags() error {
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

	if cmd.Option == types.ProposalCrypto {
		if cmd.CalldataOption == types.CalldataTransfer {
			from, err := cmd.From.Encode(enc)
			if err != nil {
				return errors.Wrapf(err, "invalid from address format, %q", cmd.From.String())
			}

			to, err := cmd.To.Encode(enc)
			if err != nil {
				return errors.Wrapf(err, "invalid to address format, %q", cmd.To.String())
			}

			amount := currencytypes.NewAmount(cmd.Amount.Big, cmd.Amount.CID)

			callData := types.NewTransferCallData(from, to, amount)
			if err := callData.IsValid(nil); err != nil {
				return err
			}

			proposal := types.NewCryptoProposal(sender, cmd.StartTime, callData)
			if err := proposal.IsValid(nil); err != nil {
				return err
			}
			cmd.proposal = proposal
		} else if cmd.CalldataOption == types.CalldataGovernance {
			whitelist := types.NewWhitelist(false, []base.Address{})

			if 0 < len(cmd.Whitelist.String()) {
				a, err := cmd.Whitelist.Encode(enc)
				if err != nil {
					return errors.Wrapf(err, "invalid whitelist account format, %q", cmd.Whitelist.String())
				}
				whitelist = types.NewWhitelist(true, []base.Address{a})
			}

			fee := currencytypes.NewAmount(cmd.Fee.Big, cmd.Fee.CID)

			policy := types.NewPolicy(
				cmd.VotingPowerToken.CID, cmd.Threshold.Big,
				fee, whitelist,
				cmd.ProposalReviewPeriod,
				cmd.RegistrationPeriod,
				cmd.PreSnapshotPeriod,
				cmd.VotingPeriod,
				cmd.PostSnapshotPeriod,
				cmd.ExecutionDelayPeriod,
				types.PercentRatio(cmd.Turnout), types.PercentRatio(cmd.Quorum),
			)
			if err := policy.IsValid(nil); err != nil {
				return err
			}

			calldata := types.NewGovernanceCallData(policy)
			if err := calldata.IsValid(nil); err != nil {
				return err
			}

			proposal := types.NewCryptoProposal(sender, cmd.StartTime, calldata)
			if err := proposal.IsValid(nil); err != nil {
				return err
			}
			cmd.proposal = proposal
		} else {
			return errors.Errorf("invalid calldata option, %s", cmd.CalldataOption)
		}
	} else if cmd.Option == types.ProposalBiz {
		proposal := types.NewBizProposal(sender, cmd.StartTime, cmd.URL, cmd.Hash, cmd.Options)
		if err := proposal.IsValid(nil); err != nil {
			return err
		}
		cmd.proposal = proposal
	} else {
		return errors.Errorf("invalid proposal option, %s", cmd.Option)
	}

	return nil
}

func (cmd *ProposeCommand) createOperation() (base.Operation, error) { // nolint:dupl}
	e := util.StringError("failed to create propose operation")

	fact := dao.NewProposeFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.ProposalID,
		cmd.proposal,
		cmd.Currency.CID,
	)

	op, err := dao.NewPropose(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
