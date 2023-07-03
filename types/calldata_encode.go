package types

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (cd *TransferCallData) unpack(enc encoder.Encoder, ht hint.Hint, sd, rc string, bam []byte) error {
	e := util.StringErrorFunc("failed to decode bson of TransferCallData")

	cd.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(sd, enc); {
	case err != nil:
		return e(err, "")
	default:
		cd.sender = a
	}

	switch a, err := base.DecodeAddress(rc, enc); {
	case err != nil:
		return e(err, "")
	default:
		cd.receiver = a
	}

	if hinter, err := enc.Decode(bam); err != nil {
		return e(err, "")
	} else if am, ok := hinter.(currencytypes.Amount); !ok {
		return e(util.ErrWrongType.Errorf("expected Amount, not %T", hinter), "")
	} else {
		cd.amount = am
	}

	return nil
}

func (cd *GovernanceCallData) unpack(enc encoder.Encoder, ht hint.Hint, bpo []byte) error {
	e := util.StringErrorFunc("failed to decode bson of GovernanceCallData")

	cd.BaseHinter = hint.NewBaseHinter(ht)

	if hinter, err := enc.Decode(bpo); err != nil {
		return e(err, "")
	} else if po, ok := hinter.(Policy); !ok {
		return e(util.ErrWrongType.Errorf("expected Policy, not %T", hinter), "")
	} else {
		cd.policy = po
	}

	return nil
}
