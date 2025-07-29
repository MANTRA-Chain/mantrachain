package ante

import (
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v5/x/sanction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmante "github.com/cosmos/evm/ante/evm"
)

func newEVMAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		sanctionkeeper.NewEVMBlacklistCheckDecorator(*options.SanctionKeeper),
		evmante.NewEVMMonoDecorator(
			options.EvmOptions.AccountKeeper,
			options.EvmOptions.FeeMarketKeeper,
			options.EvmOptions.EvmKeeper,
			options.EvmOptions.MaxTxGasWanted,
		),
	)
}
