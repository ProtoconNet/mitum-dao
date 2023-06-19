package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DAOOption string

func (op DAOOption) IsValid([]byte) error {
	if op != "crypto" && op != "biz" {
		return util.ErrInvalid.Errorf("invalid dao option; 'crypto' | 'biz'")
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
	daoID  types.ContractID
	policy Policy
}

func NewDesign(option DAOOption, daoID types.ContractID, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		option:     option,
		daoID:      daoID,
		policy:     policy,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
		de.option,
		de.daoID,
		de.policy,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid Design: %w", err)
	}

	return nil
}

func (de Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		de.option.Bytes(),
		de.daoID.Bytes(),
		de.policy.Bytes(),
	)
}

func (de Design) Option() DAOOption {
	return de.option
}

func (de Design) DAOID() types.ContractID {
	return de.daoID
}

func (de Design) Policy() Policy {
	return de.policy
}
