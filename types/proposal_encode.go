package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (p *CryptoProposal) unpack(enc encoder.Encoder, ht hint.Hint, st uint64, bcd []byte) error {
	e := util.StringError("failed to decode bson of CryptoProposal")

	p.BaseHinter = hint.NewBaseHinter(ht)
	p.startTime = st

	if hinter, err := enc.Decode(bcd); err != nil {
		return e.Wrap(err)
	} else if cd, ok := hinter.(CallData); !ok {
		return e.Wrap(errors.Errorf("expected CallData, not %T", hinter))
	} else {
		p.callData = cd
	}

	return nil
}

func (p *BizProposal) unpack(_ encoder.Encoder, ht hint.Hint, st uint64, url, hash string) error {
	p.BaseHinter = hint.NewBaseHinter(ht)

	p.startTime = st
	p.url = URL(url)
	p.hash = hash

	return nil
}
