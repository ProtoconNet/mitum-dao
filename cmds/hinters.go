package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: types.WhitelistHint, Instance: types.Whitelist{}},
	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.TransferCalldataHint, Instance: types.TransferCallData{}},
	{Hint: types.GovernanceCalldataHint, Instance: types.GovernanceCallData{}},
	{Hint: types.CryptoProposalHint, Instance: types.CryptoProposal{}},
	{Hint: types.BizProposalHint, Instance: types.BizProposal{}},
	{Hint: types.VoterInfoHint, Instance: types.VoterInfo{}},
	{Hint: types.VotingPowerHint, Instance: types.VotingPower{}},
	{Hint: types.VotingPowerBoxHint, Instance: types.VotingPowerBox{}},
	{Hint: types.DelegatorInfoHint, Instance: types.DelegatorInfo{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.ProposalStateValueHint, Instance: state.ProposalStateValue{}},
	{Hint: state.VotersStateValueHint, Instance: state.VotersStateValue{}},
	{Hint: state.DelegatorsStateValueHint, Instance: state.DelegatorsStateValue{}},
	{Hint: state.VotingPowerBoxStateValueHint, Instance: state.VotingPowerBoxStateValue{}},

	{Hint: dao.CreateDAOHint, Instance: dao.CreateDAO{}},
	{Hint: dao.ProposeHint, Instance: dao.Propose{}},
	{Hint: dao.CancelProposalHint, Instance: dao.CancelProposal{}},
	{Hint: dao.RegisterHint, Instance: dao.Register{}},
	{Hint: dao.PreSnapHint, Instance: dao.PreSnap{}},
	{Hint: dao.VoteHint, Instance: dao.Vote{}},
	{Hint: dao.PostSnapHint, Instance: dao.PostSnap{}},
	{Hint: dao.ExecuteHint, Instance: dao.Execute{}}}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: dao.CreateDAOFactHint, Instance: dao.CreateDAOFact{}},
	{Hint: dao.ProposeFactHint, Instance: dao.ProposeFact{}},
	{Hint: dao.CancelProposalFactHint, Instance: dao.CancelProposalFact{}},
	{Hint: dao.RegisterFactHint, Instance: dao.RegisterFact{}},
	{Hint: dao.PreSnapFactHint, Instance: dao.PreSnapFact{}},
	{Hint: dao.VoteFactHint, Instance: dao.VoteFact{}},
	{Hint: dao.PostSnapFactHint, Instance: dao.PostSnapFact{}},
	{Hint: dao.ExecuteFactHint, Instance: dao.ExecuteFact{}},
}

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	allExtendedLen := currencyExtendedLen + len(AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:], AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	allSupportedExtendedLen := currencySupportedExtendedLen + len(AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:], AddedSupportedHinters)

}

func LoadHinters(enc encoder.Encoder) error {
	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
