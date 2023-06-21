package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (p *CryptoProposal) unpack(enc encoder.Encoder, ht hint.Hint, st uint64, bcd []byte) error {
	e := util.StringErrorFunc("failed to decode bson of CryptoProposal")

	p.BaseHinter = hint.NewBaseHinter(ht)
	p.starttime = st

	if hinter, err := enc.Decode(bcd); err != nil {
		return e(err, "")
	} else if cd, ok := hinter.(Calldata); !ok {
		return e(util.ErrWrongType.Errorf("expected Calldata, not %T", hinter), "")
	} else {
		p.calldata = cd
	}

	return nil
}

func (p *BizProposal) unpack(enc encoder.Encoder, ht hint.Hint, st uint64, url, hash string) error {
	p.BaseHinter = hint.NewBaseHinter(ht)

	p.starttime = st
	p.url = URL(url)
	p.hash = hash

	return nil
}
