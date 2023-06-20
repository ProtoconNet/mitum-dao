package types

import (
	"net/url"
	"strings"

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
	Calldata() Calldata
}

type CryptoProposal struct {
	hint.BaseHinter
	calldata Calldata
}

func NewCryptoProposal(calldata Calldata) CryptoProposal {
	return CryptoProposal{
		BaseHinter: hint.NewBaseHinter(CryptoProposalHint),
		calldata:   calldata,
	}
}

func (CryptoProposal) Type() string {
	return ProposalCrypto
}

func (p CryptoProposal) Bytes() []byte {
	return util.ConcatBytesSlice(p.calldata.Bytes())
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

type BizProposal struct {
	hint.BaseHinter
	url  URL
	hash string
}

func NewBizProposal(url URL, hash string) BizProposal {
	return BizProposal{
		BaseHinter: hint.NewBaseHinter(BizProposalHint),
		url:        url,
		hash:       hash,
	}
}

func (BizProposal) Type() string {
	return ProposalBiz
}

func (p BizProposal) Bytes() []byte {
	return util.ConcatBytesSlice(p.url.Bytes(), []byte(p.hash))
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
