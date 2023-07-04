package types

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

//var (
//	VotingPowerHint = hint.MustNewHint("mitum-dao-voting-power-v0.0.1")
//)

//	type VotingPower struct {
//		hint.BaseHinter
//		account base.Address
//		amount  common.Big
//	}
//
//	func NewVotingPower(account base.Address, votingPower common.Big) VotingPower {
//		return VotingPower{
//			BaseHinter: hint.NewBaseHinter(VotingPowerHint),
//			account:    account,
//			amount:     votingPower,
//		}
//	}
//
//	func (vp VotingPower) Hint() hint.Hint {
//		return vp.BaseHinter.Hint()
//	}
//
//	func (vp VotingPower) IsValid([]byte) error {
//		e := util.ErrInvalid.Errorf("invalid Amount")
//
//		if err := vp.BaseHinter.IsValid(VotingPowerHint.Type().Bytes()); err != nil {
//			return e.Wrap(err)
//		}
//
//		if err := util.CheckIsValiders(nil, false, vp.account, vp.amount); err != nil {
//			return e.Wrap(err)
//		}
//
//		return nil
//	}
//
//	func (vp VotingPower) Bytes() []byte {
//		return util.ConcatBytesSlice(
//			vp.account.Bytes(),
//			vp.amount.Bytes(),
//		)
//	}
//
//	func (vp VotingPower) Account() base.Address {
//		return vp.account
//	}
//
//	func (vp VotingPower) Amount() common.Big {
//		return vp.amount
//	}

var (
	VotingPowerBoxHint = hint.MustNewHint("mitum-dao-voting-power-box-v0.0.1")
)

type VotingPowerBox struct {
	hint.BaseHinter
	total        common.Big
	votingPowers map[base.Address]common.Big
	result       map[uint8]common.Big
}

func NewVotingPowerBox(total common.Big, votingPowers map[base.Address]common.Big) VotingPowerBox {
	return VotingPowerBox{
		BaseHinter:   hint.NewBaseHinter(VotingPowerBoxHint),
		total:        total,
		votingPowers: votingPowers,
		result:       map[uint8]common.Big{},
	}
}

func (vp VotingPowerBox) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotingPowerBox")

	if err := vp.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	total := common.ZeroBig
	for _, vp := range vp.votingPowers {
		if err := vp.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		total = total.Add(vp)
	}

	if total.Compare(vp.total) != 0 {
		return e.Wrap(errors.Errorf("invalid voting power total, %q != %q", total, vp.total))
	}

	return nil
}

func (vp VotingPowerBox) Bytes() []byte {
	bs := make([][]byte, 3)
	bs[0] = vp.total.Bytes()
	if vp.votingPowers != nil {
		votingPowers, _ := json.Marshal(vp.votingPowers)
		bs[1] = valuehash.NewSHA256(votingPowers).Bytes()
	} else {
		bs[1] = []byte{}
	}

	if vp.result != nil {
		result, _ := json.Marshal(vp.result)
		bs[2] = valuehash.NewSHA256(result).Bytes()
	} else {
		bs[2] = []byte{}
	}

	return util.ConcatBytesSlice(bs...)
}

func (vp VotingPowerBox) Total() common.Big {
	return vp.total
}

func (vp VotingPowerBox) VotingPowers() map[base.Address]common.Big {
	return vp.votingPowers
}

func (vp VotingPowerBox) Result() map[uint8]common.Big {
	return vp.result
}

func (vp *VotingPowerBox) SetTotal(total common.Big) {
	vp.total = total
}

func (vp *VotingPowerBox) SetVotingPowers(votingPowers map[base.Address]common.Big) {
	vp.votingPowers = votingPowers
}

func (vp *VotingPowerBox) SetResult(result map[uint8]common.Big) {
	vp.result = result
}
