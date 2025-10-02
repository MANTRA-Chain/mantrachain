package ante

import (
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v6/x/sanction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmante "github.com/cosmos/evm/ante/evm"
	chainante "github.com/cosmos/evm/evmd/ante"
)

func newEVMAnteHandler(options HandlerOptions) sdk.AnteHandler {
	decorators := []sdk.AnteDecorator{
		sanctionkeeper.NewEVMBlacklistCheckDecorator(*options.SanctionKeeper),
		evmante.NewEVMMonoDecorator(
			options.EvmOptions.AccountKeeper,
			options.EvmOptions.FeeMarketKeeper,
			options.EvmOptions.EvmKeeper,
			options.EvmOptions.MaxTxGasWanted,
		),
	}
	if options.EvmOptions.PendingTxListener != nil {
		decorators = append(decorators, chainante.NewTxListenerDecorator(options.EvmOptions.PendingTxListener))
	}
	return sdk.ChainAnteDecorators(decorators...)
}
