package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

type DAOOption string

func (op DAOOption) IsValid([]byte) error {
	if op != "crypto" && op != "biz" {
		return common.ErrValueInvalid.Wrap(errors.Errorf("dao option must be crypto or biz, got %v", op))
	}

	return nil
}

func (op DAOOption) Bytes() []byte {
	return []byte(op)
}

var DesignHint = hint.MustNewHint("mitum-dao-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	option DAOOption
	policy Policy
}

func NewDesign(option DAOOption, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		option:     option,
		policy:     policy,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.option,
		de.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid Design: %v", err)
	}

	return nil
}

func (de Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		de.option.Bytes(),
		de.policy.Bytes(),
	)
}

func (de Design) Option() DAOOption {
	return de.option
}

func (de Design) Policy() Policy {
	return de.policy
}
