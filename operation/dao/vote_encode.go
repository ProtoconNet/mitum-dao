package dao

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *VoteFact) unpack(enc encoder.Encoder,
	sa, ca, pid string, vt uint8, cid string,
) error {
	fact.proposalID = pid
	fact.vote = vt
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

	return nil
}
