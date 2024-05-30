package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	RegisterFactHint = hint.MustNewHint("mitum-dao-register-operation-fact-v0.0.1")
	RegisterHint     = hint.MustNewHint("mitum-dao-register-operation-v0.0.1")
)

type RegisterFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	proposalID string
	delegated  base.Address
	currency   currencytypes.CurrencyID
}

func NewRegisterFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	proposalID string,
	delegated base.Address,
	currency currencytypes.CurrencyID,
) RegisterFact {
	bf := base.NewBaseFact(RegisterFactHint, token)
	fact := RegisterFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		proposalID: proposalID,
		delegated:  delegated,
		currency:   currency,
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
		[]byte(fact.proposalID),
		fact.delegated.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact RegisterFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.contract,
		fact.currency,
		fact.delegated,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if len(fact.proposalID) == 0 {
		return common.ErrFactInvalid.Wrap(common.ErrValOOR.Wrap(errors.Errorf("empty proposal ID")))
	}

	if !currencytypes.ReSpcecialChar.Match([]byte(fact.proposalID)) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Wrap(
				errors.Errorf("proposal ID %s, must match regex `^[^\\s:/?#\\[\\]@]*$`", fact.proposalID)))
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(
				errors.Errorf("contract address is same with sender, %q", fact.sender)))
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
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

func (fact RegisterFact) ProposalID() string {
	return fact.proposalID
}

func (fact RegisterFact) Delegated() base.Address {
	return fact.delegated
}

func (fact RegisterFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact RegisterFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 3)

	as[0] = fact.sender
	as[1] = fact.contract
	as[2] = fact.delegated

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
