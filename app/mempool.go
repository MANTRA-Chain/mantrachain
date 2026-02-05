package app

import (
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmempool "github.com/cosmos/cosmos-sdk/types/mempool"

	evmconfig "github.com/cosmos/evm/config"
	evmmempool "github.com/cosmos/evm/mempool"
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

	mempoolConfig := &evmmempool.EVMMempoolConfig{
		AnteHandler:      app.AnteHandler(),
		LegacyPoolConfig: evmconfig.GetLegacyPoolConfig(appOpts, logger),
		BlockGasLimit:    evmconfig.GetBlockGasLimit(appOpts, logger),
		MinTip:           evmconfig.GetMinTip(appOpts, logger),
	}

	evmMempool := evmmempool.NewExperimentalEVMMempool(app.CreateQueryContext, logger, app.EVMKeeper, app.FeeMarketKeeper, app.txConfig, app.clientCtx, mempoolConfig, cosmosPoolMaxTx)
	app.EVMMempool = evmMempool
	app.SetMempool(evmMempool)

	checkTxHandler := evmmempool.NewCheckTxHandler(evmMempool)
	app.SetCheckTxHandler(checkTxHandler)

	abciProposalHandler := baseapp.NewDefaultProposalHandler(evmMempool, app)
	abciProposalHandler.SetSignerExtractionAdapter(evmmempool.NewEthSignerExtractionAdapter(sdkmempool.NewDefaultSignerExtractionAdapter()))

	return abciProposalHandler.PrepareProposalHandler(), baseapp.NoOpProcessProposal()
}
