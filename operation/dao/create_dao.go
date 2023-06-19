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
	CreateDAOFactHint = hint.MustNewHint("mitum-dao-create-dao-operation-fact-v0.0.1")
	CreateDAOHint     = hint.MustNewHint("mitum-dao-create-dao-operation-v0.0.1")
)

type CreateDAOFact struct {
	base.BaseFact
	sender           base.Address
	contract         base.Address
	daoID            currencytypes.ContractID
	option           types.DAOOption
	votingPowerToken currencytypes.CurrencyID
	threshold        currencytypes.Amount
	fee              currencytypes.Amount
	whitelist        types.Whitelist
	delaytime        uint64
	snaptime         uint64
	timelock         uint64
	turnout          float64
	quorum           float64
	currency         currencytypes.CurrencyID
}

func NewCreateDAOFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	daoID currencytypes.ContractID,
	option types.DAOOption,
	votingPowerToken currencytypes.CurrencyID,
	threshold, fee currencytypes.Amount,
	whitelist types.Whitelist,
	delaytime, snaptime, timelock uint64,
	turnout, quorum float64,
	currency currencytypes.CurrencyID,
) CreateDAOFact {
	bf := base.NewBaseFact(CreateDAOFactHint, token)
	fact := CreateDAOFact{
		BaseFact:         bf,
		sender:           sender,
		contract:         contract,
		daoID:            daoID,
		option:           option,
		votingPowerToken: votingPowerToken,
		threshold:        threshold,
		fee:              fee,
		whitelist:        whitelist,
		delaytime:        delaytime,
		snaptime:         snaptime,
		timelock:         timelock,
		turnout:          turnout,
		quorum:           quorum,
		currency:         currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateDAOFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateDAOFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateDAOFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.daoID.Bytes(),
		fact.votingPowerToken.Bytes(),
		fact.fee.Bytes(),
		fact.threshold.Bytes(),
		fact.whitelist.Bytes(),
		util.Uint64ToBytes(fact.delaytime),
		util.Uint64ToBytes(fact.snaptime),
		util.Uint64ToBytes(fact.timelock),
		util.Float64ToBytes(fact.turnout),
		util.Float64ToBytes(fact.quorum),
		fact.currency.Bytes(),
	)
}

func (fact CreateDAOFact) IsValid(b []byte) error {
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
		fact.votingPowerToken,
		fact.fee,
		fact.threshold,
		fact.whitelist,
		fact.currency,
	); err != nil {
		return err
	}

	if fact.sender.Equal(fact.contract) {
		return util.ErrInvalid.Errorf("contract address is same with sender, %q", fact.sender)
	}

	return nil
}

func (fact CreateDAOFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CreateDAOFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateDAOFact) Contract() base.Address {
	return fact.contract
}

func (fact CreateDAOFact) DAOID() currencytypes.ContractID {
	return fact.daoID
}

func (fact CreateDAOFact) VotingPowerToken() currencytypes.CurrencyID {
	return fact.votingPowerToken
}

func (fact CreateDAOFact) Fee() currencytypes.Amount {
	return fact.fee
}

func (fact CreateDAOFact) Threshold() currencytypes.Amount {
	return fact.threshold
}

func (fact CreateDAOFact) Whitelist() types.Whitelist {
	return fact.whitelist
}

func (fact CreateDAOFact) DelayTime() uint64 {
	return fact.delaytime
}

func (fact CreateDAOFact) SnapTime() uint64 {
	return fact.snaptime
}

func (fact CreateDAOFact) TimeLock() uint64 {
	return fact.timelock
}

func (fact CreateDAOFact) Turnout() float64 {
	return fact.turnout
}

func (fact CreateDAOFact) Quorum() float64 {
	return fact.quorum
}

func (fact CreateDAOFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CreateDAOFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2+len(fact.whitelist.Accounts()))

	as[0] = fact.sender
	as[1] = fact.contract

	for i, ac := range fact.whitelist.Accounts() {
		as[i+2] = ac
	}

	return as, nil
}

type CreateDAO struct {
	common.BaseOperation
}

func NewCreateDAO(fact CreateDAOFact) (CreateDAO, error) {
	return CreateDAO{BaseOperation: common.NewBaseOperation(CreateDAOHint, fact)}, nil
}

func (op *CreateDAO) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
