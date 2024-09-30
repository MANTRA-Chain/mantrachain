package app

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// registerBankModule sets up the bank module for the app.
func (app *App) registerBankModule() error {
	// Initialize the bank keeper
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		app.AppCodec(),
		runtime.NewKVStoreService(app.GetKey(banktypes.StoreKey)),
		app.AccountKeeper,
		BlockedAddresses(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.Logger(),
	)

	// Set the hooks on the BankKeeper
	app.BankKeeper.SetHooks(
		banktypes.NewMultiBankHooks(
			// Add any other hooks here
			app.TokenFactoryKeeper.Hooks(),
		),
	)

	return nil
}
