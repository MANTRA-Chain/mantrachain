package app

/*

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (app *App) registerBankModule() {
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
			app.TokenFactoryKeeper.Hooks(),
			// Add any other hooks here, separated by commas
		),
	)

	// Ensure the token factory hooks are always included
	currentHooks := app.BankKeeper.BaseSendKeeper.Hooks()
	if tfHooks, ok := currentHooks.(banktypes.MultiBankHooks); ok {
		tfHooks = append(tfHooks, app.TokenFactoryKeeper.Hooks())
		app.BankKeeper.SetHooks(tfHooks)
	} else {
		app.Logger().Error("Failed to add TokenFactory hooks to BankKeeper")
	}
}

*/
