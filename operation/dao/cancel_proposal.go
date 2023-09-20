package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CancelProposalFactHint = hint.MustNewHint("mitum-dao-cancel-proposal-operation-fact-v0.0.1")
	CancelProposalHint     = hint.MustNewHint("mitum-dao-cancel-proposal-operation-v0.0.1")
)

type CancelProposalFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	proposalID string
	currency   currencytypes.CurrencyID
}

func NewCancelProposalFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	proposalID string,
	currency currencytypes.CurrencyID,
) CancelProposalFact {
	bf := base.NewBaseFact(CancelProposalFactHint, token)
	fact := CancelProposalFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		proposalID: proposalID,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CancelProposalFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CancelProposalFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CancelProposalFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.proposalID),
		fact.currency.Bytes(),
	)
}

func (fact CancelProposalFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return err
	}

	if len(fact.proposalID) == 0 {
		return util.ErrInvalid.Errorf("empty propose id")
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	return nil
}

func (fact CancelProposalFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CancelProposalFact) Sender() base.Address {
	return fact.sender
}

func (fact CancelProposalFact) Contract() base.Address {
	return fact.contract
}

func (fact CancelProposalFact) ProposalID() string {
	return fact.proposalID
}

func (fact CancelProposalFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CancelProposalFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type CancelProposal struct {
	common.BaseOperation
}

func NewCancelProposal(fact CancelProposalFact) (CancelProposal, error) {
	return CancelProposal{BaseOperation: common.NewBaseOperation(CancelProposalHint, fact)}, nil
}

func (op *CancelProposal) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
