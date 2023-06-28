package dao

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CreateDAOFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Owner            base.Address             `json:"sender"`
	Contract         base.Address             `json:"contract"`
	DAOID            currencytypes.ContractID `json:"daoid"`
	Option           types.DAOOption          `json:"option"`
	VotingPowerToken currencytypes.CurrencyID `json:"voting_power_token"`
	Threshold        currencytypes.Amount     `json:"threshold"`
	Fee              currencytypes.Amount     `json:"fee"`
	Whitelist        types.Whitelist          `json:"whitelist"`
	Delaytime        uint64                   `json:"delaytime"`
	Snaptime         uint64                   `json:"snaptime"`
	Voteperiod       uint64                   `json:"voteperiod"`
	Timelock         uint64                   `json:"timelock"`
	Turnout          types.PercentRatio       `json:"turnout"`
	Quorum           types.PercentRatio       `json:"quorum"`

	Currency currencytypes.CurrencyID `json:"currency"`
}

func (fact CreateDAOFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateDAOFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Owner:                 fact.sender,
		Contract:              fact.contract,
		DAOID:                 fact.daoID,
		Option:                fact.option,
		VotingPowerToken:      fact.votingPowerToken,
		Threshold:             fact.threshold,
		Fee:                   fact.fee,
		Whitelist:             fact.whitelist,
		Delaytime:             fact.delaytime,
		Snaptime:              fact.snaptime,
		Voteperiod:            fact.voteperiod,
		Timelock:              fact.timelock,
		Turnout:               fact.turnout,
		Quorum:                fact.quorum,
		Currency:              fact.currency,
	})
}

type CreateDAOFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Owner            string          `json:"sender"`
	Contract         string          `json:"contract"`
	DAOID            string          `json:"daoid"`
	Option           string          `json:"option"`
	VotingPowerToken string          `json:"voting_power_token"`
	Threshold        json.RawMessage `json:"threshold"`
	Fee              json.RawMessage `json:"fee"`
	Whitelist        json.RawMessage `json:"whitelist"`
	Delaytime        uint64          `json:"delaytime"`
	Snaptime         uint64          `json:"snaptime"`
	Voteperiod       uint64          `json:"voteperiod"`
	Timelock         uint64          `json:"timelock"`
	Turnout          uint            `json:"turnout"`
	Quorum           uint            `json:"quorum"`
	Currency         string          `json:"currency"`
}

func (fact *CreateDAOFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CreateDAOFact")

	var uf CreateDAOFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return e(err, "")
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	return fact.unpack(enc,
		uf.Owner,
		uf.Contract,
		uf.DAOID,
		uf.Option,
		uf.VotingPowerToken,
		uf.Threshold,
		uf.Fee,
		uf.Whitelist,
		uf.Delaytime,
		uf.Snaptime,
		uf.Voteperiod,
		uf.Timelock,
		uf.Turnout,
		uf.Quorum,
		uf.Currency,
	)
}

type CreateDAOMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateDAO) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateDAOMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateDAO) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of CreateDAO")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e(err, "")
	}

	op.BaseOperation = ubo

	return nil
}
