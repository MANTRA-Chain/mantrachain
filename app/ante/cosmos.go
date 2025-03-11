package ante

import (
	"errors"

	corestoretypes "cosmossdk.io/core/store"
	circuitante "cosmossdk.io/x/circuit/ante"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v4/x/sanction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	"github.com/cosmos/ibc-go/v8/modules/core/keeper"
	evmosante "github.com/evmos/evmos/v20/app/ante"
	evmoscosmosante "github.com/evmos/evmos/v20/app/ante/cosmos"
	evmante "github.com/evmos/evmos/v20/app/ante/evm"
	evmtypes "github.com/evmos/evmos/v20/x/evm/types"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	EvmosOptions          evmosante.HandlerOptions
	IBCKeeper             *keeper.Keeper
	WasmConfig            *wasmTypes.NodeConfig
	WasmKeeper            *wasmkeeper.Keeper
	TXCounterStoreService corestoretypes.KVStoreService
	CircuitKeeper         *circuitkeeper.Keeper
	SanctionKeeper        *sanctionkeeper.Keeper
}

// Validate checks if the keepers are defined
func (options HandlerOptions) Validate() error {
	if options.EvmosOptions.Validate() != nil {
		return options.EvmosOptions.Validate()
	}
	if options.IBCKeeper == nil {
		return errors.New("ibc keeper is required for ante builder")
	}
	if options.WasmConfig == nil {
		return errors.New("wasm config is required for ante builder")
	}
	if options.WasmKeeper == nil {
		return errors.New("wasm keeper is required for ante builder")
	}
	if options.TXCounterStoreService == nil {
		return errors.New("wasm store service is required for ante builder")
	}
	if options.CircuitKeeper == nil {
		return errors.New("circuit keeper is required for ante builder")
	}
	if options.SanctionKeeper == nil {
		return errors.New("sanction keeper is required for ante builder")
	}
	return nil
}

// newCosmosAnteHandler constructor
func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	anteDecorators := []sdk.AnteDecorator{
		evmoscosmosante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		evmoscosmosante.NewAuthzLimiterDecorator( // disable the Msg types that cannot be included on an authz.MsgExec msgs field
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
			sdk.MsgTypeURL(&vestingtypes.MsgCreateVestingAccount{}),
		),
		ante.NewSetUpContextDecorator(),
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		sanctionkeeper.NewBlacklistCheckDecorator(*options.SanctionKeeper),
		ante.NewExtensionOptionsDecorator(options.EvmosOptions.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.EvmosOptions.AccountKeeper),
		evmoscosmosante.NewMinGasPriceDecorator(options.EvmosOptions.FeeMarketKeeper, options.EvmosOptions.EvmKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.EvmosOptions.AccountKeeper),
		ante.NewDeductFeeDecorator(
			options.EvmosOptions.AccountKeeper,
			options.EvmosOptions.BankKeeper,
			options.EvmosOptions.FeegrantKeeper,
			options.EvmosOptions.TxFeeChecker,
		),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.EvmosOptions.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.EvmosOptions.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.EvmosOptions.AccountKeeper, options.EvmosOptions.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.EvmosOptions.AccountKeeper, options.EvmosOptions.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.EvmosOptions.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		evmante.NewGasWantedDecorator(options.EvmosOptions.EvmKeeper, options.EvmosOptions.FeeMarketKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...)
}
