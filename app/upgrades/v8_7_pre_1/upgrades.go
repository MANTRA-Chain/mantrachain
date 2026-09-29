package v8_7_pre_1

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
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
		ctx.Logger().Info("Starting v8.7.0-pre.1 upgrade...")

		// Before module migrations, which may decode ContractInfo.
		ctx.Logger().Info("Migrating ContractInfo encoding...")
		migrated, err := upgrades.MigrateContractInfo(ctx, storekeys[wasmtypes.StoreKey])
		if err != nil {
			return vm, err
		}
		ctx.Logger().Info("Migrated ContractInfo records", "count", migrated)

		ctx.Logger().Info("Running module migrations...")
		vm, err = mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v8.7.0-pre.1 complete")
		return vm, nil
	}
}
