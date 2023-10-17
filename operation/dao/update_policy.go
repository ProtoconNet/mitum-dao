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
	UpdatePolicyFactHint = hint.MustNewHint("mitum-dao-update-policy-operation-fact-v0.0.1")
	UpdatePolicyHint     = hint.MustNewHint("mitum-dao-update-policy-operation-v0.0.1")
)

type UpdatePolicyFact struct {
	base.BaseFact
	sender               base.Address
	contract             base.Address
	option               types.DAOOption
	votingPowerToken     currencytypes.CurrencyID
	threshold            common.Big
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

func NewUpdatePolicyFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	option types.DAOOption,
	votingPowerToken currencytypes.CurrencyID,
	threshold common.Big,
	fee currencytypes.Amount,
	whitelist types.Whitelist,
	proposalReviewPeriod,
	registrationPeriod,
	preSnapshotPeriod,
	votingPeriod,
	postSnapshotPeriod,
	executionDelayPeriod uint64,
	turnout, quorum types.PercentRatio,
	currency currencytypes.CurrencyID,
) UpdatePolicyFact {
	bf := base.NewBaseFact(UpdatePolicyFactHint, token)
	fact := UpdatePolicyFact{
		BaseFact:             bf,
		sender:               sender,
		contract:             contract,
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

func (fact UpdatePolicyFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact UpdatePolicyFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UpdatePolicyFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.option.Bytes(),
		fact.votingPowerToken.Bytes(),
		fact.threshold.Bytes(),
		fact.fee.Bytes(),
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

func (fact UpdatePolicyFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
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

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	return nil
}

func (fact UpdatePolicyFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact UpdatePolicyFact) Sender() base.Address {
	return fact.sender
}

func (fact UpdatePolicyFact) Contract() base.Address {
	return fact.contract
}

func (fact UpdatePolicyFact) Option() types.DAOOption {
	return fact.option
}

func (fact UpdatePolicyFact) VotingPowerToken() currencytypes.CurrencyID {
	return fact.votingPowerToken
}

func (fact UpdatePolicyFact) Fee() currencytypes.Amount {
	return fact.fee
}

func (fact UpdatePolicyFact) Threshold() common.Big {
	return fact.threshold
}

func (fact UpdatePolicyFact) Whitelist() types.Whitelist {
	return fact.whitelist
}

func (fact UpdatePolicyFact) ProposalReviewPeriod() uint64 {
	return fact.proposalReviewPeriod
}

func (fact UpdatePolicyFact) RegistrationPeriod() uint64 {
	return fact.registrationPeriod
}

func (fact UpdatePolicyFact) PreSnapshotPeriod() uint64 {
	return fact.preSnapshotPeriod
}

func (fact UpdatePolicyFact) VotingPeriod() uint64 {
	return fact.votingPeriod
}

func (fact UpdatePolicyFact) PostSnapshotPeriod() uint64 {
	return fact.postSnapshotPeriod
}

func (fact UpdatePolicyFact) ExecutionDelayPeriod() uint64 {
	return fact.executionDelayPeriod
}

func (fact UpdatePolicyFact) Turnout() types.PercentRatio {
	return fact.turnout
}

func (fact UpdatePolicyFact) Quorum() types.PercentRatio {
	return fact.quorum
}

func (fact UpdatePolicyFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact UpdatePolicyFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 2+len(fact.whitelist.Accounts()))

	as[0] = fact.sender
	as[1] = fact.contract

	for i, ac := range fact.whitelist.Accounts() {
		as[i+2] = ac
	}

	return as, nil
}

type UpdatePolicy struct {
	common.BaseOperation
}

func NewUpdatePolicy(fact UpdatePolicyFact) (UpdatePolicy, error) {
	return UpdatePolicy{BaseOperation: common.NewBaseOperation(UpdatePolicyHint, fact)}, nil
}

func (op *UpdatePolicy) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
