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
	SnapFactHint = hint.MustNewHint("mitum-dao-snap-operation-fact-v0.0.1")
	SnapHint     = hint.MustNewHint("mitum-dao-snap-operation-v0.0.1")
)

type SnapFact struct {
	base.BaseFact
	sender    base.Address
	contract  base.Address
	daoID     currencytypes.ContractID
	proposeID string
	currency  currencytypes.CurrencyID
}

func NewSnapFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	daoID currencytypes.ContractID,
	proposeID string,
	currency currencytypes.CurrencyID,
) SnapFact {
	bf := base.NewBaseFact(SnapFactHint, token)
	fact := SnapFact{
		BaseFact:  bf,
		sender:    sender,
		contract:  contract,
		daoID:     daoID,
		proposeID: proposeID,
		currency:  currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact SnapFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact SnapFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SnapFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.daoID.Bytes(),
		[]byte(fact.proposeID),
		fact.currency.Bytes(),
	)
}

func (fact SnapFact) IsValid(b []byte) error {
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

	return nil
}

func (fact SnapFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact SnapFact) Sender() base.Address {
	return fact.sender
}

func (fact SnapFact) Contract() base.Address {
	return fact.contract
}

func (fact SnapFact) DAOID() currencytypes.ContractID {
	return fact.daoID
}

func (fact SnapFact) ProposeID() string {
	return fact.proposeID
}

func (fact SnapFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact SnapFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2)

	as[0] = fact.sender
	as[1] = fact.contract

	return as, nil
}

type Snap struct {
	common.BaseOperation
}

func NewSnap(fact SnapFact) (Snap, error) {
	return Snap{BaseOperation: common.NewBaseOperation(SnapHint, fact)}, nil
}

func (op *Snap) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
