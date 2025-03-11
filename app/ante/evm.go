package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmosante "github.com/evmos/evmos/v20/app/ante"
	evmante "github.com/evmos/evmos/v20/app/ante/evm"
)

func newMonoEVMAnteHandler(options evmosante.HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		evmante.NewMonoDecorator(
			options.AccountKeeper,
			options.BankKeeper,
			options.FeeMarketKeeper,
			options.EvmKeeper,
			options.DistributionKeeper,
			options.StakingKeeper,
			options.MaxTxGasWanted,
		),
	)
}
