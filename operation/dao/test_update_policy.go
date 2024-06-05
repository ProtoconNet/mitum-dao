package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	daotypes "github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"time"
)

type TestUpdatePolicyProcessor struct {
	*test.BaseTestOperationProcessorNoItem[UpdatePolicy]
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

func NewTestUpdatePolicyProcessor(
	tp *test.TestProcessor,
) TestUpdatePolicyProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[UpdatePolicy](tp)
	return TestUpdatePolicyProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestUpdatePolicyProcessor) Create() *TestUpdatePolicyProcessor {
	t.Opr, _ = NewUpdatePolicyProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestUpdatePolicyProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []types.CurrencyID, instate bool,
) *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestUpdatePolicyProcessor) SetAmount(
	am int64, cid types.CurrencyID, target []types.Amount,
) *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestUpdatePolicyProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestUpdatePolicyProcessor) SetAccount(
	priv string, amount int64, cid types.CurrencyID, target []test.Account, inState bool,
) *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestUpdatePolicyProcessor) SetProposal(
	proposer base.Address, startTime uint64, url, hash string, options uint8, target []daotypes.Proposal,
) *TestUpdatePolicyProcessor {
	pr := daotypes.NewBizProposal(proposer, startTime, daotypes.URL(url), hash, options)
	test.UpdateSlice[daotypes.Proposal](pr, target)

	return t
}

func (t *TestUpdatePolicyProcessor) SetBlockMap(
	proposedAt int64, target []base.BlockMap,
) *TestUpdatePolicyProcessor {
	bm := BlockMap{
		manifest: Manifest{proposedAt: time.Unix(proposedAt, 0)},
	}
	test.UpdateSlice[base.BlockMap](bm, target)

	return t
}

func (t *TestUpdatePolicyProcessor) SetWhitelist(whitelist []test.Account, active bool,
) *TestUpdatePolicyProcessor {
	var adrs []base.Address

	for i := range whitelist {
		adrs = append(adrs, whitelist[i].Address())
	}

	t.whitelist = daotypes.NewWhitelist(active, adrs)

	return t
}

func (t *TestUpdatePolicyProcessor) SetDAO(
	option daotypes.DAOOption,
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
	quorum daotypes.PercentRatio,
) *TestUpdatePolicyProcessor {
	t.option = option
	t.votingPowerToken = votingPowerToken
	t.threshold = threshold
	t.fee = fee
	t.proposalReviewPeriod = proposalReviewPeriod
	t.registrationPeriod = registrationPeriod
	t.preSnapshotPeriod = preSnapshotPeriod
	t.votingPeriod = votingPeriod
	t.postSnapshotPeriod = postSnapshotPeriod
	t.executionDelayPeriod = executionDelayPeriod
	t.turnout = turnout
	t.quorum = quorum

	return t
}

func (t *TestUpdatePolicyProcessor) LoadOperation(fileName string,
) *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestUpdatePolicyProcessor) Print(fileName string,
) *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestUpdatePolicyProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey,
	contract base.Address, currency types.CurrencyID,
) *TestUpdatePolicyProcessor {
	op := NewUpdatePolicy(
		NewUpdatePolicyFact(
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

func (t *TestUpdatePolicyProcessor) RunPreProcess() *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestUpdatePolicyProcessor) RunProcess() *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestUpdatePolicyProcessor) IsValid() *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.IsValid()

	return t
}

func (t *TestUpdatePolicyProcessor) Decode(fileName string) *TestUpdatePolicyProcessor {
	t.BaseTestOperationProcessorNoItem.Decode(fileName)

	return t
}
