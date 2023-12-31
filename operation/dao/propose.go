package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	ProposeFactHint = hint.MustNewHint("mitum-dao-propose-operation-fact-v0.0.1")
	ProposeHint     = hint.MustNewHint("mitum-dao-propose-operation-v0.0.1")
)

type ProposeFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	proposalID string
	proposal   types.Proposal
	currency   currencytypes.CurrencyID
}

func NewProposeFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	proposalID string,
	proposal types.Proposal,
	currency currencytypes.CurrencyID,
) ProposeFact {
	bf := base.NewBaseFact(ProposeFactHint, token)
	fact := ProposeFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		proposalID: proposalID,
		proposal:   proposal,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ProposeFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ProposeFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ProposeFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.proposalID),
		fact.proposal.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact ProposeFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.contract,
		fact.proposal,
		fact.currency,
	); err != nil {
		return err
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	if !fact.sender.Equal(fact.Proposal().Proposer()) {
		return util.ErrInvalid.Errorf("sender is not same with the proposer")
	}

	if len(fact.proposalID) == 0 {
		return util.ErrInvalid.Errorf("empty propose id")
	}

	if !currencytypes.ReSpcecialChar.Match([]byte(fact.proposalID)) {
		return util.ErrInvalid.Errorf("invalid proposalID due to the inclusion of special characters")
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact ProposeFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ProposeFact) Sender() base.Address {
	return fact.sender
}

func (fact ProposeFact) Contract() base.Address {
	return fact.contract
}

func (fact ProposeFact) ProposalID() string {
	return fact.proposalID
}

func (fact ProposeFact) Proposal() types.Proposal {
	return fact.proposal
}

func (fact ProposeFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact ProposeFact) Addresses() ([]base.Address, error) {
	as := fact.proposal.Addresses()

	as = append(as, fact.sender)
	as = append(as, fact.contract)

	return as, nil
}

type Propose struct {
	common.BaseOperation
}

func NewPropose(fact ProposeFact) (Propose, error) {
	return Propose{BaseOperation: common.NewBaseOperation(ProposeHint, fact)}, nil
}

func (op *Propose) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
