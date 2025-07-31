package v5rc3

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
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
		ctx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		iter := keepers.TokenFactoryKeeper.GetAllDenomsIterator(ctx)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			denom := string(iter.Value())
			if err := keepers.TokenFactoryKeeper.UpdateDenomWithERC20(ctx, denom); err != nil {
				return vm, err
			}
		}

		// enable AllowUnprotectedTxs see adr-006
		params := keepers.EVMKeeper.GetParams(ctx)
		params.AllowUnprotectedTxs = true
		keepers.EVMKeeper.SetParams(ctx, params)

		ctx.Logger().Info("Upgrade v5.0.0-rc3 complete")
		return vm, nil
	}
}
