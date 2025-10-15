package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	icstest "github.com/cosmos/interchain-security/v7/testutil/integration"
	ibcproviderkeeper "github.com/cosmos/interchain-security/v7/x/ccv/provider/keeper"
)

// TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *App) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetTxConfig implements the TestingApp interface.
func (app *App) GetTxConfig() client.TxConfig {
	return app.txConfig
}

// GetTestGovKeeper implements the TestingApp interface.
func (app *App) GetTestGovKeeper() *govkeeper.Keeper {
	return app.GovKeeper
}

func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *App) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

func (app *App) GetStakingKeeper() *stakingkeeper.Keeper {
	return app.StakingKeeper
}

func (app *App) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *App) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}

// ProviderApp interface implementations for icstest tests

// GetProviderKeeper implements the ProviderApp interface.
func (app *App) GetProviderKeeper() ibcproviderkeeper.Keeper { //nolint:nolintlint
	return app.ProviderKeeper
}

// GetTestStakingKeeper implements the ProviderApp interface.
func (app *App) GetTestStakingKeeper() icstest.TestStakingKeeper { //nolint:nolintlint
	return app.StakingKeeper
}

// GetTestBankKeeper implements the ProviderApp interface.
func (app *App) GetTestBankKeeper() icstest.TestBankKeeper { //nolint:nolintlint
	return app.BankKeeper
}

// GetTestSlashingKeeper implements the ProviderApp interface.
func (app *App) GetTestSlashingKeeper() icstest.TestSlashingKeeper { //nolint:nolintlint
	return app.SlashingKeeper
}

// GetTestDistributionKeeper implements the ProviderApp interface.
func (app *App) GetTestDistributionKeeper() icstest.TestDistributionKeeper { //nolint:nolintlint
	return app.DistrKeeper
}

func (app *App) GetTestAccountKeeper() icstest.TestAccountKeeper { //nolint:nolintlint
	return app.AccountKeeper
}
