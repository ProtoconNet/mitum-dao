package state

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"sort"
	"strings"
	"sync"
)

type VotersStateValueMerger struct {
	*common.BaseStateValueMerger
	existing []types.VoterInfo
	add      []types.VoterInfo
	sync.Mutex
}

func NewVotersStateValueMerger(height base.Height, key string, st base.State) *VotersStateValueMerger {
	nst := st
	if st == nil {
		nst = common.NewBaseState(base.NilHeight, key, nil, nil, nil)
	}

	s := &VotersStateValueMerger{
		BaseStateValueMerger: common.NewBaseStateValueMerger(height, nst.Key(), nst),
	}

	if nst.Value() != nil {
		s.existing = nst.Value().(VotersStateValue).voters //nolint:forcetypeassert //...
	}

	return s
}

func (s *VotersStateValueMerger) Merge(value base.StateValue, op util.Hash) error {
	s.Lock()
	defer s.Unlock()

	switch t := value.(type) {
	case VotersStateValue:
		s.add = append(s.add, t.voters...)
	default:
		return errors.Errorf("unsupported voters state value, %T", value)
	}

	s.AddOperation(op)

	return nil
}

func (s *VotersStateValueMerger) CloseValue() (base.State, error) {
	s.Lock()
	defer s.Unlock()

	newValue, err := s.closeValue()
	if err != nil {
		return nil, errors.WithMessage(err, "close VotersStateValueMerger")
	}

	s.BaseStateValueMerger.SetValue(newValue)

	return s.BaseStateValueMerger.CloseValue()
}

func (s *VotersStateValueMerger) closeValue() (base.StateValue, error) {
	var nvoters []types.VoterInfo
	if len(s.add) > 0 {
		nvoters = append(s.existing, s.add...)
	} else {
		nvoters = s.existing
	}

	rvoters, _ := util.RemoveDuplicatedSlice(nvoters, func(v types.VoterInfo) (string, error) { return string(v.Bytes()), nil })
	sort.Slice(rvoters, func(i, j int) bool { // NOTE sort by address
		return strings.Compare(string(rvoters[i].Bytes()), string(rvoters[j].Bytes())) < 0
	})

	return NewVotersStateValue(
		rvoters,
	), nil
}

type DelegatorsStateValueMerger struct {
	*common.BaseStateValueMerger
	existing []types.DelegatorInfo
	add      []types.DelegatorInfo
	sync.Mutex
}

func NewDelegatorsStateValueMerger(height base.Height, key string, st base.State) *DelegatorsStateValueMerger {
	nst := st
	if st == nil {
		nst = common.NewBaseState(base.NilHeight, key, nil, nil, nil)
	}

	s := &DelegatorsStateValueMerger{
		BaseStateValueMerger: common.NewBaseStateValueMerger(height, nst.Key(), nst),
	}

	if nst.Value() != nil {
		s.existing = nst.Value().(DelegatorsStateValue).delegators //nolint:forcetypeassert //...
	}

	return s
}

func (s *DelegatorsStateValueMerger) Merge(value base.StateValue, op util.Hash) error {
	s.Lock()
	defer s.Unlock()

	switch t := value.(type) {
	case DelegatorsStateValue:
		s.add = append(s.add, t.delegators...)
	default:
		return errors.Errorf("unsupported delegators state value, %T", value)
	}

	s.AddOperation(op)

	return nil
}

func (s *DelegatorsStateValueMerger) CloseValue() (base.State, error) {
	s.Lock()
	defer s.Unlock()

	newValue, err := s.closeValue()
	if err != nil {
		return nil, errors.WithMessage(err, "close DelegatorsStateValueMerger")
	}

	s.BaseStateValueMerger.SetValue(newValue)

	return s.BaseStateValueMerger.CloseValue()
}

func (s *DelegatorsStateValueMerger) closeValue() (base.StateValue, error) {
	var ndelegators []types.DelegatorInfo
	if len(s.add) > 0 {
		ndelegators = append(s.existing, s.add...)
	} else {
		ndelegators = s.existing
	}

	rdelegators, _ := util.RemoveDuplicatedSlice(ndelegators, func(v types.DelegatorInfo) (string, error) { return string(v.Bytes()), nil })
	sort.Slice(rdelegators, func(i, j int) bool { // NOTE sort by address
		return strings.Compare(string(rdelegators[i].Bytes()), string(rdelegators[j].Bytes())) < 0
	})

	return NewDelegatorsStateValue(
		rdelegators,
	), nil
}
