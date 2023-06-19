package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-dao/types"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	DAO types.Design `json:"dao"`
}

func (de DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignStateValueJSONMarshaler{
		BaseHinter: de.BaseHinter,
		DAO:        de.Design,
	})
}

type DesignStateValueJSONUnmarshaler struct {
	DAO json.RawMessage `json:"dao"`
}

func (de *DesignStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringErrorFunc("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e(err, "")
	}

	var design types.Design

	if err := design.DecodeJSON(u.DAO, enc); err != nil {
		return e(err, "")
	}

	de.Design = design

	return nil
}
