package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"time"
)

type TestRegisterProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Register]
}

func NewTestRegisterProcessor(
	tp *test.TestProcessor,
) TestRegisterProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Register](tp)
	return TestRegisterProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestRegisterProcessor) Create(bm []base.BlockMap) *TestRegisterProcessor {
	t.Opr, _ = NewRegisterProcessor(func() (base.BlockMap, bool, error) { return bm[len(bm)-1], true, nil })(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestRegisterProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestRegisterProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestRegisterProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestRegisterProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestRegisterProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestRegisterProcessor) LoadOperation(fileName string,
) *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestRegisterProcessor) Print(fileName string,
) *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestRegisterProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, proposalID string, delegated base.Address, currency types.CurrencyID,
) *TestRegisterProcessor {
	op := NewRegister(
		NewRegisterFact(
			[]byte("token"),
			sender,
			contract,
			proposalID,
			delegated,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestRegisterProcessor) RunPreProcess() *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestRegisterProcessor) RunProcess() *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestRegisterProcessor) IsValid() *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestRegisterProcessor) Decode(fileName string) *TestRegisterProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
