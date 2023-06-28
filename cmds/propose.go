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

type TransferCalldataCommand struct {
	From   AddressFlag        `name:"from" help:"calldata sender"`
	To     AddressFlag        `name:"to" help:"calldata receiver"`
	Amount CurrencyAmountFlag `name:"amount" help:"calldata amount"`
}

type GovernanceCalldataCommand struct {
	VotingPowerToken CurrencyIDFlag     `name:"voting-power-token" help:"voting power token"`
	Threshold        CurrencyAmountFlag `name:"threshold" help:"threshold to propose"`
	Fee              CurrencyAmountFlag `name:"fee" help:"fee to propose"`
	Delaytime        uint64             `name:"delaytime" help:"delaytime"`
	Snaptime         uint64             `name:"snaptime" help:"snaptime"`
	Voteperiod       uint64             `name:"voteperiod" help:"voteperiod"`
	Timelock         uint64             `name:"timelock" help:"timelock"`
	Turnout          uint               `name:"turnout" help:"turnout"`
	Quorum           uint               `name:"quorum" help:"quorum"`
	Whitelist        AddressFlag        `name:"whitelist" help:"whitelist account"`
}

type CryptoProposalCommand struct {
	CalldataOption string `name:"calldata-option" help:"calldata option; transfer | governance"`
	TransferCalldataCommand
	GovernanceCalldataCommand
}

type BizProposalCommand struct {
	URL  types.URL `name:"url" help:"proposal url"`
	Hash string    `name:"hash" help:"proposal hash"`
}

type ProposeCommand struct {
	baseCommand
	OperationFlags
	Sender    AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract  AddressFlag    `arg:"" name:"contract" help:"contract address of credential" required:"true"`
	DAO       ContractIDFlag `arg:"" name:"dao-id" help:"dao id" required:"true"`
	Option    string         `arg:"" name:"option" help:"propose option; crypto | biz" required:"true"`
	ProposeID string         `arg:"" name:"propose-id" help:"propose id" required:"true"`
	StartTime uint64         `arg:"" name:"starttime" help:"start time to register" required:"true"`
	Options   uint8          `arg:"" name:"options" help:"number of vote options" required:"true"`
	CryptoProposalCommand
	BizProposalCommand
	Currency CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
	proposal types.Proposal
}

func NewProposeCommand() ProposeCommand {
	cmd := NewbaseCommand()
	return ProposeCommand{
		baseCommand: *cmd,
	}
}

func (cmd *ProposeCommand) Run(pctx context.Context) error { // nolint:dupl
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

			calldata := types.NewTransferCalldata(from, to, amount)
			if err := calldata.IsValid(nil); err != nil {
				return err
			}

			proposal := types.NewCryptoProposal(cmd.StartTime, calldata)
			if err := proposal.IsValid(nil); err != nil {
				return err
			}
			cmd.proposal = proposal
		} else if cmd.CalldataOption == types.CalldataGovernance {
			whitelist := types.NewWhitelist(false, []base.Address{})

			if 0 < len(cmd.Whitelist.s) {
				a, err := cmd.Whitelist.Encode(enc)
				if err != nil {
					return errors.Wrapf(err, "invalid whitelist account format, %q", cmd.Whitelist.String())
				}
				whitelist = types.NewWhitelist(true, []base.Address{a})
			}

			threshold := currencytypes.NewAmount(cmd.Threshold.Big, cmd.Threshold.CID)
			fee := currencytypes.NewAmount(cmd.Fee.Big, cmd.Fee.CID)

			policy := types.NewPolicy(
				cmd.VotingPowerToken.CID,
				fee, threshold, whitelist,
				cmd.Delaytime, cmd.Snaptime, cmd.Voteperiod, cmd.Timelock,
				types.PercentRatio(cmd.Turnout), types.PercentRatio(cmd.Quorum),
			)
			if err := policy.IsValid(nil); err != nil {
				return err
			}

			calldata := types.NewGovernanceCalldata(policy)
			if err := calldata.IsValid(nil); err != nil {
				return err
			}

			proposal := types.NewCryptoProposal(cmd.StartTime, calldata)
			if err := proposal.IsValid(nil); err != nil {
				return err
			}
			cmd.proposal = proposal
		} else {
			return errors.Errorf("invalid calldata option, %s", cmd.CalldataOption)
		}
	} else if cmd.Option == types.ProposalBiz {
		proposal := types.NewBizProposal(cmd.StartTime, cmd.URL, cmd.Hash, cmd.Options)
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
	e := util.StringErrorFunc("failed to create propose operation")

	fact := dao.NewProposeFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.DAO.ID,
		cmd.ProposeID,
		cmd.StartTime,
		cmd.proposal,
		cmd.Currency.CID,
	)

	op, err := dao.NewPropose(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
