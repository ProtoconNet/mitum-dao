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
	RegisterFactHint = hint.MustNewHint("mitum-dao-register-operation-fact-v0.0.1")
	RegisterHint     = hint.MustNewHint("mitum-dao-register-operation-v0.0.1")
)

type RegisterFact struct {
	base.BaseFact
	sender    base.Address
	contract  base.Address
	daoID     currencytypes.ContractID
	proposeID string
	approved  base.Address
	currency  currencytypes.CurrencyID
}

func NewRegisterFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	daoID currencytypes.ContractID,
	proposeID string,
	approved base.Address,
	currency currencytypes.CurrencyID,
) RegisterFact {
	bf := base.NewBaseFact(RegisterFactHint, token)
	fact := RegisterFact{
		BaseFact:  bf,
		sender:    sender,
		contract:  contract,
		daoID:     daoID,
		proposeID: proposeID,
		approved:  approved,
		currency:  currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RegisterFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RegisterFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RegisterFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.daoID.Bytes(),
		[]byte(fact.proposeID),
		fact.approved.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact RegisterFact) IsValid(b []byte) error {
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

	if len(fact.proposeID) == 0 {
		return util.ErrInvalid.Errorf("empty propose id")
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	if fact.approved != nil {
		if err := fact.approved.IsValid(nil); err != nil {
			return err
		}

		if fact.sender.Equal(fact.approved) {
			return util.ErrInvalid.Errorf("sender cannot approve itself, %q", fact.sender)
		}

		if fact.approved.Equal(fact.contract) {
			return util.ErrInvalid.Errorf("contract address is same with approved, %q", fact.approved)
		}
	}

	return nil
}

func (fact RegisterFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RegisterFact) Sender() base.Address {
	return fact.sender
}

func (fact RegisterFact) Contract() base.Address {
	return fact.contract
}

func (fact RegisterFact) DAOID() currencytypes.ContractID {
	return fact.daoID
}

func (fact RegisterFact) ProposeID() string {
	return fact.proposeID
}

func (fact RegisterFact) Approved() base.Address {
	return fact.approved
}

func (fact RegisterFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact RegisterFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	if fact.approved != nil {
		as = append(as, fact.approved)
	}

	return as, nil
}

type Register struct {
	common.BaseOperation
}

func NewRegister(fact RegisterFact) (Register, error) {
	return Register{BaseOperation: common.NewBaseOperation(RegisterHint, fact)}, nil
}

func (op *Register) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
