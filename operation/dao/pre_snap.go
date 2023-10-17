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
	PreSnapFactHint = hint.MustNewHint("mitum-dao-pre-snap-operation-fact-v0.0.1")
	PreSnapHint     = hint.MustNewHint("mitum-dao-pre-snap-operation-v0.0.1")
)

type PreSnapFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	proposalID string
	currency   currencytypes.CurrencyID
}

func NewPreSnapFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	proposalID string,
	currency currencytypes.CurrencyID,
) PreSnapFact {
	bf := base.NewBaseFact(PreSnapFactHint, token)
	fact := PreSnapFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		proposalID: proposalID,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact PreSnapFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact PreSnapFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact PreSnapFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.proposalID),
		fact.currency.Bytes(),
	)
}

func (fact PreSnapFact) IsValid(b []byte) error {
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

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact PreSnapFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact PreSnapFact) Sender() base.Address {
	return fact.sender
}

func (fact PreSnapFact) Contract() base.Address {
	return fact.contract
}

func (fact PreSnapFact) ProposalID() string {
	return fact.proposalID
}

func (fact PreSnapFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact PreSnapFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type PreSnap struct {
	common.BaseOperation
}

func NewPreSnap(fact PreSnapFact) (PreSnap, error) {
	return PreSnap{BaseOperation: common.NewBaseOperation(PreSnapHint, fact)}, nil
}

func (op *PreSnap) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
