package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func newMonoEVMAnteHandler(options EVMHandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewEVMMonoDecorator(
			options.AccountKeeper,
			options.FeeMarketKeeper,
			options.EvmKeeper,
			options.MaxTxGasWanted,
		),
	)
}
