package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/util"
)

type PercentRatio uint8

func (r PercentRatio) IsValid([]byte) error {
	if 100 < r {
		return util.ErrInvalid.Errorf("percent ratio out of range, %d", r)
	}

	return nil
}

func (r PercentRatio) Bytes() []byte {
	return util.Uint8ToBytes(uint8(r))
}

func (r PercentRatio) Quorum(total common.Big) common.Big {
	if !total.OverZero() || r == 0 {
		return common.ZeroBig
	}

	if r == 100 {
		return total
	}

	return total.Mul(common.NewBig(int64(r))).Div(common.NewBig(100))
}
