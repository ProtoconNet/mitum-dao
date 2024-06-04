package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"time"
)

type TestVoteProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Vote]
}

func NewTestVoteProcessor(
	tp *test.TestProcessor,
) TestVoteProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Vote](tp)
	return TestVoteProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestVoteProcessor) Create(bm []base.BlockMap) *TestVoteProcessor {
	t.Opr, _ = NewVoteProcessor(func() (base.BlockMap, bool, error) { return bm[len(bm)-1], true, nil })(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestVoteProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestVoteProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestVoteProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestVoteProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestVoteProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestVoteProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestVoteProcessor) LoadOperation(fileName string,
) *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestVoteProcessor) Print(fileName string,
) *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestVoteProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, proposalID string, vote uint8, currency types.CurrencyID,
) *TestVoteProcessor {
	op := NewVote(
		NewVoteFact(
			[]byte("token"),
			sender,
			contract,
			proposalID,
			vote,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestVoteProcessor) RunPreProcess() *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestVoteProcessor) RunProcess() *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestVoteProcessor) IsValid() *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestVoteProcessor) Decode(fileName string) *TestVoteProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
