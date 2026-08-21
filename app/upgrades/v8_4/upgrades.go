package v8_4

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
	storekeys map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting v8.4.0 upgrade...")

		ctx.Logger().Info("Running module migrations...")
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		AccountsToBlacklist := []string{"mantra13n9sk3p8x7tpq9adgxvzv9q0qev953mld0hwva"}
		if err := keepers.SanctionKeeper.AddToBlacklist(ctx, AccountsToBlacklist); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v8.4.0 complete")
		return vm, nil
	}
}
