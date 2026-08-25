package app

import (
	"cosmossdk.io/log/v2"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmempool "github.com/cosmos/cosmos-sdk/types/mempool"

	evmmempool "github.com/cosmos/evm/mempool"
	evmconfig "github.com/cosmos/evm/server"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

func (app *App) configureEVMMempool(
	appOpts servertypes.AppOptions,
	logger log.Logger,
) (sdk.PrepareProposalHandler, sdk.ProcessProposalHandler) {
	cosmosPoolMaxTx := evmconfig.GetCosmosPoolMaxTx(appOpts, logger)
	if cosmosPoolMaxTx < 0 || evmtypes.GetChainConfig() == nil {
		// default to nop mempool
		mpool := sdkmempool.NoOpMempool{}
		app.SetMempool(mpool)

		handler := baseapp.NewDefaultProposalHandler(mpool, app)
		return handler.PrepareProposalHandler(), handler.ProcessProposalHandler()
	}

	mpConfig := evmconfig.ResolveMempoolConfig(app.AnteHandler(), appOpts, logger)
	txEncoder := evmmempool.NewTxEncoder(app.txConfig)
	evmRechecker := evmmempool.NewTxRechecker(mpConfig.AnteHandler, txEncoder)
	cosmosRechecker := evmmempool.NewTxRechecker(mpConfig.AnteHandler, txEncoder)

	evmMempool := evmmempool.NewMempool(
		app.CreateQueryContext,
		logger,
		app.EVMKeeper,
		app.FeeMarketKeeper,
		app.txConfig,
		evmRechecker,
		cosmosRechecker,
		mpConfig,
		cosmosPoolMaxTx,
	)
	app.EVMMempool = evmMempool
	app.SetMempool(evmMempool)

	checkTxTimeout := evmconfig.GetMempoolCheckTxTimeout(appOpts, logger)
	app.SetCheckTxHandler(evmMempool.NewCheckTxHandler(app.TxDecode, checkTxTimeout))
	app.SetInsertTxHandler(evmMempool.NewInsertTxHandler(app.TxDecode))
	app.SetReapTxsHandler(evmMempool.NewReapTxsHandler())

	abciProposalHandler := baseapp.NewDefaultProposalHandler(evmMempool, app)
	abciProposalHandler.SetSignerExtractionAdapter(evmmempool.NewEthSignerExtractionAdapter(sdkmempool.NewDefaultSignerExtractionAdapter()))

	return abciProposalHandler.PrepareProposalHandler(), abciProposalHandler.ProcessProposalHandler()
}
