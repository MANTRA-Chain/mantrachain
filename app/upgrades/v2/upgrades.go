package v2

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	v1_rewardstypes "github.com/MANTRA-Finance/mantrachain/app/upgrades/v2/v1_types/rewards"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	rewardskeeper "github.com/MANTRA-Finance/mantrachain/x/rewards/keeper"
	rewardstypes "github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var logger log.Logger

// MustGetStoreKey is a helper to directly access the KV-Store
func MustGetStoreKey(storeKeys []storetypes.StoreKey, storeName string) storetypes.StoreKey {
	for _, k := range storeKeys {
		if k.Name() == storeName {
			return k
		}
	}

	panic(fmt.Sprintf("failed to find store key: %s", storeName))
}

// NOTE: we're only keeping this logic for the upgrade tests
// This is not the original upgrade logic.
// Look into the previous version if want to know what the upgrade logic was
// CreateUpgradeHandler creates an SDK upgrade handler for v15.0.0
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.Codec,
	storeKeys []storetypes.StoreKey,
	rewardsKeeper rewardskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger = sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		migrateRewardsModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, rewardstypes.StoreKey), rewardsKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func migrateRewardsModule(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	rewardsStoreKey storetypes.StoreKey,
	rewardsKeeper rewardskeeper.Keeper,
) error {
	store := prefix.NewStore(ctx.KVStore(rewardsStoreKey), rewardstypes.SnapshotStoreKey(0))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	coinSupply := math.LegacyMustNewDecFromStr("0")

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val v1_rewardstypes.Snapshot
		cdc.MustUnmarshal(iterator.Value(), &val)

		pools := make([]*rewardstypes.SnapshotPool, len(val.Pools))
		for i, pool := range val.Pools {

			newPool := rewardstypes.SnapshotPool{
				PoolId:          pool.PoolId,
				CoinSupply:      coinSupply, // Not ok to use zero value, but we don't have the historical data to populate this field
				RewardsPerToken: pool.RewardsPerToken,
			}

			pools[i] = &newPool
		}

		upgradedSnapshot := rewardstypes.Snapshot{
			Id:            val.Id,
			PairId:        val.PairId,
			Pools:         pools,
			PoolIdToIdx:   val.PoolIdToIdx,
			Distributed:   val.Distributed,
			DistributedAt: val.DistributedAt,
			Remaining:     val.Remaining,
		}

		rewardsKeeper.SetSnapshot(ctx, upgradedSnapshot)
	}

	return nil
}
