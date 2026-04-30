package v8_1

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
		ctx.Logger().Info("Starting v8.1.0 upgrade...")

		// Repair before RunMigrations so reward-touching migrations don't panic.
		ctx.Logger().Info("Clamping silently-skipped slash residue on delegator starting info...")
		if err := fixSilentlySkippedSlashes(ctx, keepers.StakingKeeper, keepers.DistrKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Resetting starting info on silent-slash validators...")
		records, err := resolveSilentSlashes(plan)
		if err != nil {
			return vm, err
		}
		if err := repairSilentSlashes(ctx, keepers.StakingKeeper, keepers.DistrKeeper, records); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Running module migrations...")
		vm, err = mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v8.1.0 complete")
		return vm, nil
	}
}
