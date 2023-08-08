package types

import (
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (vp *VotingPowerBox) unpack(enc encoder.Encoder, ht hint.Hint, st string, bvp []byte, bre []byte) error {
	e := util.StringError("failed to unmarshal VotingPowerBox")

	vp.BaseHinter = hint.NewBaseHinter(ht)

	big, err := common.NewBigFromString(st)
	if err != nil {
		return e.Wrap(err)
	}
	vp.total = big

	votingPowers := make(map[string]VotingPower)
	m, err := enc.DecodeMap(bvp)
	if err != nil {
		return err
	}
	for k := range m {
		v, ok := m[k].(VotingPower)
		if !ok {
			return errors.Errorf("expected VotingPower, not %T", m[k])
		}

		if _, err := base.DecodeAddress(k, enc); err != nil {
			return e.Wrap(err)
		}

		votingPowers[k] = v
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
