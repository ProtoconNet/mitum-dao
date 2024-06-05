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
	existing map[string]types.VoterInfo
	add      map[string]types.VoterInfo
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

	s.existing = make(map[string]types.VoterInfo)
	s.add = make(map[string]types.VoterInfo)
	var voters []types.VoterInfo
	if nst.Value() != nil {
		voters = nst.Value().(VotersStateValue).voters
		for i := range voters {
			s.existing[voters[i].Account().String()] = voters[i]
		}
	}

	return s
}

func (s *VotersStateValueMerger) Merge(value base.StateValue, op util.Hash) error {
	s.Lock()
	defer s.Unlock()

	switch t := value.(type) {
	case VotersStateValue:
		for i := range t.voters {
			switch v, found := s.add[t.voters[i].Account().String()]; {
			case !found:
				s.add[t.voters[i].Account().String()] = t.voters[i]
			default:
				delegators := append(v.Delegators(), t.voters[i].Delegators()...)
				delegators, _ = util.RemoveDuplicatedSlice(delegators, func(address base.Address) (string, error) { return address.String(), nil })
				v.SetDelegators(delegators)
				s.add[t.voters[i].Account().String()] = v
			}
		}
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
		for k, v := range s.add {
			switch value, found := s.existing[k]; {
			case !found:
				s.existing[k] = v
			default:
				delegators := append(value.Delegators(), v.Delegators()...)
				delegators, _ = util.RemoveDuplicatedSlice(delegators, func(address base.Address) (string, error) { return address.String(), nil })
				value.SetDelegators(delegators)
				s.existing[k] = value
			}
		}
	}
	for _, v := range s.existing {
		nvoters = append(nvoters, v)
	}

	sort.Slice(nvoters, func(i, j int) bool { // NOTE sort by address
		return strings.Compare(nvoters[i].Account().String(), nvoters[j].Account().String()) < 0
	})

	return NewVotersStateValue(
		nvoters,
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
