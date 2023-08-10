package types

import (
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-dao/utils"
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
		return e.Wrap(err)
	}
	for k, v := range m {
		vp, ok := v.(VotingPower)
		if !ok {
			return e.Wrap(errors.Errorf("expected VotingPower, not %T", v))
		}

		if _, err := base.DecodeAddress(k, enc); err != nil {
			return e.Wrap(err)
		}

		votingPowers[k] = vp
	}
	vp.votingPowers = votingPowers

	m, err = utils.DecodeMap(bre)
	if err != nil {
		return e.Wrap(err)
	}

	result := make(map[uint8]common.Big)
	for k, v := range m {
		u, err := strconv.ParseUint(k, 10, 8)
		if err != nil {
			return e.Wrap(err)
		}

		big, err := common.NewBigFromInterface(v)
		if err != nil {
			return e.Wrap(err)
		}

		result[uint8(u)] = big
	}
	vp.result = result

	return nil
}
