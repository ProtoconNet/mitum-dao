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
	VoteFactHint = hint.MustNewHint("mitum-dao-vote-operation-fact-v0.0.1")
	VoteHint     = hint.MustNewHint("mitum-dao-vote-operation-v0.0.1")
)

type VoteFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	proposalID string
	vote       uint8
	currency   currencytypes.CurrencyID
}

func NewVoteFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	proposalID string,
	vote uint8,
	currency currencytypes.CurrencyID,
) VoteFact {
	bf := base.NewBaseFact(VoteFactHint, token)
	fact := VoteFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		proposalID: proposalID,
		vote:       vote,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact VoteFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact VoteFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact VoteFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.proposalID),
		util.Uint8ToBytes(fact.vote),
		fact.currency.Bytes(),
	)
}

func (fact VoteFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
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

	if !currencytypes.ReSpcecialChar.Match([]byte(fact.proposalID)) {
		return util.ErrInvalid.Errorf("invalid proposalID due to the inclusion of special characters")
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact VoteFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact VoteFact) Sender() base.Address {
	return fact.sender
}

func (fact VoteFact) Contract() base.Address {
	return fact.contract
}

func (fact VoteFact) ProposalID() string {
	return fact.proposalID
}

func (fact VoteFact) Vote() uint8 {
	return fact.vote
}

func (fact VoteFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact VoteFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type Vote struct {
	common.BaseOperation
}

func NewVote(fact VoteFact) (Vote, error) {
	return Vote{BaseOperation: common.NewBaseOperation(VoteHint, fact)}, nil
}

func (op *Vote) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
