package v8_9

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// CreateUpgradeHandler runs module migrations. x/sanction moves from
// consensus version 1 to 2: params get their defaults, existing blacklist
// entries are kept but lose any EIP-7702 delegation, and the sanctioner list
// starts empty.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *upgrades.UpgradeKeepers,
	_ map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting v8.9.0 upgrade...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v8.9.0 complete")
		return vm, nil
	}
}
