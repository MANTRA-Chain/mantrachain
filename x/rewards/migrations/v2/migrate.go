package v2

import (
	"encoding/binary"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/exported"
	v1types "github.com/MANTRA-Finance/mantrachain/x/rewards/migrations/v1/types"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetSnapshotIDBytes returns the byte representation of the ID
func GetSnapshotIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

func MigrateStore(
	ctx sdk.Context,
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	legacySubspace exported.Subspace,
) error {
	store := storeService.OpenKVStore(ctx)
	migrateSnapshots(store, cdc)

	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	return store.Set(types.ParamsKey, bz)
}

func migrateSnapshots(
	store store.KVStore,
	cdc codec.BinaryCodec,
) {
	storeAdapter := runtime.KVStoreAdapter(store)
	snapshotsStore := prefix.NewStore(storeAdapter, types.SnapshotStoreKey(0))
	iterator := storetypes.KVStorePrefixIterator(snapshotsStore, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val v1types.Snapshot
		cdc.MustUnmarshal(iterator.Value(), &val)

		pools := make([]*types.SnapshotPool, len(val.Pools))
		for i, pool := range val.Pools {

			newPool := types.SnapshotPool{
				PoolId:          pool.PoolId,
				CoinSupply:      math.LegacyNewDec(0), // Not ok to use zero value, but we don't have the historical data to populate this field
				RewardsPerToken: pool.RewardsPerToken,
			}

			pools[i] = &newPool
		}

		upgradedSnapshot := types.Snapshot{
			Id:            val.Id,
			PairId:        val.PairId,
			Pools:         pools,
			PoolIdToIdx:   val.PoolIdToIdx,
			Distributed:   val.Distributed,
			DistributedAt: val.DistributedAt,
			Remaining:     val.Remaining,
		}

		b := cdc.MustMarshal(&upgradedSnapshot)

		snapshotsStoreByPairId := prefix.NewStore(storeAdapter, types.SnapshotStoreKey(upgradedSnapshot.PairId))
		snapshotsStoreByPairId.Set(GetSnapshotIDBytes(upgradedSnapshot.Id), b)
	}
}
