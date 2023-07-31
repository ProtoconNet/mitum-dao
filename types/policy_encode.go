package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (wl *Whitelist) unpack(enc encoder.Encoder, ht hint.Hint, at bool, acs []string) error {
	e := util.StringError("failed to unmarshal Whitelist")

	wl.active = at

	wl.BaseHinter = hint.NewBaseHinter(ht)

	accs := make([]base.Address, len(acs))
	for i, ac := range acs {
		switch a, err := base.DecodeAddress(ac, enc); {
		case err != nil:
			return e.Wrap(err)
		default:
			accs[i] = a
		}
	}
	wl.accounts = accs

	return nil
}

func (po *Policy) unpack(enc encoder.Encoder, ht hint.Hint,
	cr, th string,
	bf, bw []byte,
	rvp, rgp, prsp, vp, psp, edp uint64,
	to, qou uint,
) error {
	e := util.StringError("failed to unmarshal Policy")

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

	if big, err := common.NewBigFromString(th); err != nil {
		return e.Wrap(err)
	} else {
		po.threshold = big
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
