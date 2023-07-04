package types

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var WhitelistHint = hint.MustNewHint("mitum-dao-whitelist-v0.0.1")

type Whitelist struct {
	hint.BaseHinter
	active   bool
	accounts []base.Address
}

func NewWhitelist(active bool, accounts []base.Address) Whitelist {
	return Whitelist{
		BaseHinter: hint.NewBaseHinter(WhitelistHint),
		active:     active,
		accounts:   accounts,
	}
}

func (wl Whitelist) Bytes() []byte {
	ab := make([]byte, 1)
	if wl.active {
		ab[0] = 1
	} else {
		ab[0] = 0
	}

	ads := make([][]byte, len(wl.accounts))
	for i := range wl.accounts {
		ads[i] = wl.accounts[i].Bytes()
	}

	return util.ConcatBytesSlice(
		ab,
		util.ConcatBytesSlice(ads...),
	)
}

func (wl Whitelist) IsValid([]byte) error {
	e := util.StringError("invalid whitelist")

	if err := util.CheckIsValiders(nil, false, wl.BaseHinter); err != nil {
		return e.Wrap(err)
	}

	for _, ac := range wl.accounts {
		if err := ac.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}

func (wl Whitelist) Active() bool {
	return wl.active
}

func (wl Whitelist) Accounts() []base.Address {
	return wl.accounts
}

func (wl Whitelist) IsExist(a base.Address) bool {
	for _, ac := range wl.accounts {
		if ac.Equal(a) {
			return true
		}
	}

	return false
}

var PolicyHint = hint.MustNewHint("mitum-dao-policy-v0.0.1")

type Policy struct {
	hint.BaseHinter
	token                currencytypes.CurrencyID
	threshold            currencytypes.Amount
	fee                  currencytypes.Amount
	whitelist            Whitelist
	proposalReviewPeriod uint64
	registrationPeriod   uint64
	preSnapshotPeriod    uint64
	votingPeriod         uint64
	postSnapshotPeriod   uint64
	executionDelayPeriod uint64
	turnout              PercentRatio
	quorum               PercentRatio
}

func NewPolicy(
	token currencytypes.CurrencyID,
	fee, threshold currencytypes.Amount,
	whitelist Whitelist,
	proposalReviewPeriod, registrationPeriod, preSnapshotPeriod, votingPeriod, postSnapshotPeriod, executionDelayPeriod uint64,
	turnout, quorum PercentRatio,
) Policy {
	return Policy{
		BaseHinter:           hint.NewBaseHinter(PolicyHint),
		token:                token,
		fee:                  fee,
		threshold:            threshold,
		whitelist:            whitelist,
		proposalReviewPeriod: proposalReviewPeriod,
		registrationPeriod:   registrationPeriod,
		preSnapshotPeriod:    preSnapshotPeriod,
		votingPeriod:         votingPeriod,
		postSnapshotPeriod:   postSnapshotPeriod,
		executionDelayPeriod: executionDelayPeriod,
		turnout:              turnout,
		quorum:               quorum,
	}
}

func (po Policy) Bytes() []byte {
	return util.ConcatBytesSlice(
		po.token.Bytes(),
		po.fee.Bytes(),
		po.threshold.Bytes(),
		po.whitelist.Bytes(),
		util.Uint64ToBytes(po.proposalReviewPeriod),
		util.Uint64ToBytes(po.registrationPeriod),
		util.Uint64ToBytes(po.preSnapshotPeriod),
		util.Uint64ToBytes(po.votingPeriod),
		util.Uint64ToBytes(po.postSnapshotPeriod),
		util.Uint64ToBytes(po.executionDelayPeriod),
		po.turnout.Bytes(),
		po.quorum.Bytes(),
	)
}

func (po Policy) IsValid([]byte) error {
	e := util.StringError("invalid dao policy")

	if err := util.CheckIsValiders(nil, false,
		po.BaseHinter,
		po.token,
		po.fee,
		po.threshold,
		po.whitelist,
		po.turnout,
		po.quorum,
	); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (po Policy) Token() currencytypes.CurrencyID {
	return po.token
}

func (po Policy) Fee() currencytypes.Amount {
	return po.fee
}

func (po Policy) Threshold() currencytypes.Amount {
	return po.threshold
}

func (po Policy) Whitelist() Whitelist {
	return po.whitelist
}

func (po Policy) ProposalReviewPeriod() uint64 {
	return po.proposalReviewPeriod
}

func (po Policy) RegistrationPeriod() uint64 {
	return po.registrationPeriod
}

func (po Policy) PreSnapshotPeriod() uint64 {
	return po.preSnapshotPeriod
}

func (po Policy) VotingPeriod() uint64 {
	return po.votingPeriod
}

func (po Policy) PostSnapshotPeriod() uint64 {
	return po.postSnapshotPeriod
}

func (po Policy) ExecutionDelayPeriod() uint64 {
	return po.executionDelayPeriod
}

func (po Policy) Turnout() PercentRatio {
	return po.turnout
}

func (po Policy) Quorum() PercentRatio {
	return po.quorum
}
