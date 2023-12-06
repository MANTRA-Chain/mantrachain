package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
)

// GetSnapshotCount get the total number of snapshot
func (k Keeper) GetSnapshotCount(ctx sdk.Context) uint64 {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.SnapshotCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetSnapshotCount set the total number of snapshot
func (k Keeper) SetSnapshotCount(ctx sdk.Context, count uint64)  {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.SnapshotCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendSnapshot appends a snapshot in the store with a new id and update the count
func (k Keeper) AppendSnapshot(
    ctx sdk.Context,
    snapshot types.Snapshot,
) uint64 {
	// Create the snapshot
    count := k.GetSnapshotCount(ctx)

    // Set the ID of the appended value
    snapshot.Id = count

    store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SnapshotKey))
    appendedValue := k.cdc.MustMarshal(&snapshot)
    store.Set(GetSnapshotIDBytes(snapshot.Id), appendedValue)

    // Update snapshot count
    k.SetSnapshotCount(ctx, count+1)

    return count
}

// SetSnapshot set a specific snapshot in the store
func (k Keeper) SetSnapshot(ctx sdk.Context, snapshot types.Snapshot) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SnapshotKey))
	b := k.cdc.MustMarshal(&snapshot)
	store.Set(GetSnapshotIDBytes(snapshot.Id), b)
}

// GetSnapshot returns a snapshot from its id
func (k Keeper) GetSnapshot(ctx sdk.Context, id uint64) (val types.Snapshot, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SnapshotKey))
	b := store.Get(GetSnapshotIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSnapshot removes a snapshot from the store
func (k Keeper) RemoveSnapshot(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SnapshotKey))
	store.Delete(GetSnapshotIDBytes(id))
}

// GetAllSnapshot returns all snapshot
func (k Keeper) GetAllSnapshot(ctx sdk.Context) (list []types.Snapshot) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SnapshotKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Snapshot
		k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
	}

    return
}

// GetSnapshotIDBytes returns the byte representation of the ID
func GetSnapshotIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetSnapshotIDFromBytes returns ID in uint64 format from a byte array
func GetSnapshotIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
