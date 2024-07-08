package v2

import (
	"encoding/binary"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	v2types "github.com/MANTRA-Finance/mantrachain/x/rewards/migrations/v2/types"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(
	ctx sdk.Context,
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
) error {
	storeAdapter := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.SnapshotStoreKey(0))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	coinSupply := math.LegacyMustNewDecFromStr("0")

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val v2types.Snapshot
		cdc.MustUnmarshal(iterator.Value(), &val)

		pools := make([]*types.SnapshotPool, len(val.Pools))
		for i, pool := range val.Pools {

			newPool := types.SnapshotPool{
				PoolId:          pool.PoolId,
				CoinSupply:      coinSupply, // Not ok to use zero value, but we don't have the historical data to populate this field
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

		bz := make([]byte, 8)
		binary.BigEndian.PutUint64(bz, upgradedSnapshot.Id)

		store.Set(bz, b)
	}

	return nil
}
