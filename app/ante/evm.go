package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmante "github.com/cosmos/evm/ante/evm"
)

func newEVMAnteHandler(options EVMHandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		EVMMsgCheckDecorator{},
		evmante.NewEVMMonoDecorator(
			options.AccountKeeper,
			options.FeeMarketKeeper,
			options.EvmKeeper,
			options.MaxTxGasWanted,
		),
	)
}
