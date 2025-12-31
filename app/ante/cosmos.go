package ante

import (
	"errors"

	corestoretypes "cosmossdk.io/core/store"
	circuitante "cosmossdk.io/x/circuit/ante"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	cosmosante "github.com/cosmos/evm/ante/cosmos"
	evmante "github.com/cosmos/evm/ante/evm"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ibcante "github.com/cosmos/ibc-go/v10/modules/core/ante"
	"github.com/cosmos/ibc-go/v10/modules/core/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by r	equiring the IBC
// channel keeper.
type HandlerOptions struct {
	EvmOptions            EVMHandlerOptions
	IBCKeeper             *keeper.Keeper
	WasmConfig            *wasmTypes.NodeConfig
	WasmKeeper            *wasmkeeper.Keeper
	TXCounterStoreService corestoretypes.KVStoreService
	CircuitKeeper         *circuitkeeper.Keeper
	SanctionKeeper        *sanctionkeeper.Keeper
}

// Validate checks if the keepers are defined
func (options HandlerOptions) Validate() error {
	if options.EvmOptions.Validate() != nil {
		return options.EvmOptions.Validate()
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
func newCosmosAnteHandler(ctx sdk.Context, options HandlerOptions) sdk.AnteHandler {
	feemarketParams := options.EvmOptions.FeeMarketKeeper.GetParams(ctx)
	var txFeeChecker ante.TxFeeChecker
	if options.EvmOptions.DynamicFeeChecker {
		txFeeChecker = evmante.NewDynamicFeeChecker(&feemarketParams)
	}

	anteDecorators := []sdk.AnteDecorator{
		cosmosante.NewRejectMessagesDecorator(), // reject MsgEthereumTxs
		cosmosante.NewAuthzLimiterDecorator( // disable the Msg types that cannot be included on an authz.MsgExec msgs field
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
			sdk.MsgTypeURL(&vestingtypes.MsgCreateVestingAccount{}),
		),
		ante.NewSetUpContextDecorator(),
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		sanctionkeeper.NewBlacklistCheckDecorator(*options.SanctionKeeper),
		ante.NewExtensionOptionsDecorator(options.EvmOptions.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.EvmOptions.AccountKeeper),
		cosmosante.NewMinGasPriceDecorator(&feemarketParams),
		ante.NewConsumeGasForTxSizeDecorator(options.EvmOptions.AccountKeeper),
		ante.NewDeductFeeDecorator(options.EvmOptions.AccountKeeper, options.EvmOptions.BankKeeper, options.EvmOptions.FeegrantKeeper, txFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.EvmOptions.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.EvmOptions.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.EvmOptions.AccountKeeper, options.EvmOptions.SigGasConsumer),
		NewMultiChainIDDecorator(ante.NewSigVerificationDecorator(options.EvmOptions.AccountKeeper, options.EvmOptions.SignModeHandler)),
		ante.NewIncrementSequenceDecorator(options.EvmOptions.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		evmante.NewGasWantedDecorator(options.EvmOptions.EvmKeeper, options.EvmOptions.FeeMarketKeeper, &feemarketParams),
	}

	return sdk.ChainAnteDecorators(anteDecorators...)
}
