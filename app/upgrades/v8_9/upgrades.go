package v8_9

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// CreateUpgradeHandler removes the interchain accounts module. Its stores are
// deleted by the store loader (see Upgrade.StoreUpgrades); the handler removes
// the state it left in other modules' stores, then runs module migrations.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *upgrades.UpgradeKeepers,
	storekeys map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting v8.9.0 upgrade...")

		ctx.Logger().Info("Removing interchain accounts module state...")
		removed := RemoveICAState(ctx, storekeys, vm)
		ctx.Logger().Info("Removed interchain accounts module state", "legacy_params", removed)

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v8.9.0 complete")
		return vm, nil
	}
}

// RemoveICAState deletes what the removed interchain accounts module left
// outside its own stores and returns the number of legacy params deleted:
//
//   - its entry in the x/upgrade module version map. RunMigrations only
//     handles modules in the module manager, so the entry would otherwise be
//     kept and reported by the module_versions query forever.
//   - its legacy x/params subspaces, which the app no longer registers.
func RemoveICAState(ctx sdk.Context, storekeys map[string]*storetypes.KVStoreKey, vm module.VersionMap) int {
	delete(vm, ICAModuleName)
	versionStore := prefix.NewStore(ctx.KVStore(storekeys[upgradetypes.StoreKey]), []byte{upgradetypes.VersionMapByte})
	versionStore.Delete([]byte(ICAModuleName))

	removed := 0
	paramsStore := ctx.KVStore(storekeys[paramstypes.StoreKey])
	for _, subspace := range []string{ICAControllerStoreKey, ICAHostStoreKey} {
		subspaceStore := prefix.NewStore(paramsStore, []byte(subspace+"/"))

		var keys [][]byte
		iterator := subspaceStore.Iterator(nil, nil)
		for ; iterator.Valid(); iterator.Next() {
			keys = append(keys, iterator.Key())
		}
		iterator.Close()

		for _, key := range keys {
			subspaceStore.Delete(key)
		}
		removed += len(keys)
	}

	return removed
}
