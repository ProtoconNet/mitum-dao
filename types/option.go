package types

import (
	"fmt"
	"github.com/ProtoconNet/mitum2/util"
)

type Option uint8

type ProposalStatus Option

func (p ProposalStatus) Bytes() []byte {
	return util.Uint8ToBytes(uint8(p))
}

func (p ProposalStatus) String() string {
	if name, found := proposalStatusNames[p]; found {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", uint8(p))
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

var proposalStatusNames = map[ProposalStatus]string{
	Proposed:    "proposed",
	Canceled:    "canceled",
	PreSnapped:  "pre-snapped",
	PostSnapped: "post-snapped",
	Completed:   "completed",
	Rejected:    "rejected",
	Executed:    "executed",
	NilStatus:   "none",
}

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
