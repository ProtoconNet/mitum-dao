package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
	"strconv"
)

func (vp *VotingPowerBox) unpack(enc encoder.Encoder, ht hint.Hint, st string, bvp []byte, bre []byte) error {
	e := util.StringError("failed to decode bson of VotingPowerBox")

	vp.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(st)
	if err != nil {
		return e.Wrap(err)
	}
	vp.total = big

	votingPowers := make(map[base.Address]common.Big)
	m, err := enc.DecodeMap(bvp)
	if err != nil {
		return err
	}
	for k := range m {
		v, ok := m[k].(common.Big)
		if !ok {
			return errors.Errorf("expected common.Big, not %T", m[k])
		}
		switch ad, err := base.DecodeAddress(k, enc); {
		case err != nil:
			return e.Wrap(err)
		default:
			votingPowers[ad] = v
		}
	}
	vp.votingPowers = votingPowers

	result := make(map[uint8]common.Big)
	m, err = enc.DecodeMap(bre)
	if err != nil {
		return err
	}
	for k := range m {
		v, ok := m[k].(common.Big)
		if !ok {
			return errors.Errorf("expected common.Big, not %T", m[k])
		}
		val, err := strconv.ParseUint(k, 10, 8)
		if err != nil {
			return err
		}
		result[uint8(val)] = v
	}
	vp.result = result

	return nil
}
