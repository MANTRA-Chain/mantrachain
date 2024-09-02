package app

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	v2 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v2"
	v3 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v3"
	v4 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v4"
	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	guardtypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	lpfarmtypes "github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	marketmakertypes "github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	tokentypes "github.com/MANTRA-Finance/mantrachain/x/token/types"
	txfeestypes "github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func (app *App) RegisterUpgradeHandlers() {
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		var keyTable paramstypes.KeyTable
		var customModule bool
		switch subspace.Name() {
		case coinfactorytypes.ModuleName:
			keyTable = coinfactorytypes.ParamKeyTable()
			customModule = true
		case guardtypes.ModuleName:
			keyTable = guardtypes.ParamKeyTable()
			customModule = true
		case liquiditytypes.ModuleName:
			keyTable = liquiditytypes.ParamKeyTable()
			customModule = true
		case lpfarmtypes.ModuleName:
			keyTable = lpfarmtypes.ParamKeyTable()
			customModule = true
		case marketmakertypes.ModuleName:
			keyTable = marketmakertypes.ParamKeyTable()
			customModule = true
		case tokentypes.ModuleName:
			keyTable = tokentypes.ParamKeyTable()
			customModule = true
		case txfeestypes.ModuleName:
			keyTable = txfeestypes.ParamKeyTable()
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
		v4.CreateUpgradeHandler(app.ModuleManager, app.Configurator(), *app.GovKeeper, *app.GuardKeeper, app.IBCKeeper.ChannelKeeper, app.ConsensusParamsKeeper, *app.MarketMapKeeper),
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
	//	var upgradeHandler *upgradetypes.UpgradeHandler

	switch upgradeInfo.Name {
	case v2.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added:   []string{"bridge", "circuit"},
			Deleted: []string{"farming", "rewards"},
		}
	case v4.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				marketmaptypes.ModuleName,
				oracletypes.ModuleName,
			},
			Deleted: []string{"icahost", "icacontroller"},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}

	// previous code was:
	// if upgradeHandler != nil {
	//     app.UpgradeKeeper.SetUpgradeHandler(upgradeInfo.Name, *upgradeHandler)
	// }
	// but there is no condition where the upgradehandler is not nil

	// app.UpgradeKeeper.SetUpgradeHandler(upgradeInfo.Name, *upgradeHandler)
}
