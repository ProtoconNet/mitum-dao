package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	VotingPowerHint = hint.MustNewHint("mitum-dao-voting-power-v0.0.1")
)

type VotingPower struct {
	hint.BaseHinter
	account     base.Address
	votingPower common.Big
}

func NewVotingPower(account base.Address, votingPower common.Big) VotingPower {
	return VotingPower{
		BaseHinter:  hint.NewBaseHinter(VotingPowerHint),
		account:     account,
		votingPower: votingPower,
	}
}

func (vp VotingPower) Hint() hint.Hint {
	return vp.BaseHinter.Hint()
}

func (vp VotingPower) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotingPower")

	if err := vp.BaseHinter.IsValid(SnapHistoryHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false, vp.account, vp.votingPower); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (vp VotingPower) Bytes() []byte {
	return util.ConcatBytesSlice(
		vp.account.Bytes(),
		vp.votingPower.Bytes(),
	)
}

func (vp VotingPower) Account() base.Address {
	return vp.account
}

func (vp VotingPower) VotingPower() common.Big {
	return vp.votingPower
}

var (
	VotingPowersHint = hint.MustNewHint("mitum-dao-voting-powers-v0.0.1")
)

type VotingPowers struct {
	hint.BaseHinter
	total        common.Big
	votingPowers []VotingPower
}

func NewVotingPowers(total common.Big, votingPowers []VotingPower) VotingPowers {
	return VotingPowers{
		BaseHinter:   hint.NewBaseHinter(VotingPowersHint),
		total:        total,
		votingPowers: votingPowers,
	}
}

func (vp VotingPowers) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotingPowers")

	if err := vp.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	total := common.ZeroBig
	for _, vp := range vp.votingPowers {
		if err := vp.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		total = total.Add(vp.votingPower)
	}

	if total.Compare(vp.total) != 0 {
		return e.Wrap(errors.Errorf("invalid voting power total, %q != %q", total, vp.total))
	}

	return nil
}

func (vp VotingPowers) Bytes() []byte {
	bs := make([][]byte, len(vp.votingPowers))
	for i, v := range vp.votingPowers {
		bs[i] = v.Bytes()
	}

	return util.ConcatBytesSlice(
		vp.total.Bytes(),
		util.ConcatBytesSlice(bs...),
	)
}

func (vp VotingPowers) Total() common.Big {
	return vp.total
}

func (vp VotingPowers) VotingPowers() []VotingPower {
	return vp.votingPowers
}
