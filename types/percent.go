package types

import "github.com/ProtoconNet/mitum2/util"

type PercentRatio uint

func (r PercentRatio) IsValid([]byte) error {
	if 100 < r {
		return util.ErrInvalid.Errorf("percent ratio out of range, %d", r)
	}

	return nil
}

func (r PercentRatio) Bytes() []byte {
	return util.UintToBytes(uint(r))
}
