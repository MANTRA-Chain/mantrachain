package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	guardtypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	lpfarmtypes "github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	marketmakertypes "github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	tokentypes "github.com/MANTRA-Finance/mantrachain/x/token/types"
	txfeestypes "github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	v2 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v2"
	v3 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v3"
	v4 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v4"
)

func (app *App) RegisterUpgradeHandlers() {
	var defaultUpgradeHandler = func(goCtx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return app.ModuleManager.RunMigrations(goCtx, app.Configurator(), fromVM)
	}

	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		var keyTable paramstypes.KeyTable
		var customModule bool
		switch subspace.Name() {
		case coinfactorytypes.ModuleName:
			keyTable = coinfactorytypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		case guardtypes.ModuleName:
			keyTable = guardtypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		case liquiditytypes.ModuleName:
			keyTable = liquiditytypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		case lpfarmtypes.ModuleName:
			keyTable = lpfarmtypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		case marketmakertypes.ModuleName:
			keyTable = marketmakertypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		case tokentypes.ModuleName:
			keyTable = tokentypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		case txfeestypes.ModuleName:
			keyTable = txfeestypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		}

		if customModule && !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}

	upgradeKey := app.GetKey("upgrade")
	authority := app.AccountKeeper.GetAuthority()

	// v2 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v2.UpgradeName,
		v2.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.appCodec,
			&app.ConsensusParamsKeeper.ParamsStore,
			upgradeKey,
			authority,
		),
	)

	// v3 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v3.UpgradeName,
		v3.CreateUpgradeHandler(app.ModuleManager, app.Configurator(), *app.GovKeeper),
	)

	// v4 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v4.UpgradeName,
		defaultUpgradeHandler,
	)

	// When a planned update height is reached, the old binary will panic
	// writing on disk the height and name of the update that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades
	var upgradeHandler *upgradetypes.UpgradeHandler

	// Slinky upgrader
	slinkyUpgrade := createSlinkyUpgrader(app)

	switch upgradeInfo.Name {
	case v2.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added:   []string{"bridge", "circuit"},
			Deleted: []string{"farming", "rewards"},
		}
	case v4.UpgradeName:
		upgradeHandler = &slinkyUpgrade.Handler
		storeUpgrades = &slinkyUpgrade.StoreUpgrade
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}

	if upgradeHandler != nil {
		app.UpgradeKeeper.SetUpgradeHandler(upgradeInfo.Name, *upgradeHandler)
	}
}
