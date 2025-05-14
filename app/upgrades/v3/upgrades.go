package v3

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const BLACKLIST_ACCOUNT = "mantra14t56rzvxzw0yp9plcf9dy6rr53chyvxt4cqtt5"

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		err = keepers.SanctionKeeper.BlacklistAccounts.Set(ctx, BLACKLIST_ACCOUNT)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v3 complete")
		return vm, nil
	}
}
