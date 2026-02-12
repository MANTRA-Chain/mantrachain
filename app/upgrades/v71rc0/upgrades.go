package v71rc0

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

		iter, err := keepers.SanctionKeeper.BlacklistAccounts.Iterate(ctx, nil)
		if err != nil {
			ctx.Logger().Error("Failed to iterate blacklisted addresses", "err", err)
			return vm, nil
		}
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			addrStr, err := iter.Key()
			if err != nil {
				continue
			}

			delegatorAddr, err := sdk.AccAddressFromBech32(addrStr)
			if err != nil {
				continue
			}

			delegations, err := keepers.StakingKeeper.GetDelegatorDelegations(ctx, delegatorAddr, 1000)
			if err != nil || len(delegations) == 0 {
				continue
			}

			ctx.Logger().Info("Undelegating from blacklisted address", "address", addrStr, "delegations", len(delegations))
			for _, delegation := range delegations {
				valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
				if err != nil {
					continue
				}

				if _, _, err = keepers.StakingKeeper.Undelegate(ctx, delegatorAddr, valAddr, delegation.Shares); err != nil {
					ctx.Logger().Error("Failed to undelegate", "delegator", addrStr, "validator", delegation.ValidatorAddress, "err", err)
				}
			}
		}

		ctx.Logger().Info("Upgrade complete", "upgrade", UpgradeName)
		return vm, nil
	}
}
