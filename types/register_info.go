package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var RegisterInfoHint = hint.MustNewHint("mitum-dao-register-info-v0.0.1")

type RegisterInfo struct {
	hint.BaseHinter
	account    base.Address
	approvedBy []base.Address
}

func NewRegisterInfo(account base.Address, approvedBy []base.Address) RegisterInfo {
	return RegisterInfo{
		BaseHinter: hint.NewBaseHinter(RegisterInfoHint),
		account:    account,
		approvedBy: approvedBy,
	}
}

func (r RegisterInfo) Hint() hint.Hint {
	return r.BaseHinter.Hint()
}

func (r RegisterInfo) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid RegisterInfo")

	if err := r.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if err := r.account.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	for _, ac := range r.approvedBy {
		if err := ac.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if ac.Equal(r.account) {
			return e.Wrap(errors.Errorf("approving address is same with approved address, %q", r.Account))
		}
	}

	return nil
}

func (r RegisterInfo) Bytes() []byte {
	ba := make([][]byte, len(r.approvedBy)+1)

	ba[0] = r.account.Bytes()

	for i, ac := range r.approvedBy {
		ba[i+1] = ac.Bytes()
	}

	return util.ConcatBytesSlice(ba...)
}

func (r RegisterInfo) Account() base.Address {
	return r.account
}

func (r RegisterInfo) ApprovedBy() []base.Address {
	return r.approvedBy
}
