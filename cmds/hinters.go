package cmds

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/digest"
	digestisaac "github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	mitumcurrency "github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v3/operation/isaac"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	extensionstate "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum-dao/state"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: common.BaseStateHint, Instance: common.BaseState{}},
	{Hint: common.NodeHint, Instance: common.BaseNode{}},
	{Hint: currencytypes.AccountHint, Instance: currencytypes.Account{}},
	{Hint: currencytypes.AddressHint, Instance: currencytypes.Address{}},
	{Hint: currencytypes.AmountHint, Instance: currencytypes.Amount{}},
	{Hint: currencytypes.AccountKeysHint, Instance: currencytypes.BaseAccountKeys{}},
	{Hint: currencytypes.AccountKeyHint, Instance: currencytypes.BaseAccountKey{}},
	{Hint: mitumcurrency.CreateAccountsItemMultiAmountsHint, Instance: mitumcurrency.CreateAccountsItemMultiAmounts{}},
	{Hint: mitumcurrency.CreateAccountsItemSingleAmountHint, Instance: mitumcurrency.CreateAccountsItemSingleAmount{}},
	{Hint: mitumcurrency.CreateAccountsHint, Instance: mitumcurrency.CreateAccounts{}},
	{Hint: mitumcurrency.KeyUpdaterHint, Instance: mitumcurrency.KeyUpdater{}},
	{Hint: mitumcurrency.TransfersItemMultiAmountsHint, Instance: mitumcurrency.TransfersItemMultiAmounts{}},
	{Hint: mitumcurrency.TransfersItemSingleAmountHint, Instance: mitumcurrency.TransfersItemSingleAmount{}},
	{Hint: mitumcurrency.TransfersHint, Instance: mitumcurrency.Transfers{}},
	{Hint: mitumcurrency.SuffrageInflationHint, Instance: mitumcurrency.SuffrageInflation{}},
	{Hint: currencystate.AccountStateValueHint, Instance: currencystate.AccountStateValue{}},
	{Hint: currencystate.BalanceStateValueHint, Instance: currencystate.BalanceStateValue{}},

	{Hint: currencytypes.CurrencyDesignHint, Instance: currencytypes.CurrencyDesign{}},
	{Hint: currencytypes.CurrencyPolicyHint, Instance: currencytypes.CurrencyPolicy{}},
	{Hint: mitumcurrency.CurrencyRegisterHint, Instance: mitumcurrency.CurrencyRegister{}},
	{Hint: mitumcurrency.CurrencyPolicyUpdaterHint, Instance: mitumcurrency.CurrencyPolicyUpdater{}},
	{Hint: currencytypes.ContractAccountKeysHint, Instance: currencytypes.ContractAccountKeys{}},
	{Hint: extensioncurrency.CreateContractAccountsItemMultiAmountsHint, Instance: extensioncurrency.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extensioncurrency.CreateContractAccountsItemSingleAmountHint, Instance: extensioncurrency.CreateContractAccountsItemSingleAmount{}},
	{Hint: extensioncurrency.CreateContractAccountsHint, Instance: extensioncurrency.CreateContractAccounts{}},
	{Hint: extensioncurrency.WithdrawsItemMultiAmountsHint, Instance: extensioncurrency.WithdrawsItemMultiAmounts{}},
	{Hint: extensioncurrency.WithdrawsItemSingleAmountHint, Instance: extensioncurrency.WithdrawsItemSingleAmount{}},
	{Hint: extensioncurrency.WithdrawsHint, Instance: extensioncurrency.Withdraws{}},
	{Hint: mitumcurrency.GenesisCurrenciesFactHint, Instance: mitumcurrency.GenesisCurrenciesFact{}},
	{Hint: mitumcurrency.GenesisCurrenciesHint, Instance: mitumcurrency.GenesisCurrencies{}},
	{Hint: currencytypes.NilFeeerHint, Instance: currencytypes.NilFeeer{}},
	{Hint: currencytypes.FixedFeeerHint, Instance: currencytypes.FixedFeeer{}},
	{Hint: currencytypes.RatioFeeerHint, Instance: currencytypes.RatioFeeer{}},
	{Hint: extensionstate.ContractAccountStateValueHint, Instance: extensionstate.ContractAccountStateValue{}},
	{Hint: currencystate.CurrencyDesignStateValueHint, Instance: currencystate.CurrencyDesignStateValue{}},

	{Hint: types.WhitelistHint, Instance: types.Whitelist{}},
	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.TransferCalldataHint, Instance: types.TransferCalldata{}},
	{Hint: types.GovernanceCalldataHint, Instance: types.GovernanceCalldata{}},
	{Hint: types.CryptoProposalHint, Instance: types.CryptoProposal{}},
	{Hint: types.BizProposalHint, Instance: types.BizProposal{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.ProposalStateValueHint, Instance: state.ProposalStateValue{}},

	{Hint: dao.CreateDAOHint, Instance: dao.CreateDAO{}},
	{Hint: dao.ProposeHint, Instance: dao.Propose{}},

	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},

	{Hint: isaacoperation.GenesisNetworkPolicyHint, Instance: isaacoperation.GenesisNetworkPolicy{}},
	{Hint: isaacoperation.SuffrageCandidateHint, Instance: isaacoperation.SuffrageCandidate{}},
	{Hint: isaacoperation.SuffrageGenesisJoinHint, Instance: isaacoperation.SuffrageGenesisJoin{}},
	{Hint: isaacoperation.SuffrageDisjoinHint, Instance: isaacoperation.SuffrageDisjoin{}},
	{Hint: isaacoperation.SuffrageJoinHint, Instance: isaacoperation.SuffrageJoin{}},
	{Hint: isaacoperation.NetworkPolicyHint, Instance: isaacoperation.NetworkPolicy{}},
	{Hint: isaacoperation.NetworkPolicyStateValueHint, Instance: isaacoperation.NetworkPolicyStateValue{}},
	{Hint: isaacoperation.FixedSuffrageCandidateLimiterRuleHint, Instance: isaacoperation.FixedSuffrageCandidateLimiterRule{}},
	{Hint: isaacoperation.MajoritySuffrageCandidateLimiterRuleHint, Instance: isaacoperation.MajoritySuffrageCandidateLimiterRule{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: mitumcurrency.CreateAccountsFactHint, Instance: mitumcurrency.CreateAccountsFact{}},
	{Hint: mitumcurrency.KeyUpdaterFactHint, Instance: mitumcurrency.KeyUpdaterFact{}},
	{Hint: mitumcurrency.TransfersFactHint, Instance: mitumcurrency.TransfersFact{}},
	{Hint: mitumcurrency.SuffrageInflationFactHint, Instance: mitumcurrency.SuffrageInflationFact{}},

	{Hint: mitumcurrency.CurrencyRegisterFactHint, Instance: mitumcurrency.CurrencyRegisterFact{}},
	{Hint: mitumcurrency.CurrencyPolicyUpdaterFactHint, Instance: mitumcurrency.CurrencyPolicyUpdaterFact{}},
	{Hint: extensioncurrency.CreateContractAccountsFactHint, Instance: extensioncurrency.CreateContractAccountsFact{}},
	{Hint: extensioncurrency.WithdrawsFactHint, Instance: extensioncurrency.WithdrawsFact{}},

	{Hint: dao.CreateDAOFactHint, Instance: dao.CreateDAOFact{}},
	{Hint: dao.ProposeFactHint, Instance: dao.ProposeFact{}},

	{Hint: isaacoperation.GenesisNetworkPolicyFactHint, Instance: isaacoperation.GenesisNetworkPolicyFact{}},
	{Hint: isaacoperation.SuffrageCandidateFactHint, Instance: isaacoperation.SuffrageCandidateFact{}},
	{Hint: isaacoperation.SuffrageDisjoinFactHint, Instance: isaacoperation.SuffrageDisjoinFact{}},
	{Hint: isaacoperation.SuffrageJoinFactHint, Instance: isaacoperation.SuffrageJoinFact{}},
	{Hint: isaacoperation.SuffrageGenesisJoinFactHint, Instance: isaacoperation.SuffrageGenesisJoinFact{}},
}

func init() {
	Hinters = make([]encoder.DecodeDetail, len(launch.Hinters)+len(hinters))
	copy(Hinters, launch.Hinters)
	copy(Hinters[len(launch.Hinters):], hinters)

	SupportedProposalOperationFactHinters = make([]encoder.DecodeDetail, len(launch.SupportedProposalOperationFactHinters)+len(supportedProposalOperationFactHinters))
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[len(launch.SupportedProposalOperationFactHinters):], supportedProposalOperationFactHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for i := range Hinters {
		if err := enc.Add(Hinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := enc.Add(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
