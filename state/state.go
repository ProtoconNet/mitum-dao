package state

import (
	"fmt"
	"strings"

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

func StateKeyDAOPrefix(ca base.Address) string {
	return fmt.Sprintf("%s%s", DAOPrefix, ca.String())
}

type DesignStateValue struct {
	hint.BaseHinter
	design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		design:     design,
	}
}

func (de DesignStateValue) Hint() hint.Hint {
	return de.BaseHinter.Hint()
}

func (de DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DesignStateValue")

	if err := de.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := de.design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (de DesignStateValue) HashBytes() []byte {
	return de.design.Bytes()
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

	return d.design, nil
}

func IsStateDesignKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, DesignSuffix)
}

func StateKeyDesign(ca base.Address) string {
	return fmt.Sprintf("%s%s", StateKeyDAOPrefix(ca), DesignSuffix)
}

var (
	ProposalStateValueHint = hint.MustNewHint("mitum-dao-proposal-state-value-v0.0.1")
	ProposalSuffix         = ":dao-proposal"
)

type ProposalStateValue struct {
	hint.BaseHinter
	status   types.ProposalStatus
	proposal types.Proposal
	policy   types.Policy
}

func NewProposalStateValue(status types.ProposalStatus, proposal types.Proposal, policy types.Policy) ProposalStateValue {
	return ProposalStateValue{
		BaseHinter: hint.NewBaseHinter(ProposalStateValueHint),
		status:     status,
		proposal:   proposal,
		policy:     policy,
	}
}

func (p ProposalStateValue) Hint() hint.Hint {
	return p.BaseHinter.Hint()
}

func (p ProposalStateValue) Status() types.ProposalStatus {
	return p.status
}

func (p ProposalStateValue) Proposal() types.Proposal {
	return p.proposal
}

func (p ProposalStateValue) Policy() types.Policy {
	return p.policy
}

func (p ProposalStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid ProposalStateValue")

	if err := p.BaseHinter.IsValid(ProposalStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(
		nil, false,
		p.proposal,
		p.policy,
	); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (p ProposalStateValue) HashBytes() []byte {
	return p.proposal.Bytes()
}

func StateProposalValue(st base.State) (ProposalStateValue, error) {
	v := st.Value()
	if v == nil {
		return ProposalStateValue{}, util.ErrNotFound.Errorf("proposal not found in State")
	}

	p, ok := v.(ProposalStateValue)
	if !ok {
		return ProposalStateValue{}, errors.Errorf("invalid proposal value found, %T", v)
	}

	return p, nil
}

func IsStateProposalKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, ProposalSuffix)
}

func StateKeyProposal(ca base.Address, pid string) string {
	return fmt.Sprintf("%s:%s%s", StateKeyDAOPrefix(ca), pid, ProposalSuffix)
}

var (
	DelegatorsStateValueHint = hint.MustNewHint("mitum-dao-delegators-state-value-v0.0.1")
	DelegatorsSuffix         = ":delegators"
)

type DelegatorsStateValue struct {
	hint.BaseHinter
	delegators []types.DelegatorInfo
}

func NewDelegatorsStateValue(delegators []types.DelegatorInfo) DelegatorsStateValue {
	return DelegatorsStateValue{
		BaseHinter: hint.NewBaseHinter(DelegatorsStateValueHint),
		delegators: delegators,
	}
}

func (dg DelegatorsStateValue) Hint() hint.Hint {
	return dg.BaseHinter.Hint()
}

func (dg DelegatorsStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DelegatorsStateValue")

	if err := dg.BaseHinter.IsValid(DelegatorsStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	for _, ac := range dg.delegators {
		if err := ac.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}

func (dg DelegatorsStateValue) HashBytes() []byte {
	ba := make([][]byte, len(dg.delegators))

	for i, delegator := range dg.delegators {
		ba[i] = delegator.Bytes()
	}

	return util.ConcatBytesSlice(ba...)
}

func StateDelegatorsValue(st base.State) ([]types.DelegatorInfo, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("delegators not found in State")
	}

	ap, ok := v.(DelegatorsStateValue)
	if !ok {
		return nil, errors.Errorf("invalid delegators value found, %T", v)
	}

	return ap.delegators, nil
}

func IsStateDelegatorsKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, DelegatorsSuffix)
}

func StateKeyDelegators(ca base.Address, pid string) string {
	return fmt.Sprintf("%s:%s%s", StateKeyDAOPrefix(ca), pid, DelegatorsSuffix)
}

var (
	VotersStateValueHint = hint.MustNewHint("mitum-dao-voters-state-value-v0.0.1")
	VotersSuffix         = ":voters"
)

type VotersStateValue struct {
	hint.BaseHinter
	voters []types.VoterInfo
}

func NewVotersStateValue(voters []types.VoterInfo) VotersStateValue {
	return VotersStateValue{
		BaseHinter: hint.NewBaseHinter(VotersStateValueHint),
		voters:     voters,
	}
}

func (vt VotersStateValue) Hint() hint.Hint {
	return vt.BaseHinter.Hint()
}

func (vt VotersStateValue) Voters() []types.VoterInfo {
	return vt.voters
}

func (vt VotersStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotersStateValue")

	if err := vt.BaseHinter.IsValid(VotersStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, info := range vt.voters {
		if err := info.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if _, found := founds[info.Account().String()]; found {
			return e.Wrap(errors.Errorf("duplicate voter account found, %q", info.Account()))
		}
	}

	return nil
}

func (vt VotersStateValue) HashBytes() []byte {
	bs := make([][]byte, len(vt.voters))

	for i, br := range vt.voters {
		bs[i] = br.Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func StateVotersValue(st base.State) ([]types.VoterInfo, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("voters not found in State")
	}

	r, ok := v.(VotersStateValue)
	if !ok {
		return nil, errors.Errorf("invalid voters value found, %T", v)
	}

	return r.voters, nil
}

func IsStateVotersKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, VotersSuffix)
}

func StateKeyVoters(ca base.Address, pid string) string {
	return fmt.Sprintf("%s:%s%s", StateKeyDAOPrefix(ca), pid, VotersSuffix)
}

//var (
//	ProposalStatusStateValueHint = hint.MustNewHint("mitum-dao-proposal-status-state-value-v0.0.1")
//	ProposalStatusSuffix         = ":proposalstatus"
//)

//type ProposalStatusStateValue struct {
//	hint.BaseHinter
//	proposalStatus types.ProposalStatus
//}
//
//func NewProposalStatusStateValue(proposalStatus types.ProposalStatus) ProposalStatusStateValue {
//	return ProposalStatusStateValue{
//		BaseHinter:     hint.NewBaseHinter(ProposalStatusStateValueHint),
//		proposalStatus: proposalStatus,
//	}
//}
//
//func (ps ProposalStatusStateValue) Hint() hint.Hint {
//	return ps.BaseHinter.Hint()
//}
//
//func (ps ProposalStatusStateValue) ProposalStatus() types.ProposalStatus {
//	return ps.proposalStatus
//}
//
//func (ps ProposalStatusStateValue) IsValid([]byte) error {
//	e := util.ErrInvalid.Errorf("invalid ProposalStatusStateValue")
//
//	if err := ps.BaseHinter.IsValid(ProposalStatusStateValueHint.Type().Bytes()); err != nil {
//		return e.Wrap(err)
//	}
//
//	return nil
//}
//
//func (ps ProposalStatusStateValue) HashBytes() []byte {
//	return ps.proposalStatus.Bytes()
//}
//
//func StateProposalStatusValue(st base.State) (types.ProposalStatus, error) {
//	v := st.Value()
//	if v == nil {
//		return types.NilStatus, util.ErrNotFound.Errorf("ProposalStatus value not found in State")
//	}
//
//	ps, ok := v.(ProposalStatusStateValue)
//	if !ok {
//		return types.NilStatus, errors.Errorf("invalid ProposalStatus value found, %T", v)
//	}
//
//	return ps.proposalStatus, nil
//}
//
//func IsStateProposalStatusKey(key string) bool {
//	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, ProposalStatusSuffix)
//}
//
//func StateKeyProposalStatus(ca base.Address, daoID currencytypes.ContractID, pid string) string {
//	return fmt.Sprintf("%s-%s%s", StateKeyDAOPrefix(ca, daoID), pid, ProposalStatusSuffix)
//}

var (
	VotingPowerBoxStateValueHint = hint.MustNewHint("mitum-dao-voting-power-box-state-value-v0.0.1")
	VotingPowerBoxSuffix         = ":votingpowerbox"
)

type VotingPowerBoxStateValue struct {
	hint.BaseHinter
	votingPowerBox types.VotingPowerBox
}

func NewVotingPowerBoxStateValue(votingPowerBox types.VotingPowerBox) VotingPowerBoxStateValue {
	return VotingPowerBoxStateValue{
		BaseHinter:     hint.NewBaseHinter(VotingPowerBoxStateValueHint),
		votingPowerBox: votingPowerBox,
	}
}

func (vb VotingPowerBoxStateValue) Hint() hint.Hint {
	return vb.BaseHinter.Hint()
}

func (vb VotingPowerBoxStateValue) VotingPowerBox() types.VotingPowerBox {
	return vb.votingPowerBox
}

func (vb VotingPowerBoxStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotingPowerBoxStateValue")

	if err := vb.BaseHinter.IsValid(VotingPowerBoxStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (vb VotingPowerBoxStateValue) HashBytes() []byte {
	return vb.votingPowerBox.Bytes()
}

func StateVotingPowerBoxValue(st base.State) (types.VotingPowerBox, error) {
	v := st.Value()
	if v == nil {
		return types.VotingPowerBox{}, util.ErrNotFound.Errorf("VotingPowerBox not found in State")
	}

	r, ok := v.(VotingPowerBoxStateValue)
	if !ok {
		return types.VotingPowerBox{}, errors.Errorf("invalid VotingPowerBox value found, %T", v)
	}

	return r.votingPowerBox, nil
}

func IsStateVotingPowerBoxKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, VotingPowerBoxSuffix)
}

func StateKeyVotingPowerBox(ca base.Address, pid string) string {
	return fmt.Sprintf("%s:%s%s", StateKeyDAOPrefix(ca), pid, VotingPowerBoxSuffix)
}

func ParseStateKey(key string, Prefix string) ([]string, error) {
	parsedKey := strings.Split(key, ":")
	if parsedKey[0] != Prefix[:len(Prefix)-1] {
		return nil, errors.Errorf("State Key not include Prefix, %s", parsedKey)
	}
	if len(parsedKey) < 3 {
		return nil, errors.Errorf("parsing State Key string failed, %s", parsedKey)
	} else {
		return parsedKey, nil
	}
}
