package dao

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *ProposeFact) unpack(enc encoder.Encoder,
	sa, ca, pid string,
	bp []byte,
	cid string,
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

	if hinter, err := enc.Decode(bp); err != nil {
		return err
	} else if proposal, ok := hinter.(types.Proposal); !ok {
		return common.ErrTypeMismatch.Wrap(errors.Errorf("expected Proposal, not %T", hinter))
	} else {
		fact.proposal = proposal
	}

	return nil
}
