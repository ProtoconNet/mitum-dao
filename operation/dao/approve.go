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
	ApproveFactHint = hint.MustNewHint("mitum-dao-approve-operation-fact-v0.0.1")
	ApproveHint     = hint.MustNewHint("mitum-dao-approve-operation-v0.0.1")
)

type ApproveFact struct {
	base.BaseFact
	sender    base.Address
	contract  base.Address
	daoID     currencytypes.ContractID
	proposeID string
	target    base.Address
	currency  currencytypes.CurrencyID
}

func NewApproveFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	daoID currencytypes.ContractID,
	proposeID string,
	target base.Address,
	currency currencytypes.CurrencyID,
) ApproveFact {
	bf := base.NewBaseFact(ApproveFactHint, token)
	fact := ApproveFact{
		BaseFact:  bf,
		sender:    sender,
		contract:  contract,
		daoID:     daoID,
		proposeID: proposeID,
		target:    target,
		currency:  currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ApproveFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ApproveFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.daoID.Bytes(),
		[]byte(fact.proposeID),
		fact.target.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact ApproveFact) IsValid(b []byte) error {
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
		fact.target,
		fact.currency,
	); err != nil {
		return err
	}

	if len(fact.proposeID) == 0 {
		return util.ErrInvalid.Errorf("empty propose id")
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	if fact.sender.Equal(fact.target) {
		return util.ErrInvalid.Errorf("sender cannot approve itself, %q", fact.sender)
	}

	if fact.target.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with target, %q", fact.target)
	}

	return nil
}

func (fact ApproveFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ApproveFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveFact) Contract() base.Address {
	return fact.contract
}

func (fact ApproveFact) DAOID() currencytypes.ContractID {
	return fact.daoID
}

func (fact ApproveFact) ProposeID() string {
	return fact.proposeID
}

func (fact ApproveFact) Target() base.Address {
	return fact.target
}

func (fact ApproveFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact ApproveFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 3)

	as[0] = fact.sender
	as[1] = fact.contract
	as[2] = fact.target

	return as, nil
}

type Approve struct {
	common.BaseOperation
}

func NewApprove(fact ApproveFact) (Approve, error) {
	return Approve{BaseOperation: common.NewBaseOperation(ApproveHint, fact)}, nil
}

func (op *Approve) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
