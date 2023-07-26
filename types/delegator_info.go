package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var DelegatorInfoHint = hint.MustNewHint("mitum-dao-delegator-info-v0.0.1")

type DelegatorInfo struct {
	hint.BaseHinter
	account   base.Address
	delegatee base.Address
}

func NewDelegatorInfo(account base.Address, delegatee base.Address) DelegatorInfo {
	return DelegatorInfo{
		BaseHinter: hint.NewBaseHinter(DelegatorInfoHint),
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
