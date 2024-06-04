package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"time"
)

type TestPostSnapProcessor struct {
	*test.BaseTestOperationProcessorNoItem[PostSnap]
}

func NewTestPostSnapProcessor(
	tp *test.TestProcessor,
) TestPostSnapProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[PostSnap](tp)
	return TestPostSnapProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestPostSnapProcessor) Create(bm []base.BlockMap) *TestPostSnapProcessor {
	t.Opr, _ = NewPostSnapProcessor(func() (base.BlockMap, bool, error) { return bm[len(bm)-1], true, nil })(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestPostSnapProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestPostSnapProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestPostSnapProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestPostSnapProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestPostSnapProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestPostSnapProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestPostSnapProcessor) LoadOperation(fileName string,
) *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestPostSnapProcessor) Print(fileName string,
) *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestPostSnapProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, proposalID string, currency types.CurrencyID,
) *TestPostSnapProcessor {
	op := NewPostSnap(
		NewPostSnapFact(
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

func (t *TestPostSnapProcessor) RunPreProcess() *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestPostSnapProcessor) RunProcess() *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestPostSnapProcessor) IsValid() *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestPostSnapProcessor) Decode(fileName string) *TestPostSnapProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
