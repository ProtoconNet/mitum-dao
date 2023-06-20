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
	e := util.ErrInvalid.Errorf("invalid p ProposalStateValue")

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
