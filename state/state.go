package state

import (
	"fmt"
	"strings"

	"github.com/ProtoconNet/mitum-currency/v3/common"
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

var RegisterInfoHint = hint.MustNewHint("mitum-dao-register-info-v0.0.1")

type RegisterInfo struct {
	hint.BaseHinter
	account    base.Address
	approvedBy []base.Address
}

func NewRegisterInfo(account base.Address, approvedBy []base.Address) RegisterInfo {
	return RegisterInfo{
		BaseHinter: hint.NewBaseHinter(RegisterInfoHint),
		account:    account,
		approvedBy: approvedBy,
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

	if err := r.account.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	for _, ac := range r.approvedBy {
		if err := ac.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if ac.Equal(r.account) {
			return e.Wrap(errors.Errorf("approving address is same with approved address, %q", r.Account))
		}
	}

	return nil
}

func (r RegisterInfo) Bytes() []byte {
	ba := make([][]byte, len(r.approvedBy)+1)

	ba[0] = r.account.Bytes()

	for i, ac := range r.approvedBy {
		ba[i+1] = ac.Bytes()
	}

	return util.ConcatBytesSlice(ba...)
}

func (r RegisterInfo) Account() base.Address {
	return r.account
}

func (r RegisterInfo) ApprovedBy() []base.Address {
	return r.approvedBy
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

		if _, found := founds[info.account.String()]; found {
			return e.Wrap(errors.Errorf("duplicate register account found, %q", info.account))
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

var (
	VotingPowerHint = hint.MustNewHint("mitum-dao-voting-power-v0.0.1")
)

type VotingPower struct {
	hint.BaseHinter
	account     base.Address
	votingPower common.Big
}

func NewVotingPower(account base.Address, votingPower common.Big) VotingPower {
	return VotingPower{
		BaseHinter:  hint.NewBaseHinter(VotingPowerHint),
		account:     account,
		votingPower: votingPower,
	}
}

func (vp VotingPower) Hint() hint.Hint {
	return vp.BaseHinter.Hint()
}

func (vp VotingPower) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotingPower")

	if err := vp.BaseHinter.IsValid(SnapHistoryHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false, vp.account, vp.votingPower); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (vp VotingPower) Bytes() []byte {
	return util.ConcatBytesSlice(
		vp.account.Bytes(),
		vp.votingPower.Bytes(),
	)
}

func (vp VotingPower) Account() base.Address {
	return vp.account
}

func (vp VotingPower) VotingPower() common.Big {
	return vp.votingPower
}

var (
	SnapHistoryHint = hint.MustNewHint("mitum-dao-snap-history-v0.0.1")
)

type SnapHistory struct {
	hint.BaseHinter
	timestamp uint64
	snaps     []VotingPower
}

func NewSnapHistory(timestamp uint64, snaps []VotingPower) SnapHistory {
	return SnapHistory{
		BaseHinter: hint.NewBaseHinter(SnapHistoryHint),
		timestamp:  timestamp,
		snaps:      snaps,
	}
}

func (sh SnapHistory) Hint() hint.Hint {
	return sh.BaseHinter.Hint()
}

func (sh SnapHistory) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid SnapHistory")

	if err := sh.BaseHinter.IsValid(SnapHistoryHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, snap := range sh.snaps {
		if err := snap.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if _, found := founds[snap.account.String()]; found {
			return e.Wrap(errors.Errorf("duplicate snap account found, %q", snap.account))
		}

		founds[snap.account.String()] = struct{}{}
	}

	return nil
}

func (sh SnapHistory) Bytes() []byte {
	bs := make([][]byte, len(sh.snaps))

	for i, snap := range sh.snaps {
		bs[i] = snap.Bytes()
	}

	return util.ConcatBytesSlice(
		util.Uint64ToBytes(sh.timestamp),
		util.ConcatBytesSlice(bs...),
	)
}

func (sh SnapHistory) TimeStamp() uint64 {
	return sh.timestamp
}

func (sh SnapHistory) Snaps() []VotingPower {
	return sh.snaps
}

var (
	SnapHistoriesStateValueHint = hint.MustNewHint("mitum-dao-snap-histories-state-value-v0.0.1")
	SnapHistoriesSuffix         = ":snap-histories"
)

type SnapHistoriesStateValue struct {
	hint.BaseHinter
	Histories []SnapHistory
}

func NewSnapHistoriesStateValue(histories []SnapHistory) SnapHistoriesStateValue {
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

func StateSnapHistoriesValue(st base.State) ([]SnapHistory, error) {
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
	VotingPowersHint = hint.MustNewHint("mitum-dao-voting-powers-v0.0.1")
)

type VotingPowers struct {
	hint.BaseHinter
	total        common.Big
	votingPowers []VotingPower
}

func NewVotingPowers(total common.Big, votingPowers []VotingPower) VotingPowers {
	return VotingPowers{
		BaseHinter:   hint.NewBaseHinter(VotingPowersHint),
		total:        total,
		votingPowers: votingPowers,
	}
}

func (vp VotingPowers) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid VotingPowers")

	if err := vp.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	total := common.ZeroBig
	for _, vp := range vp.votingPowers {
		if err := vp.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		total = total.Add(vp.votingPower)
	}

	if total.Compare(vp.total) != 0 {
		return e.Wrap(errors.Errorf("invalid voting power total, %q != %q", total, vp.total))
	}

	return nil
}

func (vp VotingPowers) Bytes() []byte {
	bs := make([][]byte, len(vp.votingPowers))
	for i, v := range vp.votingPowers {
		bs[i] = v.Bytes()
	}

	return util.ConcatBytesSlice(
		vp.total.Bytes(),
		util.ConcatBytesSlice(bs...),
	)
}

func (vp VotingPowers) Total() common.Big {
	return vp.total
}

func (vp VotingPowers) VotingPowers() []VotingPower {
	return vp.votingPowers
}

var (
	VotesStateValueHint = hint.MustNewHint("mitum-dao-votes-state-value-v0.0.1")
	VotesSuffix         = ":votes"
)

type VotesStateValue struct {
	hint.BaseHinter
	Active bool
	Result uint8
	Votes  []VotingPowers
}

func NewVotesStateValue(active bool, result uint8, votes []VotingPowers) VotesStateValue {
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
