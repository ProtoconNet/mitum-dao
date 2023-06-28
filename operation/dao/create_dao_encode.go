package dao

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateDAOFact) unpack(enc encoder.Encoder,
	sa, ca, did, op, tk string,
	bth, bf, bw []byte,
	dt, rp, st, vp, tl uint64,
	to, qou uint,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal CreateDAOFact")

	fact.daoID = currencytypes.ContractID(did)
	fact.currency = currencytypes.CurrencyID(cid)
	fact.option = types.DAOOption(op)
	fact.votingPowerToken = currencytypes.CurrencyID(tk)
	fact.delaytime = dt
	fact.registerperiod = rp
	fact.snaptime = st
	fact.voteperiod = vp
	fact.timelock = tl
	fact.turnout = types.PercentRatio(to)
	fact.quorum = types.PercentRatio(qou)

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.contract = a
	}

	if hinter, err := enc.Decode(bth); err != nil {
		return e(err, "")
	} else if am, ok := hinter.(currencytypes.Amount); !ok {
		return e(util.ErrWrongType.Errorf("expected Amount, not %T", hinter), "")
	} else {
		fact.threshold = am
	}

	if hinter, err := enc.Decode(bf); err != nil {
		return e(err, "")
	} else if am, ok := hinter.(currencytypes.Amount); !ok {
		return e(util.ErrWrongType.Errorf("expected Amount, not %T", hinter), "")
	} else {
		fact.fee = am
	}

	if hinter, err := enc.Decode(bw); err != nil {
		return e(err, "")
	} else if wl, ok := hinter.(types.Whitelist); !ok {
		return e(util.ErrWrongType.Errorf("expected Whitelist, not %T", hinter), "")
	} else {
		fact.whitelist = wl
	}

	return nil
}
