package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	SnapHistoryHint = hint.MustNewHint("mitum-dao-snap-history-v0.0.1")
)

type SnapHistory struct {
	hint.BaseHinter
	timestamp uint64
	snaps     []VotingPower
}

func NewSnapHistory(timestamp uint64, snaps []VotingPower) SnapHistory {
	return SnapHistory{
		BaseHinter: hint.NewBaseHinter(SnapHistoryHint),
		timestamp:  timestamp,
		snaps:      snaps,
	}
}

func (sh SnapHistory) Hint() hint.Hint {
	return sh.BaseHinter.Hint()
}

func (sh SnapHistory) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid SnapHistory")

	if err := sh.BaseHinter.IsValid(SnapHistoryHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, snap := range sh.snaps {
		if err := snap.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if _, found := founds[snap.account.String()]; found {
			return e.Wrap(errors.Errorf("duplicate snap account found, %q", snap.account))
		}

		founds[snap.account.String()] = struct{}{}
	}

	return nil
}

func (sh SnapHistory) Bytes() []byte {
	bs := make([][]byte, len(sh.snaps))

	for i, snap := range sh.snaps {
		bs[i] = snap.Bytes()
	}

	return util.ConcatBytesSlice(
		util.Uint64ToBytes(sh.timestamp),
		util.ConcatBytesSlice(bs...),
	)
}

func (sh SnapHistory) TimeStamp() uint64 {
	return sh.timestamp
}

func (sh SnapHistory) Snaps() []VotingPower {
	return sh.snaps
}
