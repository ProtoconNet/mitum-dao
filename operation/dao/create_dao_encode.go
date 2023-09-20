package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *CreateDAOFact) unpack(enc encoder.Encoder,
	sa, ca, op, tk, th string,
	bf, bw []byte,
	prp, rp, prsp, vp, psp, edp uint64,
	to, qou uint,
	cid string,
) error {
	e := util.StringError("failed to unmarshal CreateDAOFact")

	fact.currency = currencytypes.CurrencyID(cid)
	fact.option = types.DAOOption(op)
	fact.votingPowerToken = currencytypes.CurrencyID(tk)
	fact.proposalReviewPeriod = prp
	fact.registrationPeriod = rp
	fact.preSnapshotPeriod = prsp
	fact.votingPeriod = vp
	fact.postSnapshotPeriod = psp
	fact.executionDelayPeriod = edp
	fact.turnout = types.PercentRatio(to)
	fact.quorum = types.PercentRatio(qou)

	if big, err := common.NewBigFromString(th); err != nil {
		return e.Wrap(err)
	} else {
		fact.threshold = big
	}

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.contract = a
	}

	if hinter, err := enc.Decode(bf); err != nil {
		return e.Wrap(err)
	} else if am, ok := hinter.(currencytypes.Amount); !ok {
		return e.Wrap(errors.Errorf("expected Amount, not %T", hinter))
	} else {
		fact.fee = am
	}

	if hinter, err := enc.Decode(bw); err != nil {
		return e.Wrap(err)
	} else if wl, ok := hinter.(types.Whitelist); !ok {
		return e.Wrap(errors.Errorf("expected Whitelist, not %T", hinter))
	} else {
		fact.whitelist = wl
	}

	return nil
}
