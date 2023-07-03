package digest

import (
	"context"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	isaacblock "github.com/ProtoconNet/mitum2/isaac/block"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/fixedtree"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"sort"
	"sync"
	"time"
)

type Digester struct {
	sync.RWMutex
	*util.ContextDaemon
	*logging.Logging
	database    *currencydigest.Database
	localfsRoot string
	blockChan   chan base.BlockMap
	errChan     chan error
}

func NewDigester(st *currencydigest.Database, root string, errChan chan error) *Digester {
	di := &Digester{
		Logging: logging.NewLogging(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "digester")
		}),
		database:    st,
		localfsRoot: root,
		blockChan:   make(chan base.BlockMap, 100),
		errChan:     errChan,
	}

	di.ContextDaemon = util.NewContextDaemon(di.start)

	return di
}

func (di *Digester) start(ctx context.Context) error {
	errch := func(err currencydigest.DigestError) {
		if di.errChan == nil {
			return
		}

		di.errChan <- err
	}

end:
	for {
		select {
		case <-ctx.Done():
			di.Log().Debug().Msg("stopped")

			break end
		case blk := <-di.blockChan:
			err := util.Retry(ctx, func() (bool, error) {
				if err := di.digest(ctx, blk); err != nil {
					go errch(currencydigest.NewDigestError(err, blk.Manifest().Height()))

					if errors.Is(err, context.Canceled) {
						return false, isaac.ErrStopProcessingRetry.Wrap(err)
					}

					return true, err
				}

				return false, nil
			}, 15, time.Second*1)
			if err != nil {
				di.Log().Error().Err(err).Int64("block", blk.Manifest().Height().Int64()).Msg("failed to digest block")
			} else {
				di.Log().Info().Int64("block", blk.Manifest().Height().Int64()).Msg("block digested")
			}

			go errch(currencydigest.NewDigestError(err, blk.Manifest().Height()))
		}
	}

	return nil
}

func (di *Digester) Digest(blocks []base.BlockMap) {
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Manifest().Height() < blocks[j].Manifest().Height()
	})

	for i := range blocks {
		blk := blocks[i]
		di.Log().Debug().Int64("block", blk.Manifest().Height().Int64()).Msg("start to digest block")

		di.blockChan <- blk
	}
}

func (di *Digester) digest(ctx context.Context, blk base.BlockMap) error {
	di.Lock()
	defer di.Unlock()
	reader, err := isaacblock.NewLocalFSReaderFromHeight(di.localfsRoot, blk.Manifest().Height(), di.database.DatabaseEncoders().Find(jsonenc.JSONEncoderHint))
	if err != nil {
		return err
	}
	var ops []base.Operation
	switch v, found, err := reader.Item(base.BlockMapItemTypeOperations); {
	case err != nil:
		return err
	case found:
		ops = v.([]base.Operation) //nolint:forcetypeassert //...
	}
	var opstree fixedtree.Tree
	switch v, found, err := reader.Item(base.BlockMapItemTypeOperationsTree); {
	case err != nil:
		return err
	case found:
		opstree = v.(fixedtree.Tree) //nolint:forcetypeassert //...
	}
	var sts []base.State
	switch v, found, err := reader.Item(base.BlockMapItemTypeStates); {
	case err != nil:
		return err
	case found:
		sts = v.([]base.State) //nolint:forcetypeassert //...
	}

	if err := DigestBlock(ctx, di.database, blk, ops, opstree, sts); err != nil {
		return err
	}

	return di.database.SetLastBlock(blk.Manifest().Height())
}

func DigestBlock(ctx context.Context, st *currencydigest.Database, blk base.BlockMap, ops []base.Operation, opstree fixedtree.Tree, sts []base.State) error {
	bs, err := NewBlockSession(st, blk, ops, opstree, sts)
	if err != nil {
		return err
	}
	defer func() {
		_ = bs.Close()
	}()

	if err := bs.Prepare(); err != nil {
		return err
	}

	return bs.Commit(ctx)
}
