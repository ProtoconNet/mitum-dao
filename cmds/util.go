package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v3/operation/isaac"
	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum-dao/operation/processor"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
	); err != nil {
		return pctx, err
	}

	limiterF, err := launch.NewSuffrageCandidateLimiterFunc(pctx)
	if err != nil {
		return pctx, err
	}

	set := hint.NewCompatibleSet()

	opr := processor.NewOperationProcessor()
	if err := opr.SetProcessor(
		currency.CreateAccountsHint,
		currency.NewCreateAccountsProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		currency.KeyUpdaterHint,
		currency.NewKeyUpdaterProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		currency.TransfersHint,
		currency.NewTransfersProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		currency.CurrencyRegisterHint,
		currency.NewCurrencyRegisterProcessor(isaacParams.Threshold()),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		currency.CurrencyPolicyUpdaterHint,
		currency.NewCurrencyPolicyUpdaterProcessor(isaacParams.Threshold()),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		currency.SuffrageInflationHint,
		currency.NewSuffrageInflationProcessor(isaacParams.Threshold()),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		extension.CreateContractAccountsHint,
		extension.NewCreateContractAccountsProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		extension.WithdrawsHint,
		extension.NewWithdrawsProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.CreateDAOHint,
		dao.NewCreateDAOProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.ProposeHint,
		dao.NewProposeProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.ProposeHint,
		dao.NewCancelProposalProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.CancelProposalHint,
		dao.NewCancelProposalProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.RegisterHint,
		dao.NewRegisterProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.PreSnapHint,
		dao.NewPreSnapProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.VoteHint,
		dao.NewVoteProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.PostSnapHint,
		dao.NewPostSnapProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.ExecuteHint,
		dao.NewPostSnapProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	}

	_ = set.Add(currency.CreateAccountsHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.KeyUpdaterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.TransfersHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.CurrencyRegisterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.CurrencyPolicyUpdaterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(currency.SuffrageInflationHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(extension.CreateContractAccountsHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(extension.WithdrawsHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.CreateDAOHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.ProposeHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.CancelProposalHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.RegisterHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.PreSnapHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.VoteHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.PostSnapHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(dao.ExecuteHint, func(height base.Height) (base.OperationProcessor, error) {
		return opr.New(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(isaacoperation.SuffrageCandidateHint, func(height base.Height) (base.OperationProcessor, error) {
		policy := db.LastNetworkPolicy()
		if policy == nil { // NOTE Usually it means empty block data
			return nil, nil
		}

		return isaacoperation.NewSuffrageCandidateProcessor(
			height,
			db.State,
			limiterF,
			nil,
			policy.SuffrageCandidateLifespan(),
		)
	})

	_ = set.Add(isaacoperation.SuffrageJoinHint, func(height base.Height) (base.OperationProcessor, error) {
		policy := db.LastNetworkPolicy()
		if policy == nil { // NOTE Usually it means empty block data
			return nil, nil
		}

		return isaacoperation.NewSuffrageJoinProcessor(
			height,
			isaacParams.Threshold(),
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(isaac.SuffrageExpelOperationHint, func(height base.Height) (base.OperationProcessor, error) {
		policy := db.LastNetworkPolicy()
		if policy == nil { // NOTE Usually it means empty block data
			return nil, nil
		}

		return isaacoperation.NewSuffrageExpelProcessor(
			height,
			db.State,
			nil,
			nil,
		)
	})

	_ = set.Add(isaacoperation.SuffrageDisjoinHint, func(height base.Height) (base.OperationProcessor, error) {
		return isaacoperation.NewSuffrageDisjoinProcessor(
			height,
			db.State,
			nil,
			nil,
		)
	})

	var f currencycmds.ProposalOperationFactHintFunc = IsSupportedProposalOperationFactHintFunc

	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, set) //revive:disable-line:modifies-parameter
	pctx = context.WithValue(pctx, currencycmds.ProposalOperationFactHintContextKey, f)

	return pctx, nil
}

func IsSupportedProposalOperationFactHintFunc() func(hint.Hint) bool {
	return func(ht hint.Hint) bool {
		for i := range SupportedProposalOperationFactHinters {
			s := SupportedProposalOperationFactHinters[i].Hint
			if ht.Type() != s.Type() {
				continue
			}

			return ht.IsCompatible(s)
		}

		return false
	}
}
