package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	"github.com/ProtoconNet/mitum-dao/operation/dao"
	"github.com/ProtoconNet/mitum-dao/operation/processor"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-dao-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database
	var opr *currencyprocessor.OperationProcessor
	var set *hint.CompatibleSet[isaac.NewOperationProcessorInternalFunc]

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
		currencycmds.OperationProcessorContextKey, &opr,
		launch.OperationProcessorsMapContextKey, &set,
	); err != nil {
		return pctx, err
	}

	//err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	//if err != nil {
	//	return pctx, err
	//}
	err := opr.SetGetNewProcessorFunc(processor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}

	if err := opr.SetProcessor(
		dao.CreateDAOHint,
		dao.NewCreateDAOProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.UpdatePolicyHint,
		dao.NewUpdatePolicyProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		dao.ProposeHint,
		dao.NewProposeProcessor(),
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
		dao.NewExecuteProcessor(db.LastBlockMap),
	); err != nil {
		return pctx, err
	}

	_ = set.Add(dao.CreateDAOHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.UpdatePolicyHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.ProposeHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.CancelProposalHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.RegisterHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.PreSnapHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.VoteHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.PostSnapHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(dao.ExecuteHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	pctx = context.WithValue(pctx, currencycmds.OperationProcessorContextKey, opr)
	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, set) //revive:disable-line:modifies-parameter

	return pctx, nil
}
