package dao

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *ProposeFact) unpack(enc encoder.Encoder,
	sa, ca, did, pid string,
	st uint64,
	bp []byte,
	cid string,
) error {
	e := util.StringErrorFunc("failed to unmarshal ProposeFact")

	fact.daoID = currencytypes.ContractID(did)
	fact.proposalID = pid
	fact.startTime = st
	fact.currency = currencytypes.CurrencyID(cid)

	if hinter, err := enc.Decode(bp); err != nil {
		return e(err, "")
	} else if proposal, ok := hinter.(types.Proposal); !ok {
		return e(util.ErrWrongType.Errorf("expected Proposal, not %T", hinter), "")
	} else {
		fact.proposal = proposal
	}

	return nil
}
