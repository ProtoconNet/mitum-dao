package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) handleAccountState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if rs, err := currencydigest.NewAccountValue(st); err != nil {
		return nil, err
	} else if doc, err := currencydigest.NewAccountDoc(rs, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, nil
	}
}

func (bs *BlockSession) handleBalanceState(st mitumbase.State) ([]mongo.WriteModel, string, error) {
	doc, address, err := currencydigest.NewBalanceDoc(st, bs.st.DatabaseEncoder())
	if err != nil {
		return nil, "", err
	}
	return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, address, nil
}

func (bs *BlockSession) handleContractAccountState(st mitumbase.State) ([]mongo.WriteModel, error) {
	doc, err := currencydigest.NewContractAccountStatusDoc(st, bs.st.DatabaseEncoder())
	if err != nil {
		return nil, err
	}
	return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, nil
}

func (bs *BlockSession) handleCurrencyState(st mitumbase.State) ([]mongo.WriteModel, error) {
	doc, err := currencydigest.NewCurrencyDoc(st, bs.st.DatabaseEncoder())
	if err != nil {
		return nil, err
	}
	return []mongo.WriteModel{mongo.NewInsertOneModel().SetDocument(doc)}, nil
}
