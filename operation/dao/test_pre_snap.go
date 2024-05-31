package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"time"
)

type TestPreSnapProcessor struct {
	*test.BaseTestOperationProcessorNoItem[PreSnap]
}

func NewTestPreSnapProcessor(encs *encoder.Encoders) TestPreSnapProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[PreSnap](encs)
	return TestPreSnapProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestPreSnapProcessor) Create(bm []base.BlockMap) *TestPreSnapProcessor {
	t.Opr, _ = NewPreSnapProcessor(func() (base.BlockMap, bool, error) { return bm[len(bm)-1], true, nil })(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestPreSnapProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestPreSnapProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestPreSnapProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestPreSnapProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestPreSnapProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestPreSnapProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestPreSnapProcessor) LoadOperation(fileName string,
) *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestPreSnapProcessor) Print(fileName string,
) *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestPreSnapProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, proposalID string, currency types.CurrencyID,
) *TestPreSnapProcessor {
	op := NewPreSnap(
		NewPreSnapFact(
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

func (t *TestPreSnapProcessor) RunPreProcess() *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestPreSnapProcessor) RunProcess() *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestPreSnapProcessor) IsValid() *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestPreSnapProcessor) Decode(fileName string) *TestPreSnapProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
