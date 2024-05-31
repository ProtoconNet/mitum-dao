package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	daotypes "github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"time"
)

type TestCreateDAOProcessor struct {
	*test.BaseTestOperationProcessorNoItem[CreateDAO]
	option               daotypes.DAOOption
	votingPowerToken     types.CurrencyID
	threshold            common.Big
	fee                  types.Amount
	whitelist            daotypes.Whitelist
	proposalReviewPeriod uint64
	registrationPeriod   uint64
	preSnapshotPeriod    uint64
	votingPeriod         uint64
	postSnapshotPeriod   uint64
	executionDelayPeriod uint64
	turnout              daotypes.PercentRatio
	quorum               daotypes.PercentRatio
}

func NewTestCreateDAOProcessor(encs *encoder.Encoders) TestCreateDAOProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[CreateDAO](encs)
	return TestCreateDAOProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestCreateDAOProcessor) Create(bm []base.BlockMap) *TestCreateDAOProcessor {
	t.Opr, _ = NewCreateDAOProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestCreateDAOProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestCreateDAOProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestCreateDAOProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestCreateDAOProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestCreateDAOProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestCreateDAOProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestCreateDAOProcessor) SetWhitelist(whitelist []test.Account, active bool,
) *TestCreateDAOProcessor {
	var adrs []base.Address

	for i := range whitelist {
		adrs = append(adrs, whitelist[i].Address())
	}

	t.whitelist = daotypes.NewWhitelist(active, adrs)

	return t
}

func (t *TestCreateDAOProcessor) SetDAO(
	option string,
	votingPowerToken types.CurrencyID,
	threshold common.Big,
	fee types.Amount,
	proposalReviewPeriod,
	registrationPeriod,
	preSnapshotPeriod,
	votingPeriod,
	postSnapshotPeriod,
	executionDelayPeriod uint64,
	turnout,
	quorum uint8,
) *TestCreateDAOProcessor {
	t.option = daotypes.DAOOption(option)
	t.votingPowerToken = votingPowerToken
	t.threshold = threshold
	t.fee = fee
	t.proposalReviewPeriod = proposalReviewPeriod
	t.registrationPeriod = registrationPeriod
	t.preSnapshotPeriod = preSnapshotPeriod
	t.votingPeriod = votingPeriod
	t.postSnapshotPeriod = postSnapshotPeriod
	t.executionDelayPeriod = executionDelayPeriod
	t.turnout = daotypes.PercentRatio(turnout)
	t.quorum = daotypes.PercentRatio(quorum)

	return t
}

func (t *TestCreateDAOProcessor) LoadOperation(fileName string,
) *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestCreateDAOProcessor) Print(fileName string,
) *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestCreateDAOProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey,
	contract base.Address, proposalID string,
	proposal daotypes.Proposal, currency types.CurrencyID,
) *TestCreateDAOProcessor {
	op := NewCreateDAO(
		NewCreateDAOFact(
			[]byte("token"),
			sender,
			contract,
			t.option,
			t.votingPowerToken,
			t.threshold,
			t.fee,
			t.whitelist,
			t.proposalReviewPeriod,
			t.registrationPeriod,
			t.preSnapshotPeriod,
			t.votingPeriod,
			t.postSnapshotPeriod,
			t.executionDelayPeriod,
			t.turnout,
			t.quorum,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestCreateDAOProcessor) RunPreProcess() *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestCreateDAOProcessor) RunProcess() *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestCreateDAOProcessor) IsValid() *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestCreateDAOProcessor) Decode(fileName string) *TestCreateDAOProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
