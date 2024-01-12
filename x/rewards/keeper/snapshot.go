package keeper

import (
	"encoding/binary"

	"github.com/AumegaChain/aumega/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetSnapshotStartId(ctx sdk.Context, pairId uint64) (val types.SnapshotStartId, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStartIdStoreKey(pairId))
	byteKey := types.KeyPrefix(types.SnapshotStartIdKey)
	bz := store.Get(byteKey)

	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

func (k Keeper) SetSnapshotStartId(ctx sdk.Context, snapshotStartId types.SnapshotStartId) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStartIdStoreKey(snapshotStartId.PairId))
	byteKey := types.KeyPrefix(types.SnapshotStartIdKey)
	b := k.cdc.MustMarshal(&snapshotStartId)
	store.Set(byteKey, b)
}

func (k Keeper) GetAllSnapshotStartId(ctx sdk.Context) (list []types.SnapshotStartId) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStartIdStoreKey(0))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SnapshotStartId
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetSnapshotsLastDistributedAt(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.SnapshotsLastDistributedAtKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetSnapshotsLastDistributedAt(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.SnapshotsLastDistributedAtKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

func (k Keeper) GetSnapshotCount(ctx sdk.Context, pairId uint64) (val types.SnapshotCount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotCountStoreKey(pairId))
	byteKey := types.KeyPrefix(types.SnapshotCountKey)
	bz := store.Get(byteKey)

	if bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

func (k Keeper) SetSnapshotCount(ctx sdk.Context, snapshotCount types.SnapshotCount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotCountStoreKey(snapshotCount.PairId))
	byteKey := types.KeyPrefix(types.SnapshotCountKey)
	b := k.cdc.MustMarshal(&snapshotCount)
	store.Set(byteKey, b)
}

func (k Keeper) GetAllSnapshotCount(ctx sdk.Context) (list []types.SnapshotCount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotCountStoreKey(0))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SnapshotCount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// AppendSnapshot appends a snapshot in the store with a new id and update the count
func (k Keeper) AppendSnapshot(
	ctx sdk.Context,
	snapshot types.Snapshot,
) uint64 {
	// Create the snapshot
	snapshotCount, found := k.GetSnapshotCount(ctx, snapshot.PairId)

	if !found {
		snapshotCount = types.SnapshotCount{PairId: snapshot.PairId, Count: 0}
	}

	// Set the ID of the appended value
	snapshot.Id = snapshotCount.Count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStoreKey(snapshot.PairId))
	appendedValue := k.cdc.MustMarshal(&snapshot)
	store.Set(GetSnapshotIDBytes(snapshot.Id), appendedValue)

	// Update snapshot count
	snapshotCount.Count += 1
	k.SetSnapshotCount(ctx, snapshotCount)

	return snapshot.Id
}

// SetSnapshot set a specific snapshot in the store
func (k Keeper) SetSnapshot(ctx sdk.Context, snapshot types.Snapshot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStoreKey(snapshot.PairId))
	b := k.cdc.MustMarshal(&snapshot)
	store.Set(GetSnapshotIDBytes(snapshot.Id), b)
}

// GetSnapshot returns a snapshot from its id
func (k Keeper) GetSnapshot(ctx sdk.Context, pairId uint64, id uint64) (val types.Snapshot, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStoreKey(pairId))
	b := store.Get(GetSnapshotIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetLastSnapshot(ctx sdk.Context, pairId uint64) (val types.Snapshot, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStoreKey(pairId))
	snapshotCount, found := k.GetSnapshotCount(ctx, pairId)

	if !found {
		return val, false
	}

	b := store.Get(GetSnapshotIDBytes(snapshotCount.Count - 1))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSnapshot removes a snapshot from the store
func (k Keeper) RemoveSnapshot(ctx sdk.Context, pairId uint64, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStoreKey(pairId))
	store.Delete(GetSnapshotIDBytes(id))
}

// GetAllSnapshot returns all snapshot
func (k Keeper) GetAllSnapshot(ctx sdk.Context, pairId uint64) (list []types.Snapshot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStoreKey(pairId))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Snapshot
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetSnapshotsInRange(ctx sdk.Context, pairId uint64, startId uint64, endId uint64) (list []types.Snapshot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SnapshotStoreKey(pairId))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Snapshot
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.Id >= startId && val.Id <= endId {
			list = append(list, val)
		}
	}

	return
}

func (k Keeper) GetEndClaimedSnapshotId(ctx sdk.Context, pairId uint64) (endClaimedSnapshotId uint64) {
	lastSnapshot, found := k.GetLastSnapshot(ctx, pairId)

	if found {
		if lastSnapshot.Distributed {
			endClaimedSnapshotId = lastSnapshot.Id
		} else if lastSnapshot.Id > 0 { // if this snapshot is not distributed, then the last claimed snapshot is the previous one
			endClaimedSnapshotId = lastSnapshot.Id - 1
		}
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
