package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (p *CryptoProposal) unpack(enc encoder.Encoder, ht hint.Hint, bcd []byte) error {
	e := util.StringErrorFunc("failed to decode bson of CryptoProposal")

	if hinter, err := enc.Decode(bcd); err != nil {
		return e(err, "")
	} else if cd, ok := hinter.(Calldata); !ok {
		return e(util.ErrWrongType.Errorf("expected Calldata, not %T", hinter), "")
	} else {
		p.calldata = cd
	}

	return nil
}

func (p *BizProposal) unpack(enc encoder.Encoder, ht hint.Hint, url, hash string) error {
	p.url = URL(url)
	p.hash = hash

	return nil
}
