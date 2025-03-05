package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	evmosante "github.com/evmos/evmos/v20/app/ante"
	ethante "github.com/evmos/evmos/v20/app/ante/evm"
	evmostypes "github.com/evmos/evmos/v20/types"
)

type EvmosAnteHandlerOptions evmosante.HandlerOptions

func NewEvmosAnteHandlerOptionsFromApp(app *App, txConfig client.TxConfig, maxGasWanted uint64) *EvmosAnteHandlerOptions {
	return &EvmosAnteHandlerOptions{
		Cdc:                    app.appCodec,
		AccountKeeper:          app.AccountKeeper,
		BankKeeper:             app.BankKeeper,
		ExtensionOptionChecker: evmostypes.HasDynamicFeeExtensionOption,
		EvmKeeper:              app.EvmKeeper,
		FeegrantKeeper:         app.FeeGrantKeeper,
		IBCKeeper:              app.IBCKeeper,
		FeeMarketKeeper:        app.FeeMarketKeeper,
		SignModeHandler:        txConfig.SignModeHandler(),
		SigGasConsumer:         evmosante.SigVerificationGasConsumer,
		MaxTxGasWanted:         maxGasWanted,
		TxFeeChecker:           ethante.NewDynamicFeeChecker(app.EvmKeeper),
		StakingKeeper:          app.StakingKeeper,
		DistributionKeeper:     app.DistrKeeper,
	}
}

func (aa *EvmosAnteHandlerOptions) Validate() error {
	return (*evmosante.HandlerOptions)(aa).Validate()
}

func (aa *EvmosAnteHandlerOptions) Options() evmosante.HandlerOptions {
	return evmosante.HandlerOptions(*aa)
}

func (aa *EvmosAnteHandlerOptions) WithCodec(cdc codec.BinaryCodec) *EvmosAnteHandlerOptions {
	aa.Cdc = cdc
	return aa
}

func (aa *EvmosAnteHandlerOptions) WithMaxTxGasWanted(maxTxGasWanted uint64) *EvmosAnteHandlerOptions {
	aa.MaxTxGasWanted = maxTxGasWanted
	return aa
}
