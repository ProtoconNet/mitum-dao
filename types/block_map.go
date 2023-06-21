package types

import "github.com/ProtoconNet/mitum2/base"

type GetLastBlockFunc func() (base.BlockMap, bool, error)
