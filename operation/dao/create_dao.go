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
	sender               base.Address
	contract             base.Address
	daoID                currencytypes.ContractID
	option               types.DAOOption
	votingPowerToken     currencytypes.CurrencyID
	threshold            currencytypes.Amount
	fee                  currencytypes.Amount
	whitelist            types.Whitelist
	proposalReviewPeriod uint64
	registrationPeriod   uint64
	preSnapshotPeriod    uint64
	votingPeriod         uint64
	postSnapshotPeriod   uint64
	executionDelayPeriod uint64
	turnout              types.PercentRatio
	quorum               types.PercentRatio
	currency             currencytypes.CurrencyID
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
	proposalReviewPeriod,
	registrationPeriod,
	preSnapshotPeriod,
	votingPeriod,
	postSnapshotPeriod,
	executionDelayPeriod uint64,
	turnout, quorum types.PercentRatio,
	currency currencytypes.CurrencyID,
) CreateDAOFact {
	bf := base.NewBaseFact(CreateDAOFactHint, token)
	fact := CreateDAOFact{
		BaseFact:             bf,
		sender:               sender,
		contract:             contract,
		daoID:                daoID,
		option:               option,
		votingPowerToken:     votingPowerToken,
		threshold:            threshold,
		fee:                  fee,
		whitelist:            whitelist,
		proposalReviewPeriod: proposalReviewPeriod,
		registrationPeriod:   registrationPeriod,
		preSnapshotPeriod:    preSnapshotPeriod,
		votingPeriod:         votingPeriod,
		executionDelayPeriod: executionDelayPeriod,
		postSnapshotPeriod:   postSnapshotPeriod,
		turnout:              turnout,
		quorum:               quorum,
		currency:             currency,
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
		fact.option.Bytes(),
		fact.votingPowerToken.Bytes(),
		fact.fee.Bytes(),
		fact.threshold.Bytes(),
		fact.whitelist.Bytes(),
		util.Uint64ToBytes(fact.proposalReviewPeriod),
		util.Uint64ToBytes(fact.registrationPeriod),
		util.Uint64ToBytes(fact.preSnapshotPeriod),
		util.Uint64ToBytes(fact.votingPeriod),
		util.Uint64ToBytes(fact.postSnapshotPeriod),
		util.Uint64ToBytes(fact.executionDelayPeriod),
		fact.turnout.Bytes(),
		fact.quorum.Bytes(),
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
		fact.turnout,
		fact.quorum,
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

func (fact CreateDAOFact) Option() types.DAOOption {
	return fact.option
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

func (fact CreateDAOFact) ProposalReviewPeriod() uint64 {
	return fact.proposalReviewPeriod
}

func (fact CreateDAOFact) RegistrationPeriod() uint64 {
	return fact.registrationPeriod
}

func (fact CreateDAOFact) PreSnapshotPeriod() uint64 {
	return fact.preSnapshotPeriod
}

func (fact CreateDAOFact) VotingPeriod() uint64 {
	return fact.votingPeriod
}

func (fact CreateDAOFact) PostSnapshotPeriod() uint64 {
	return fact.postSnapshotPeriod
}

func (fact CreateDAOFact) ExecutionDelayPeriod() uint64 {
	return fact.executionDelayPeriod
}

func (fact CreateDAOFact) Turnout() types.PercentRatio {
	return fact.turnout
}

func (fact CreateDAOFact) Quorum() types.PercentRatio {
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
