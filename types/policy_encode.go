package types

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (wl *Whitelist) unpack(enc encoder.Encoder, ht hint.Hint, at bool, bacs []byte) error {
	e := util.StringError("failed to decode bson of Whitelist")

	wl.active = at

	wl.BaseHinter = hint.NewBaseHinter(ht)

	hacs, err := enc.DecodeSlice(bacs)
	if err != nil {
		return e.Wrap(err)
	}

	accounts := make([]base.Address, len(hacs))
	for i := range hacs {
		j, ok := hacs[i].(base.Address)
		if !ok {
			return e.Wrap(errors.Errorf("expected base.Address, not %T", hacs[i]))
		}

		accounts[i] = j
	}
	wl.accounts = accounts

	return nil
}

func (po *Policy) unpack(enc encoder.Encoder, ht hint.Hint,
	cr string,
	bth, bf, bw []byte,
	rvp, rgp, prsp, vp, psp, edp uint64,
	to, qou uint,
) error {
	e := util.StringError("failed to decode bson of Policy")

	po.BaseHinter = hint.NewBaseHinter(ht)
	po.token = currencytypes.CurrencyID(cr)
	po.proposalReviewPeriod = rvp
	po.registrationPeriod = rgp
	po.preSnapshotPeriod = prsp
	po.votingPeriod = vp
	po.postSnapshotPeriod = psp
	po.executionDelayPeriod = edp
	po.turnout = PercentRatio(to)
	po.quorum = PercentRatio(qou)

	if hinter, err := enc.Decode(bth); err != nil {
		return e.Wrap(err)
	} else if am, ok := hinter.(currencytypes.Amount); !ok {
		return e.Wrap(errors.Errorf("expected Amount, not %T", hinter))
	} else {
		po.threshold = am
	}

	if hinter, err := enc.Decode(bf); err != nil {
		return e.Wrap(err)
	} else if am, ok := hinter.(currencytypes.Amount); !ok {
		return e.Wrap(errors.Errorf("expected Amount, not %T", hinter))
	} else {
		po.fee = am
	}

	if hinter, err := enc.Decode(bw); err != nil {
		return e.Wrap(err)
	} else if wl, ok := hinter.(Whitelist); !ok {
		return e.Wrap(errors.Errorf("expected Whitelist, not %T", hinter))
	} else {
		po.whitelist = wl
	}

	return nil
}
