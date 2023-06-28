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
	Active   bool
	Proposal types.Proposal
}

func NewProposalStateValue(active bool, proposal types.Proposal) ProposalStateValue {
	return ProposalStateValue{
		BaseHinter: hint.NewBaseHinter(ProposalStateValueHint),
		Active:     active,
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
	b := make([]byte, 1)
	if p.Active {
		b[0] = 1
	} else {
		b[0] = 0
	}

	return util.ConcatBytesSlice(
		b,
		p.Proposal.Bytes(),
	)
}

func StateProposalValue(st base.State) (ProposalStateValue, error) {
	v := st.Value()
	if v == nil {
		return ProposalStateValue{}, util.ErrNotFound.Errorf("proposal not found in State")
	}

	d, ok := v.(ProposalStateValue)
	if !ok {
		return ProposalStateValue{}, errors.Errorf("invalid proposal value found, %T", v)
	}

	return d, nil
}

func IsStateProposalKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, ProposalSuffix)
}

func StateKeyProposal(ca base.Address, daoid currencytypes.ContractID, pid string) string {
	return fmt.Sprintf("%s-%s%s", StateKeyDAOPrefix(ca, daoid), pid, ProposalSuffix)
}

var (
	ApprovingListStateValueHint = hint.MustNewHint("mitum-dao-approving-list-state-value-v0.0.1")
	ApprovingListSuffix         = ":approving-list"
)

type ApprovingListStateValue struct {
	hint.BaseHinter
	Accounts []base.Address
}

func NewApprovingListStateValue(accounts []base.Address) ApprovingListStateValue {
	return ApprovingListStateValue{
		BaseHinter: hint.NewBaseHinter(ApprovingListStateValueHint),
		Accounts:   accounts,
	}
}

func (ap ApprovingListStateValue) Hint() hint.Hint {
	return ap.BaseHinter.Hint()
}

func (ap ApprovingListStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid ApprovingListStateValue")

	if err := ap.BaseHinter.IsValid(ApprovingListStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	for _, ac := range ap.Accounts {
		if err := ac.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}

func (ap ApprovingListStateValue) HashBytes() []byte {
	ba := make([][]byte, len(ap.Accounts))

	for i, ac := range ap.Accounts {
		ba[i] = ac.Bytes()
	}

	return util.ConcatBytesSlice(ba...)
}

func StateApprovingListValue(st base.State) ([]base.Address, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("approving list not found in State")
	}

	ap, ok := v.(ApprovingListStateValue)
	if !ok {
		return nil, errors.Errorf("invalid approving list value found, %T", v)
	}

	return ap.Accounts, nil
}

func IsStateApprovingListKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, ApprovingListSuffix)
}

func StateKeyApprovingList(ca base.Address, daoid currencytypes.ContractID, pid string, ac base.Address) string {
	return fmt.Sprintf("%s-%s-%s%s", StateKeyDAOPrefix(ca, daoid), pid, ac.String(), ApprovingListSuffix)
}

var (
	RegisterListStateValueHint = hint.MustNewHint("mitum-dao-register-list-state-value-v0.0.1")
	RegisterListSuffix         = ":register-list"
)

type RegisterListStateValue struct {
	hint.BaseHinter
	Registers []types.RegisterInfo
}

func NewRegisterListStateValue(registers []types.RegisterInfo) RegisterListStateValue {
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

		if _, found := founds[info.Account().String()]; found {
			return e.Wrap(errors.Errorf("duplicate register account found, %q", info.Account()))
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

func StateRegisterListValue(st base.State) ([]types.RegisterInfo, error) {
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

var (
	SnapHistoriesStateValueHint = hint.MustNewHint("mitum-dao-snap-histories-state-value-v0.0.1")
	SnapHistoriesSuffix         = ":snap-histories"
)

type SnapHistoriesStateValue struct {
	hint.BaseHinter
	Histories []types.SnapHistory
}

func NewSnapHistoriesStateValue(histories []types.SnapHistory) SnapHistoriesStateValue {
	return SnapHistoriesStateValue{
		BaseHinter: hint.NewBaseHinter(SnapHistoriesStateValueHint),
		Histories:  histories,
	}
}

func (sh SnapHistoriesStateValue) Hint() hint.Hint {
	return sh.BaseHinter.Hint()
}

func (sh SnapHistoriesStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid SnapHistoriesStateValue")

	if err := sh.BaseHinter.IsValid(SnapHistoriesStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	for _, h := range sh.Histories {
		if err := h.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}

func (sh SnapHistoriesStateValue) HashBytes() []byte {
	bs := make([][]byte, len(sh.Histories))

	for i, h := range sh.Histories {
		bs[i] = h.Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func StateSnapHistoriesValue(st base.State) ([]types.SnapHistory, error) {
	v := st.Value()
	if v == nil {
		return nil, util.ErrNotFound.Errorf("snap histories not found in State")
	}

	hs, ok := v.(SnapHistoriesStateValue)
	if !ok {
		return nil, errors.Errorf("invalid snap histories value found, %T", v)
	}

	return hs.Histories, nil
}

func IsStateSnapHistoriesKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, SnapHistoriesSuffix)
}

func StateKeySnapHistories(ca base.Address, daoid currencytypes.ContractID, pid string) string {
	return fmt.Sprintf("%s-%s%s", StateKeyDAOPrefix(ca, daoid), pid, SnapHistoriesSuffix)
}

var (
	VotesStateValueHint = hint.MustNewHint("mitum-dao-votes-state-value-v0.0.1")
	VotesSuffix         = ":votes"
)

type VotesStateValue struct {
	hint.BaseHinter
	Active bool
	Result uint8
	Votes  []types.VotingPowers
}

func NewVotesStateValue(active bool, result uint8, votes []types.VotingPowers) VotesStateValue {
	return VotesStateValue{
		BaseHinter: hint.NewBaseHinter(VotesStateValueHint),
		Active:     active,
		Result:     result,
		Votes:      votes,
	}
}

func (v VotesStateValue) Hint() hint.Hint {
	return v.BaseHinter.Hint()
}

func (v VotesStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotesStateValue")

	if err := v.BaseHinter.IsValid(VotesStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (v VotesStateValue) HashBytes() []byte {
	b := make([]byte, 1)
	if v.Active {
		b[0] = 1
	} else {
		b[0] = 0
	}

	rs := make([][]byte, len(v.Votes))

	for i, t := range v.Votes {
		rs[i+1] = t.Bytes()
	}

	return util.ConcatBytesSlice(
		b,
		util.Uint8ToBytes(v.Result),
		util.ConcatBytesSlice(rs...),
	)
}

func StateVotesValue(st base.State) (VotesStateValue, error) {
	v := st.Value()
	if v == nil {
		return VotesStateValue{}, util.ErrNotFound.Errorf("voting votes not found in State")
	}

	r, ok := v.(VotesStateValue)
	if !ok {
		return VotesStateValue{}, errors.Errorf("invalid voting votes value found, %T", v)
	}

	return r, nil
}

func IsVotesKey(key string) bool {
	return strings.HasPrefix(key, DAOPrefix) && strings.HasSuffix(key, VotesSuffix)
}

func StateKeyVotes(ca base.Address, daoid currencytypes.ContractID, pid string) string {
	return fmt.Sprintf("%s-%s%s", StateKeyDAOPrefix(ca, daoid), pid, VotesSuffix)
}
