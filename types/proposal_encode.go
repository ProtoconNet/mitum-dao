package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (p *CryptoProposal) unpack(enc encoder.Encoder, ht hint.Hint, pr string, st uint64, bcd []byte) error {
	p.BaseHinter = hint.NewBaseHinter(ht)
	p.startTime = st

	switch a, err := base.DecodeAddress(pr, enc); {
	case err != nil:
		return err
	default:
		p.proposer = a
	}

	if hinter, err := enc.Decode(bcd); err != nil {
		return err
	} else if cd, ok := hinter.(CallData); !ok {
		return common.ErrTypeMismatch.Wrap(errors.Errorf("expected CallData, not %T", hinter))
	} else {
		p.callData = cd
	}

	return nil
}

func (p *BizProposal) unpack(enc encoder.Encoder, ht hint.Hint, pr string, st uint64, url, hash string, opt uint8) error {
	e := util.StringError("failed to unmarshal BizProposal")

	p.BaseHinter = hint.NewBaseHinter(ht)

	p.startTime = st
	p.url = URL(url)
	p.hash = hash
	p.options = opt

	switch a, err := base.DecodeAddress(pr, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		p.proposer = a
	}

	return nil
}
