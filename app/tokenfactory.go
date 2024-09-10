package app

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/osmosis-labs/osmosis/v26/x/tokenfactory"
	tokenfactorykeeper "github.com/osmosis-labs/osmosis/v26/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/osmosis-labs/osmosis/v26/x/tokenfactory/types"
)

func (app *App) registerTokenFactoryModule() error {
	// Set up store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(tokenfactorytypes.StoreKey),
	); err != nil {
		return err
	}

	// Create TokenFactory Keeper
	app.TokenFactoryKeeper = tokenfactorykeeper.NewKeeper(
		app.appCodec,
		app.GetKey(tokenfactorytypes.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create TokenFactory Module
	tokenFactoryModule := tokenfactory.NewAppModule(app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper)

	// Register the TokenFactory module
	if err := app.RegisterModules(tokenFactoryModule); err != nil {
		return fmt.Errorf("failed to register TokenFactory module: %w", err)
	}

	// Add TokenFactory module to the begin blocker
	app.ModuleManager.SetOrderBeginBlockers(tokenfactorytypes.ModuleName)

	// Add TokenFactory module to the end blocker
	app.ModuleManager.SetOrderEndBlockers(tokenfactorytypes.ModuleName)

	// Add TokenFactory module to the init genesis
	app.ModuleManager.SetOrderInitGenesis(tokenfactorytypes.ModuleName)

	return nil
}

// RegisterTokenFactory registers the TokenFactory module's types for the given codec.
func RegisterTokenFactory(registry codec.InterfaceRegistry) {
	tokenfactory.RegisterInterfaces(registry)
}
