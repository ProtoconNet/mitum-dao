package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"time"
)

type TestExecuteProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Execute]
}

func NewTestExecuteProcessor(
	tp *test.TestProcessor,
) TestExecuteProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Execute](tp)
	return TestExecuteProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestExecuteProcessor) Create(bm []base.BlockMap) *TestExecuteProcessor {
	t.Opr, _ = NewExecuteProcessor(func() (base.BlockMap, bool, error) { return bm[len(bm)-1], true, nil })(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestExecuteProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestExecuteProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestExecuteProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestExecuteProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestExecuteProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestExecuteProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestExecuteProcessor) LoadOperation(fileName string,
) *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestExecuteProcessor) Print(fileName string,
) *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestExecuteProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, proposalID string, currency types.CurrencyID,
) *TestExecuteProcessor {
	op := NewExecute(
		NewExecuteFact(
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

func (t *TestExecuteProcessor) RunPreProcess() *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestExecuteProcessor) RunProcess() *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestExecuteProcessor) IsValid() *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestExecuteProcessor) Decode(fileName string) *TestExecuteProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
