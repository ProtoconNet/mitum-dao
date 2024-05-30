package dao

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *PreSnapFact) unpack(enc encoder.Encoder,
	sa, ca, pid, cid string,
) error {
	fact.proposalID = pid
	fact.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:

	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:

	default:
		fact.contract = a
	}

	return nil
}
