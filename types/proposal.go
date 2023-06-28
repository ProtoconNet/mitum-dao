package types

import (
	"net/url"
	"strings"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type URL string

func (u URL) IsValid([]byte) error {
	if _, err := url.Parse(string(u)); err != nil {
		return err
	}

	if u != "" && strings.TrimSpace(string(u)) == "" {
		return util.ErrInvalid.Errorf("empty url")
	}

	return nil
}

func (u URL) Bytes() []byte {
	return []byte(u)
}

func (u URL) String() string {
	return string(u)
}

const (
	ProposalCrypto = "crypto"
	ProposalBiz    = "biz"
)

var (
	CryptoProposalHint = hint.MustNewHint("mitum-dao-crypto-proposal-v0.0.1")
	BizProposalHint    = hint.MustNewHint("mitum-dao-biz-proposal-v0.0.1")
)

type Proposal interface {
	util.IsValider
	hint.Hinter
	Type() string
	Bytes() []byte
	StartTime() uint64
	Options() uint8
	Addresses() []base.Address
}

type CryptoProposal struct {
	hint.BaseHinter
	starttime uint64
	calldata  Calldata
}

func NewCryptoProposal(starttime uint64, calldata Calldata) CryptoProposal {
	return CryptoProposal{
		BaseHinter: hint.NewBaseHinter(CryptoProposalHint),
		starttime:  starttime,
		calldata:   calldata,
	}
}

func (CryptoProposal) Type() string {
	return ProposalCrypto
}

func (CryptoProposal) Options() uint8 {
	return 3
}

func (p CryptoProposal) Bytes() []byte {
	return util.ConcatBytesSlice(util.Uint64ToBytes(p.starttime), p.calldata.Bytes())
}

func (p CryptoProposal) StartTime() uint64 {
	return p.starttime
}

func (p CryptoProposal) Calldata() Calldata {
	return p.calldata
}

func (p CryptoProposal) IsValid([]byte) error {
	if err := p.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := p.calldata.IsValid(nil); err != nil {
		return err
	}

	return nil
}

func (p CryptoProposal) Addresses() []base.Address {
	return p.calldata.Addresses()
}

type BizProposal struct {
	hint.BaseHinter
	starttime uint64
	url       URL
	hash      string
	options   uint8
}

func NewBizProposal(starttime uint64, url URL, hash string, options uint8) BizProposal {
	return BizProposal{
		BaseHinter: hint.NewBaseHinter(BizProposalHint),
		starttime:  starttime,
		url:        url,
		hash:       hash,
	}
}

func (BizProposal) Type() string {
	return ProposalBiz
}

func (p BizProposal) Options() uint8 {
	return p.options
}

func (p BizProposal) Bytes() []byte {
	return util.ConcatBytesSlice(util.Uint64ToBytes(p.starttime), p.url.Bytes(), []byte(p.hash), util.Uint8ToBytes(p.options))
}

func (p BizProposal) StartTime() uint64 {
	return p.starttime
}

func (p BizProposal) Url() URL {
	return p.url
}

func (p BizProposal) Hash() string {
	return p.hash
}

func (p BizProposal) IsValid([]byte) error {
	if err := p.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := p.url.IsValid(nil); err != nil {
		return err
	}

	if len(p.hash) == 0 {
		return util.ErrInvalid.Errorf("biz - empty hash")
	}

	return nil
}

func (p BizProposal) Addresses() []base.Address {
	return []base.Address{}
}
