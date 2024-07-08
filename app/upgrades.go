package app

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	rewardskeeper "github.com/MANTRA-Finance/mantrachain/x/rewards/keeper"

	v2 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v2"
)

func (app *App) setupUpgradeHandlers(rewardsKeeper rewardskeeper.Keeper) {
	// v2 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v2.UpgradeName,
		v2.CreateUpgradeHandler(
			app.ModuleManager, app.Configurator(), app.appCodec, app.GetStoreKeys(), rewardsKeeper,
		),
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

	switch upgradeInfo.Name {
	case v2.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{"bridge"},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
