package processor

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-dao/operation/dao"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender   currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency currencytypes.DuplicationType = "currency"
	DuplicationTypeContract currencytypes.DuplicationType = "contract"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var duplicationTypeSenderID string
	var duplicationTypeCurrencyID string
	var duplicationTypeContract string
	var newAddresses []mitumbase.Address

	switch t := op.(type) {
	case currency.CreateAccount:
		fact, ok := t.Fact().(currency.CreateAccountFact)
		if !ok {
			return errors.Errorf("expected CreateAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case currency.UpdateKey:
		fact, ok := t.Fact().(currency.UpdateKeyFact)
		if !ok {
			return errors.Errorf("expected UpdateKeyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Target().String()
	case currency.Transfer:
		fact, ok := t.Fact().(currency.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case currency.RegisterCurrency:
		fact, ok := t.Fact().(currency.RegisterCurrencyFact)
		if !ok {
			return errors.Errorf("expected RegisterCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = fact.Currency().Currency().String()
	case currency.UpdateCurrency:
		fact, ok := t.Fact().(currency.UpdateCurrencyFact)
		if !ok {
			return errors.Errorf("expected UpdateCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Currency().String()
	case currency.Mint:
	case extensioncurrency.CreateContractAccount:
		fact, ok := t.Fact().(extensioncurrency.CreateContractAccountFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case extensioncurrency.Withdraw:
		fact, ok := t.Fact().(extensioncurrency.WithdrawFact)
		if !ok {
			return errors.Errorf("expected WithdrawFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.CreateDAO:
		fact, ok := t.Fact().(dao.CreateDAOFact)
		if !ok {
			return errors.Errorf("expected CreateDAOFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.UpdatePolicy:
		fact, ok := t.Fact().(dao.UpdatePolicyFact)
		if !ok {
			return errors.Errorf("expected UpdatePolicyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.Propose:
		fact, ok := t.Fact().(dao.ProposeFact)
		if !ok {
			return errors.Errorf("expected ProposeFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.CancelProposal:
		fact, ok := t.Fact().(dao.CancelProposalFact)
		if !ok {
			return errors.Errorf("expected CancelProposalFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.Register:
		fact, ok := t.Fact().(dao.RegisterFact)
		if !ok {
			return errors.Errorf("expected RegisterFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.PreSnap:
		fact, ok := t.Fact().(dao.PreSnapFact)
		if !ok {
			return errors.Errorf("expected PreSnapFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.Vote:
		fact, ok := t.Fact().(dao.VoteFact)
		if !ok {
			return errors.Errorf("expected VoteFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.PostSnap:
		fact, ok := t.Fact().(dao.PostSnapFact)
		if !ok {
			return errors.Errorf("expected PostSnapFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case dao.Execute:
		fact, ok := t.Fact().(dao.ExecuteFact)
		if !ok {
			return errors.Errorf("expected ExecuteFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	default:
		return nil
	}

	if len(duplicationTypeSenderID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeSenderID]; found {
			return errors.Errorf("proposal cannot have duplicate sender, %v", duplicationTypeSenderID)
		}

		opr.Duplicated[duplicationTypeSenderID] = DuplicationTypeSender
	}
	if len(duplicationTypeCurrencyID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeCurrencyID]; found {
			return errors.Errorf(
				"cannot register duplicate currency id, %v within a proposal",
				duplicationTypeCurrencyID,
			)
		}

		opr.Duplicated[duplicationTypeCurrencyID] = DuplicationTypeCurrency
	}
	if len(duplicationTypeContract) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContract]; found {
			return errors.Errorf(
				"cannot use a duplicated contract for registering in contract model , %v within a proposal",
				duplicationTypeSenderID,
			)
		}

		opr.Duplicated[duplicationTypeContract] = DuplicationTypeContract
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}

func GetNewProcessor(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) (mitumbase.OperationProcessor, bool, error) {
	switch i, err := opr.GetNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.CreateAccount,
		currency.UpdateKey,
		currency.Transfer,
		extensioncurrency.CreateContractAccount,
		extensioncurrency.Withdraw,
		currency.RegisterCurrency,
		currency.UpdateCurrency,
		currency.Mint,
		dao.CreateDAO,
		dao.UpdatePolicy,
		dao.Propose,
		dao.CancelProposal,
		dao.Register,
		dao.PreSnap,
		dao.Vote,
		dao.PostSnap,
		dao.Execute:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}
