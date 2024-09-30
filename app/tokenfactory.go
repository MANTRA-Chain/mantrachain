package app

import (
	tokenfactory "github.com/MANTRA-Chain/mantrachain/x/tokenfactory"
	tokenfactorykeeper "github.com/MANTRA-Chain/mantrachain/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/MANTRA-Chain/mantrachain/x/tokenfactory/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (app *App) registerTokenFactoryModule(
	appCodec codec.Codec,
) {
	// Initialize the token factory keeper
	tokenFactoryKeeper := tokenfactorykeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(app.GetKey(tokenfactorytypes.StoreKey)),
		knownModules(),
		app.AccountKeeper,
		app.BankKeeper,
		app.WasmKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.TokenFactoryKeeper = tokenFactoryKeeper

	// Set the hooks on the BankKeeper
	app.BankKeeper.SetHooks(
		banktypes.NewMultiBankHooks(
			app.TokenFactoryKeeper.Hooks(),
		),
	)

	// Add TokenFactoryKeeper to the app's module manager
	app.ModuleManager.Modules[tokenfactorytypes.ModuleName] = tokenfactory.NewAppModule(appCodec, app.TokenFactoryKeeper)
}
