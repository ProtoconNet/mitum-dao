package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var VoterInfoHint = hint.MustNewHint("mitum-dao-voter-info-v0.0.1")

type VoterInfo struct {
	hint.BaseHinter
	account    base.Address
	delegators []base.Address
}

func NewVoterInfo(account base.Address, delegators []base.Address) VoterInfo {
	return VoterInfo{
		BaseHinter: hint.NewBaseHinter(VoterInfoHint),
		account:    account,
		delegators: delegators,
	}
}

func (r VoterInfo) Hint() hint.Hint {
	return r.BaseHinter.Hint()
}

func (r VoterInfo) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VoterInfo")

	if err := r.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := r.account.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	for _, ac := range r.delegators {
		if err := ac.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if ac.Equal(r.account) {
			return e.Wrap(errors.Errorf("approving address is same with approved address, %q", r.Account))
		}
	}

	return nil
}

func (r VoterInfo) Bytes() []byte {
	ba := make([][]byte, len(r.delegators)+1)

	ba[0] = r.account.Bytes()

	for i, ac := range r.delegators {
		ba[i+1] = ac.Bytes()
	}

	return util.ConcatBytesSlice(ba...)
}

func (r VoterInfo) Account() base.Address {
	return r.account
}

func (r VoterInfo) Delegators() []base.Address {
	return r.delegators
}

var DelegatorInfoHint = hint.MustNewHint("mitum-dao-delegator-info-v0.0.1")

type DelegatorInfo struct {
	hint.BaseHinter
	account   base.Address
	delegatee base.Address
}

func NewDelegatorInfo(account base.Address, delegatee base.Address) DelegatorInfo {
	return DelegatorInfo{
		BaseHinter: hint.NewBaseHinter(VoterInfoHint),
		account:    account,
		delegatee:  delegatee,
	}
}

func (r DelegatorInfo) Hint() hint.Hint {
	return r.BaseHinter.Hint()
}

func (r DelegatorInfo) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VoterInfo")

	if err := r.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := r.account.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := r.delegatee.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (r DelegatorInfo) Bytes() []byte {
	ba := make([][]byte, 2)

	ba[0] = r.account.Bytes()
	ba[1] = r.delegatee.Bytes()

	return util.ConcatBytesSlice(ba...)
}

func (r DelegatorInfo) Account() base.Address {
	return r.account
}

func (r DelegatorInfo) Delegatee() base.Address {
	return r.delegatee
}
