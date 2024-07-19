package v2

import (
	"encoding/binary"
	"fmt"

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
	migrateSnapshots(ctx, store, cdc)

	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	return store.Set(types.ParamsKey, bz)
}

func migrateSnapshots(
	ctx sdk.Context,
	store store.KVStore,
	cdc codec.BinaryCodec,
) {
	storeAdapter := runtime.KVStoreAdapter(store)
	snapshotsStore := prefix.NewStore(storeAdapter, types.SnapshotStoreKey(0))
	iterator := storetypes.KVStorePrefixIterator(snapshotsStore, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var snapshot v1types.Snapshot
		cdc.MustUnmarshal(iterator.Value(), &snapshot)

		snapshotsStoreByPairId := prefix.NewStore(storeAdapter, types.SnapshotStoreKey(snapshot.PairId))

		lastSnapshot, found := getLastSnapshot(ctx, store, cdc, snapshot.PairId)

		if !found {
			panic(fmt.Sprintf("last snapshot not found %d", snapshot.PairId))
		}

		pools := make([]*types.SnapshotPool, len(snapshot.Pools))
		for i, pool := range snapshot.Pools {
			coinSupply := math.LegacyNewDec(0)

			// If the last snapshot is the same as the current one, then we need to calculate the coin supply
			// for the pool, otwerwise we can set it to 0. It's not ideal but it's the best we can do.
			if lastSnapshot.Id == snapshot.Id {
				coinSupply = getPoolCoinSupply(ctx, store, cdc, snapshot.PairId, pool.PoolId)
			}

			newPool := types.SnapshotPool{
				PoolId:          pool.PoolId,
				CoinSupply:      coinSupply,
				RewardsPerToken: pool.RewardsPerToken,
			}

			pools[i] = &newPool
		}

		upgradedSnapshot := types.Snapshot{
			Id:            snapshot.Id,
			PairId:        snapshot.PairId,
			Pools:         pools,
			PoolIdToIdx:   snapshot.PoolIdToIdx,
			Distributed:   snapshot.Distributed,
			DistributedAt: snapshot.DistributedAt,
			Remaining:     snapshot.Remaining,
		}

		b := cdc.MustMarshal(&upgradedSnapshot)
		snapshotsStoreByPairId.Set(GetSnapshotIDBytes(upgradedSnapshot.Id), b)
	}
}

func getLastSnapshot(
	ctx sdk.Context,
	store store.KVStore,
	cdc codec.BinaryCodec,
	pairId uint64,
) (val v1types.Snapshot, found bool) {
	storeAdapter := runtime.KVStoreAdapter(store)
	snapshotsStore := prefix.NewStore(storeAdapter, types.SnapshotStoreKey(pairId))
	snapshotCount, found := getSnapshotCount(ctx, store, cdc, pairId)

	if !found {
		return val, false
	}

	b := snapshotsStore.Get(GetSnapshotIDBytes(snapshotCount.Count - 1))
	if b == nil {
		return val, false
	}
	cdc.MustUnmarshal(b, &val)
	return val, true
}

func getSnapshotCount(
	_ sdk.Context,
	store store.KVStore,
	cdc codec.BinaryCodec,
	pairId uint64,
) (val v1types.SnapshotCount, found bool) {
	storeAdapter := runtime.KVStoreAdapter(store)
	snapshotsStore := prefix.NewStore(storeAdapter, types.SnapshotCountStoreKey(pairId))
	byteKey := types.KeyPrefix(types.SnapshotCountKey)
	bz := snapshotsStore.Get(byteKey)

	if bz == nil {
		return val, false
	}

	cdc.MustUnmarshal(bz, &val)
	return val, true
}

func getPoolCoinSupply(
	_ sdk.Context,
	store store.KVStore,
	cdc codec.BinaryCodec,
	pairId uint64,
	poolId uint64,
) math.LegacyDec {
	storeAdapter := runtime.KVStoreAdapter(store)
	providerStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProviderKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(providerStore, []byte{})
	coinSupply := math.LegacyNewDec(0)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var provider v1types.Provider
		cdc.MustUnmarshal(iterator.Value(), &provider)

		pairIdx, found := provider.PairIdToIdx[pairId]
		if !found {
			continue
		}

		providerPair := provider.Pairs[pairIdx]

		poolIdx, found := providerPair.PoolIdToIdx[poolId]
		if !found {
			continue
		}

		balance := providerPair.Balances[poolIdx]

		if balance.IsPositive() {
			coinSupply = coinSupply.Add(balance.Amount.ToLegacyDec())
		}
	}

	return coinSupply

}
