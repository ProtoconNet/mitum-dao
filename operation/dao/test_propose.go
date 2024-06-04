package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	daotypes "github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"time"
)

type TestProposeProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Propose]
}

func NewTestProposeProcessor(
	tp *test.TestProcessor,
) TestProposeProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Propose](tp)
	return TestProposeProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestProposeProcessor) Create(bm []base.BlockMap) *TestProposeProcessor {
	t.Opr, _ = NewProposeProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestProposeProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestProposeProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestProposeProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestProposeProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestProposeProcessor) SetProposal(
	proposer base.Address, startTime uint64, url, hash string, options uint8, target []daotypes.Proposal,
) *TestProposeProcessor {
	pr := daotypes.NewBizProposal(proposer, startTime, daotypes.URL(url), hash, options)
	test.UpdateSlice[daotypes.Proposal](pr, target)

	return t
}

func (t *TestProposeProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestProposeProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestProposeProcessor) LoadOperation(fileName string,
) *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestProposeProcessor) Print(fileName string,
) *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestProposeProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, proposalID string, proposal daotypes.Proposal, currency types.CurrencyID,
) *TestProposeProcessor {
	op := NewPropose(
		NewProposeFact(
			[]byte("token"),
			sender,
			contract,
			proposalID,
			proposal,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestProposeProcessor) RunPreProcess() *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestProposeProcessor) RunProcess() *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestProposeProcessor) IsValid() *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestProposeProcessor) Decode(fileName string) *TestProposeProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
