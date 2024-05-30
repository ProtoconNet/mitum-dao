package dao

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RegisterFact) unpack(enc encoder.Encoder,
	sa, ca, pid, ta, cid string,
) error {
	fact.proposalID = pid
	fact.currency = currencytypes.CurrencyID(cid)

	switch a, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return err
	default:
		fact.sender = a
	}

	switch a, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		fact.contract = a
	}

	if ta != "" {
		switch a, err := base.DecodeAddress(ta, enc); {
		case err != nil:
			return err
		default:
			fact.delegated = a
		}
	} else {
		fact.delegated = nil
	}

	return nil
}
