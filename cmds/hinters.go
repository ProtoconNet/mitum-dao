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
	{Hint: types.BizProposalHint, Instance: types.BizProposal{}},
	{Hint: types.CryptoProposalHint, Instance: types.CryptoProposal{}},
	{Hint: types.DelegatorInfoHint, Instance: types.DelegatorInfo{}},
	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.GovernanceCalldataHint, Instance: types.GovernanceCallData{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.TransferCalldataHint, Instance: types.TransferCallData{}},
	{Hint: types.VoterInfoHint, Instance: types.VoterInfo{}},
	{Hint: types.VotingPowerHint, Instance: types.VotingPower{}},
	{Hint: types.VotingPowerBoxHint, Instance: types.VotingPowerBox{}},
	{Hint: types.WhitelistHint, Instance: types.Whitelist{}},

	{Hint: state.DelegatorsStateValueHint, Instance: state.DelegatorsStateValue{}},
	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.ProposalStateValueHint, Instance: state.ProposalStateValue{}},
	{Hint: state.VotersStateValueHint, Instance: state.VotersStateValue{}},
	{Hint: state.VotingPowerBoxStateValueHint, Instance: state.VotingPowerBoxStateValue{}},

	{Hint: dao.CancelProposalHint, Instance: dao.CancelProposal{}},
	{Hint: dao.CreateDAOHint, Instance: dao.CreateDAO{}},
	{Hint: dao.ExecuteHint, Instance: dao.Execute{}},
	{Hint: dao.PostSnapHint, Instance: dao.PostSnap{}},
	{Hint: dao.PreSnapHint, Instance: dao.PreSnap{}},
	{Hint: dao.ProposeHint, Instance: dao.Propose{}},
	{Hint: dao.RegisterHint, Instance: dao.Register{}},
	{Hint: dao.UpdatePolicyHint, Instance: dao.UpdatePolicy{}},
	{Hint: dao.VoteHint, Instance: dao.Vote{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: dao.CancelProposalFactHint, Instance: dao.CancelProposalFact{}},
	{Hint: dao.CreateDAOFactHint, Instance: dao.CreateDAOFact{}},
	{Hint: dao.ExecuteFactHint, Instance: dao.ExecuteFact{}},
	{Hint: dao.PostSnapFactHint, Instance: dao.PostSnapFact{}},
	{Hint: dao.PreSnapFactHint, Instance: dao.PreSnapFact{}},
	{Hint: dao.ProposeFactHint, Instance: dao.ProposeFact{}},
	{Hint: dao.RegisterFactHint, Instance: dao.RegisterFact{}},
	{Hint: dao.UpdatePolicyFactHint, Instance: dao.UpdatePolicyFact{}},
	{Hint: dao.VoteFactHint, Instance: dao.VoteFact{}},
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
