package app

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	farmingtypes "github.com/MANTRA-Finance/mantrachain/x/farming/types"
	guardtypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	lpfarmtypes "github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	marketmakertypes "github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	rewardstypes "github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	tokentypes "github.com/MANTRA-Finance/mantrachain/x/token/types"
	txfeestypes "github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	v2 "github.com/MANTRA-Finance/mantrachain/app/upgrades/v2"
)

func (app *App) RegisterUpgradeHandlers() {
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		var keyTable paramstypes.KeyTable
		var customModule bool
		switch subspace.Name() {
		case coinfactorytypes.ModuleName:
			keyTable = coinfactorytypes.ParamKeyTable() //nolint:staticcheck
			customModule = true
		case farmingtypes.ModuleName:
			keyTable = farmingtypes.ParamKeyTable() //nolint:staticcheck
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
		case rewardstypes.ModuleName:
			keyTable = rewardstypes.ParamKeyTable() //nolint:staticcheck
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
			Added: []string{"bridge", "circuit"},
		}
	default:
		// no-op
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
