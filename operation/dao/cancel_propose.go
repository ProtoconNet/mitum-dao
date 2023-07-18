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
	CancelProposeFactHint = hint.MustNewHint("mitum-dao-cancel-propose-operation-fact-v0.0.1")
	CancelProposeHint     = hint.MustNewHint("mitum-dao-cancel-propose-operation-v0.0.1")
)

type CancelProposeFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	daoID      currencytypes.ContractID
	proposalID string
	currency   currencytypes.CurrencyID
}

func NewCancelProposeFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	daoID currencytypes.ContractID,
	proposalID string,
	currency currencytypes.CurrencyID,
) CancelProposeFact {
	bf := base.NewBaseFact(CancelProposeFactHint, token)
	fact := CancelProposeFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		daoID:      daoID,
		proposalID: proposalID,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CancelProposeFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CancelProposeFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CancelProposeFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.daoID.Bytes(),
		[]byte(fact.proposalID),
		fact.currency.Bytes(),
	)
}

func (fact CancelProposeFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.daoID,
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

func (fact CancelProposeFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CancelProposeFact) Sender() base.Address {
	return fact.sender
}

func (fact CancelProposeFact) Contract() base.Address {
	return fact.contract
}

func (fact CancelProposeFact) DAOID() currencytypes.ContractID {
	return fact.daoID
}

func (fact CancelProposeFact) ProposalID() string {
	return fact.proposalID
}

func (fact CancelProposeFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CancelProposeFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type CancelPropose struct {
	common.BaseOperation
}

func NewCancelPropose(fact CancelProposeFact) (CancelPropose, error) {
	return CancelPropose{BaseOperation: common.NewBaseOperation(CancelProposeHint, fact)}, nil
}

func (op *CancelPropose) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
