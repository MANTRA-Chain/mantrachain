package v6rc0

import (
	"context"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v6/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	erc20types "github.com/cosmos/evm/x/erc20/types"
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

		// update contract owner for all existing tokenfactory token_pairs
		pairs := keepers.Erc20Keeper.GetTokenPairs(ctx)
		for _, pair := range pairs {
			if strings.HasPrefix(pair.Denom, "factory/") {
				pair.ContractOwner = erc20types.OWNER_MODULE
				keepers.Erc20Keeper.SetTokenPair(ctx, pair)
			}
		}

		ctx.Logger().Info("Upgrade v6.0.0-rc0 complete")
		return vm, nil
	}
}
