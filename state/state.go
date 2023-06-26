package state

import (
	"fmt"
	"strings"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	DAOPrefix            = "dao:"
	DesignStateValueHint = hint.MustNewHint("mitum-dao-design-state-value-v0.0.1")
	DesignSuffix         = ":design"
)

func StateKeyDAOPrefix(ca base.Address, daoID currencytypes.ContractID) string {
	return fmt.Sprintf("%s%s:%s", DAOPrefix, ca.String(), daoID)
}

type DesignStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		Design:     design,
	}
}

func (hd DesignStateValue) Hint() hint.Hint {
	return hd.BaseHinter.Hint()
}

func (hd DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DesignStateValue")

	if err := hd.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := hd.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (hd DesignStateValue) HashBytes() []byte {
	return hd.Design.Bytes()
}

func StateDesignValue(st base.State) (types.Design, error) {
	v := st.Value()
	if v == nil {
		return types.Design{}, util.ErrNotFound.Errorf("dao design not found in State")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return types.Design{}, errors.Errorf("invalid dao design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func StateKeyDesign(ca base.Address, daoid currencytypes.ContractID) string {
	return fmt.Sprintf("%s%s", StateKeyDAOPrefix(ca, daoid), DesignSuffix)
}

var (
	ProposalStateValueHint = hint.MustNewHint("mitum-dao-proposal-state-value-v0.0.1")
	ProposalSuffix         = ":dao-proposal"
)

type ProposalStateValue struct {
	hint.BaseHinter
	Proposal types.Proposal
}

func NewProposalStateValue(proposal types.Proposal) ProposalStateValue {
	return ProposalStateValue{
		BaseHinter: hint.NewBaseHinter(ProposalStateValueHint),
		Proposal:   proposal,
	}
}

func (p ProposalStateValue) Hint() hint.Hint {
	return p.BaseHinter.Hint()
}

func (p ProposalStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid ProposalStateValue")

	if err := p.BaseHinter.IsValid(ProposalStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := p.Proposal.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (p ProposalStateValue) HashBytes() []byte {
	return p.Proposal.Bytes()
}

func StateProposalValue(st base.State) (*types.Proposal, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("proposal not found in State")
	}

	d, ok := v.(ProposalStateValue)
	if !ok {
		return nil, errors.Errorf("invalid proposal value found, %T", v)
	}

	p := d.Proposal

	return &p, nil
}

func IsStateProposalKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, ProposalSuffix)
}

func StateKeyProposal(ca base.Address, daoid currencytypes.ContractID, pid string) string {
	return fmt.Sprintf("%s-%s%s", StateKeyDAOPrefix(ca, daoid), pid, ProposalSuffix)
}

var RegisterInfoHint = hint.MustNewHint("mitum-dao-register-info-v0.0.1")

type RegisterInfo struct {
	hint.BaseHinter
	Account    base.Address
	ApprovedBy []base.Address
}

func NewRegisterInfo(account base.Address, approvedBy []base.Address) RegisterInfo {
	return RegisterInfo{
		BaseHinter: hint.NewBaseHinter(RegisterInfoHint),
		Account:    account,
		ApprovedBy: approvedBy,
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

	if err := r.Account.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	for _, ac := range r.ApprovedBy {
		if err := ac.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if ac.Equal(r.Account) {
			return e.Wrap(errors.Errorf("approving address is same with approved address, %q", r.Account))
		}
	}

	return nil
}

func (r RegisterInfo) Bytes() []byte {
	ba := make([][]byte, len(r.ApprovedBy)+1)

	ba[0] = r.Account.Bytes()

	for i, ac := range r.ApprovedBy {
		ba[i+1] = ac.Bytes()
	}

	return util.ConcatBytesSlice(ba...)
}

var (
	RegisterListStateValueHint = hint.MustNewHint("mitum-dao-register-list-state-value-v0.0.1")
	RegisterListSuffix         = ":register-list"
)

type RegisterListStateValue struct {
	hint.BaseHinter
	Registers []RegisterInfo
}

func NewRegisterListStateValue(registers []RegisterInfo) RegisterListStateValue {
	return RegisterListStateValue{
		BaseHinter: hint.NewBaseHinter(RegisterListStateValueHint),
		Registers:  registers,
	}
}

func (r RegisterListStateValue) Hint() hint.Hint {
	return r.BaseHinter.Hint()
}

func (r RegisterListStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid RegisterListStateValue")

	if err := r.BaseHinter.IsValid(RegisterListStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, info := range r.Registers {
		if err := info.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if _, found := founds[info.Account.String()]; found {
			return e.Wrap(errors.Errorf("duplicate register account found, %q", info.Account))
		}
	}

	return nil
}

func (r RegisterListStateValue) HashBytes() []byte {
	bs := make([][]byte, len(r.Registers))

	for i, br := range r.Registers {
		bs[i] = br.Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func StateRegisterListValue(st base.State) ([]RegisterInfo, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("register list not found in State")
	}

	r, ok := v.(RegisterListStateValue)
	if !ok {
		return nil, errors.Errorf("invalid register list value found, %T", v)
	}

	return r.Registers, nil
}

func IsStateRegisterListKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, RegisterListSuffix)
}

func StateKeyRegisterList(ca base.Address, daoid currencytypes.ContractID, pid string) string {
	return fmt.Sprintf("%s-%s%s", StateKeyDAOPrefix(ca, daoid), pid, RegisterListSuffix)
}
