package types

import (
	"net/url"
	"strings"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type URL string

func (u URL) IsValid([]byte) error {
	if _, err := url.Parse(string(u)); err != nil {
		return err
	}

	if u != "" && strings.TrimSpace(string(u)) == "" {
		return util.ErrInvalid.Errorf("empty url")
	}

	return nil
}

func (u URL) Bytes() []byte {
	return []byte(u)
}

func (u URL) String() string {
	return string(u)
}

const (
	ProposalCrypto = "crypto"
	ProposalBiz    = "biz"
)

var (
	CryptoProposalHint = hint.MustNewHint("mitum-dao-crypto-proposal-v0.0.1")
	BizProposalHint    = hint.MustNewHint("mitum-dao-biz-proposal-v0.0.1")
)

type Proposal interface {
	util.IsValider
	hint.Hinter
	Type() string
	Bytes() []byte
	Proposer() base.Address
	StartTime() uint64
	Options() uint8
	Addresses() []base.Address
}

type CryptoProposal struct {
	hint.BaseHinter
	proposer  base.Address
	startTime uint64
	callData  CallData
}

func NewCryptoProposal(proposer base.Address, startTime uint64, callData CallData) CryptoProposal {
	return CryptoProposal{
		BaseHinter: hint.NewBaseHinter(CryptoProposalHint),
		proposer:   proposer,
		startTime:  startTime,
		callData:   callData,
	}
}

func (CryptoProposal) Type() string {
	return ProposalCrypto
}

func (CryptoProposal) Options() uint8 {
	return 3
}

func (p CryptoProposal) Bytes() []byte {
	return util.ConcatBytesSlice(
		p.proposer.Bytes(),
		util.Uint64ToBytes(p.startTime),
		p.callData.Bytes(),
	)
}

func (p CryptoProposal) Proposer() base.Address {
	return p.proposer
}

func (p CryptoProposal) StartTime() uint64 {
	return p.startTime
}

func (p CryptoProposal) CallData() CallData {
	return p.callData
}

func (p CryptoProposal) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		p.BaseHinter,
		p.proposer,
		p.callData,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid CryptoProposal: %v", err)
	}

	return nil
}

func (p CryptoProposal) Addresses() []base.Address {
	return p.callData.Addresses()
}

type BizProposal struct {
	hint.BaseHinter
	proposer  base.Address
	startTime uint64
	url       URL
	hash      string
	options   uint8
}

func NewBizProposal(proposer base.Address, startTime uint64, url URL, hash string, options uint8) BizProposal {
	return BizProposal{
		BaseHinter: hint.NewBaseHinter(BizProposalHint),
		proposer:   proposer,
		startTime:  startTime,
		url:        url,
		hash:       hash,
		options:    options,
	}
}

func (BizProposal) Type() string {
	return ProposalBiz
}

func (p BizProposal) Options() uint8 {
	return p.options
}

func (p BizProposal) Bytes() []byte {
	return util.ConcatBytesSlice(
		p.proposer.Bytes(),
		util.Uint64ToBytes(p.startTime),
		p.url.Bytes(),
		[]byte(p.hash),
		util.Uint8ToBytes(p.options),
	)
}

func (p BizProposal) Proposer() base.Address {
	return p.proposer
}

func (p BizProposal) StartTime() uint64 {
	return p.startTime
}

func (p BizProposal) Url() URL {
	return p.url
}

func (p BizProposal) Hash() string {
	return p.hash
}

func (p BizProposal) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		p.BaseHinter,
		p.proposer,
		p.url,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid BizProposal: %v", err)
	}

	if len(p.hash) == 0 {
		return util.ErrInvalid.Errorf("biz - empty hash")
	}

	return nil
}

func (p BizProposal) Addresses() []base.Address {
	return []base.Address{}
}

func GetPeriodOfCurrentTime(
	policy Policy,
	proposal Proposal,
	blockmap base.BlockMap,
) (Period, int64 /*period start time*/, int64 /*period end time*/) {
	blockTime := uint64(blockmap.Manifest().ProposedAt().Unix())
	startTime := proposal.StartTime()
	registrationTime := startTime + policy.ProposalReviewPeriod()
	preSnapTime := registrationTime + policy.RegistrationPeriod()
	votingTime := preSnapTime + policy.PreSnapshotPeriod()
	postSnapTime := votingTime + policy.VotingPeriod()
	executionDelayTime := postSnapTime + policy.PostSnapshotPeriod()
	executeTime := executionDelayTime + policy.ExecutionDelayPeriod()

	switch {
	case blockTime < startTime:
		return PreLifeCycle, 0, int64(startTime)
	case blockTime < registrationTime:
		return ProposalReview, int64(startTime), int64(registrationTime)
	case blockTime < preSnapTime:
		return Registration, int64(registrationTime), int64(preSnapTime)
	case blockTime < votingTime:
		return PreSnapshot, int64(preSnapTime), int64(votingTime)
	case blockTime < postSnapTime:
		return Voting, int64(votingTime), int64(postSnapTime)
	case blockTime < executionDelayTime:
		return PostSnapshot, int64(postSnapTime), int64(executionDelayTime)
	case blockTime < executeTime:
		return ExecutionDelay, int64(executionDelayTime), int64(executeTime)
	case blockTime >= executeTime:
		return Execute, int64(executeTime), 0
	}

	return NilPeriod, 0, 0
}

func GetPeriod(
	period Period,
	policy Policy,
	proposal Proposal,
) (int64, int64) {
	startTime := proposal.StartTime()
	registrationTime := startTime + policy.ProposalReviewPeriod()
	preSnapTime := registrationTime + policy.RegistrationPeriod()
	votingTime := preSnapTime + policy.PreSnapshotPeriod()
	postSnapTime := votingTime + policy.VotingPeriod()
	executionDelayTime := postSnapTime + policy.PostSnapshotPeriod()
	executeTime := executionDelayTime + policy.ExecutionDelayPeriod()

	switch period {
	case PreLifeCycle:
		return 0, int64(startTime)
	case ProposalReview:
		return int64(startTime), int64(registrationTime)
	case Registration:
		return int64(registrationTime), int64(preSnapTime)
	case PreSnapshot:
		return int64(preSnapTime), int64(votingTime)
	case Voting:
		return int64(votingTime), int64(postSnapTime)
	case PostSnapshot:
		return int64(postSnapTime), int64(executionDelayTime)
	case ExecutionDelay:
		return int64(executionDelayTime), int64(executeTime)
	case Execute:
		return int64(executeTime), 0
	}

	return 0, 0
}
