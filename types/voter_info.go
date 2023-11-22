package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
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

		//if ac.Equal(r.account) {
		//	return e.Wrap(errors.Errorf("approving address is same with approved address, %q", r.Account()))
		//}
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
