package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *Design) unpack(enc encoder.Encoder, ht hint.Hint, op, dao string, bpo []byte) error {
	e := util.StringErrorFunc("failed to decode bson of Design")

	de.BaseHinter = hint.NewBaseHinter(ht)
	de.option = DAOOption(dao)
	de.daoID = types.ContractID(dao)

	if hinter, err := enc.Decode(bpo); err != nil {
		return e(err, "")
	} else if po, ok := hinter.(Policy); !ok {
		return e(util.ErrWrongType.Errorf("expected Policy, not %T", hinter), "")
	} else {
		de.policy = po
	}

	return nil
}
