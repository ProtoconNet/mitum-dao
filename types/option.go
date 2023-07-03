package types

import "github.com/ProtoconNet/mitum2/util"

type Option uint8

type ProposalStatus Option

func (p ProposalStatus) Bytes() []byte {
	return util.Uint8ToBytes(uint8(p))
}

const (
	Proposed ProposalStatus = iota
	Canceled
	PreSnapped
	PostSnapped
	Completed
	Rejected
	Executed
	NilStatus
)

type Period Option

func (p Period) Bytes() []byte {
	return util.Uint8ToBytes(uint8(p))
}

const (
	PreLifeCycle Period = iota
	ProposalReview
	PreSnapshot
	Registration
	Voting
	PostSnapshot
	ExecutionDelay
	Execute
	NilPeriod
)
