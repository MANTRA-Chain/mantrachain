package app

import (
	"fmt"

	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
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

	// Ensure the subspace is properly initialized
	app.ParamsKeeper.Subspace(tokenfactorytypes.ModuleName)

	// Create TokenFactory Keeper
	tokenFactoryKeeper := tokenfactorykeeper.NewKeeper(
		app.GetKey(tokenfactorytypes.StoreKey),
		app.GetSubspace(tokenfactorytypes.ModuleName),
		GetMaccPerms(),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
	)
	app.TokenFactoryKeeper = tokenFactoryKeeper

	// Create TokenFactory Module
	tokenFactoryModule := tokenfactory.NewAppModule(app.TokenFactoryKeeper, app.AccountKeeper, app.BankKeeper)

	// Register the TokenFactory module
	if err := app.RegisterModules(tokenFactoryModule); err != nil {
		return fmt.Errorf("failed to register TokenFactory module: %w", err)
	}

	return nil
}

// RegisterTokenFactory registers the TokenFactory module with the given interface registry.
func RegisterTokenFactory(registry codectypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		tokenfactorytypes.ModuleName: tokenfactory.AppModule{},
	}

	for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}

	return modules
}
