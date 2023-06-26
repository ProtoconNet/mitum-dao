package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
)

type OperationCommand struct {
	CreateAccount         currencycmds.CreateAccountCommand         `cmd:"" name:"create-account" help:"create new account"`
	KeyUpdater            currencycmds.KeyUpdaterCommand            `cmd:"" name:"key-updater" help:"update account keys"`
	Transfer              currencycmds.TransferCommand              `cmd:"" name:"transfer" help:"transfer amounts to receiver"`
	CreateContractAccount currencycmds.CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw              currencycmds.WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw amounts from target contract account"`
	CreateDAO             CreateDAOCommand                          `cmd:"" name:"create-dao" help:"create dao to contract account"`
	Propose               ProposeCommand                            `cmd:"" name:"propose" help:"propose new proposal"`
	Register              RegisterCommand                           `cmd:"" name:"register" help:"register to vote"`
	CurrencyRegister      currencycmds.CurrencyRegisterCommand      `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater currencycmds.CurrencyPolicyUpdaterCommand `cmd:"" name:"currency-policy-updater" help:"update currency policy"`
	SuffrageInflation     currencycmds.SuffrageInflationCommand     `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"`
	SuffrageCandidate     currencycmds.SuffrageCandidateCommand     `cmd:"" name:"suffrage-candidate" help:"suffrage candidate operation"`
	SuffrageJoin          currencycmds.SuffrageJoinCommand          `cmd:"" name:"suffrage-join" help:"suffrage join operation"`
	SuffrageDisjoin       currencycmds.SuffrageDisjoinCommand       `cmd:"" name:"suffrage-disjoin" help:"suffrage disjoin operation"` // revive:disable-line:line-length-limit
}

func NewOperationCommand() OperationCommand {
	return OperationCommand{
		CreateAccount:         currencycmds.NewCreateAccountCommand(),
		KeyUpdater:            currencycmds.NewKeyUpdaterCommand(),
		Transfer:              currencycmds.NewTransferCommand(),
		CreateContractAccount: currencycmds.NewCreateContractAccountCommand(),
		Withdraw:              currencycmds.NewWithdrawCommand(),
		CreateDAO:             NewCreateDAOCommand(),
		Propose:               NewProposeCommand(),
		Register:              NewRegisterCommand(),
		CurrencyRegister:      currencycmds.NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater: currencycmds.NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:     currencycmds.NewSuffrageInflationCommand(),
		SuffrageCandidate:     currencycmds.NewSuffrageCandidateCommand(),
		SuffrageJoin:          currencycmds.NewSuffrageJoinCommand(),
		SuffrageDisjoin:       currencycmds.NewSuffrageDisjoinCommand(),
	}
}
