package app

import (
	"github.com/MANTRA-Chain/mantrachain/v4/app/ante"
	"github.com/cosmos/cosmos-sdk/client"
	evmante "github.com/cosmos/evm/ante"
	cosmosevmante "github.com/cosmos/evm/ante/evm"
	cosmosevmtypes "github.com/cosmos/evm/types"
)

func NewEVMAnteHandlerOptionsFromApp(app *App, txConfig client.TxConfig, maxGasWanted uint64) *ante.EVMHandlerOptions {
	return &ante.EVMHandlerOptions{
		Cdc:                    app.appCodec,
		AccountKeeper:          app.AccountKeeper,
		BankKeeper:             app.BankKeeper,
		ExtensionOptionChecker: cosmosevmtypes.HasDynamicFeeExtensionOption,
		EvmKeeper:              app.EVMKeeper,
		FeegrantKeeper:         app.FeeGrantKeeper,
		IBCKeeper:              app.IBCKeeper,
		FeeMarketKeeper:        app.FeeMarketKeeper,
		SignModeHandler:        txConfig.SignModeHandler(),
		SigGasConsumer:         evmante.SigVerificationGasConsumer,
		MaxTxGasWanted:         maxGasWanted,
		TxFeeChecker:           cosmosevmante.NewDynamicFeeChecker(app.FeeMarketKeeper),
	}
}
