package v7rc3

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v7/app/upgrades"
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

		if keepers.OracleKeeper != nil {
			pairs := keepers.OracleKeeper.GetAllCurrencyPairs(ctx)
			for _, pair := range pairs {
				if pair.Base == "OM" || pair.Quote == "OM" {
					if err := keepers.OracleKeeper.RemoveCurrencyPair(ctx, pair); err != nil {
						ctx.Logger().Error("Failed to remove OM currency pair", "pair", pair.String(), "err", err)
					}
				}
			}
		}

		ctx.Logger().Info("Migrating x/distribution state...")
		// Unblock all module accounts for the duration of the migration.
		// This is temporary and only applies to the bank keeper instance used within this upgrade handler.
		keepers.BankKeeper = keepers.BankKeeper.WithBlockedAddrs(nil)
		if err = migrateDistr(ctx, keepers.DistrKeeper, keepers.AccountKeeper, keepers.BankKeeper); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete", "name", UpgradeName)
		return vm, nil
	}
}
