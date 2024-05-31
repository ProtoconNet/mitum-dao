package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"time"
)

type BlockMap struct {
	base.BlockMap
	manifest Manifest
}

func (b BlockMap) Manifest() base.Manifest {
	return b.manifest
}

type Manifest struct {
	base.Manifest
	proposedAt time.Time
}

func (m Manifest) ProposedAt() time.Time {
	return m.proposedAt
}

type TestCancelProposalProcessor struct {
	*test.BaseTestOperationProcessorNoItem[CancelProposal]
}

func NewTestCancelProposalProcessor(encs *encoder.Encoders) TestCancelProposalProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[CancelProposal](encs)
	return TestCancelProposalProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestCancelProposalProcessor) Create(bm []base.BlockMap) *TestCancelProposalProcessor {
	t.Opr, _ = NewCancelProposalProcessor(func() (base.BlockMap, bool, error) { return bm[len(bm)-1], true, nil })(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestCancelProposalProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestCancelProposalProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestCancelProposalProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestCancelProposalProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestCancelProposalProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestCancelProposalProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestCancelProposalProcessor) LoadOperation(fileName string,
) *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestCancelProposalProcessor) Print(fileName string,
) *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestCancelProposalProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, proposalID string, currency types.CurrencyID,
) *TestCancelProposalProcessor {
	op := NewCancelProposal(
		NewCancelProposalFact(
			[]byte("token"),
			sender,
			contract,
			proposalID,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestCancelProposalProcessor) RunPreProcess() *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestCancelProposalProcessor) RunProcess() *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestCancelProposalProcessor) IsValid() *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestCancelProposalProcessor) Decode(fileName string) *TestCancelProposalProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
