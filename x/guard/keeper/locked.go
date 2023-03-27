package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetLocked set a specific locked in the store from its index
func (k Keeper) SetLocked(
	ctx sdk.Context,
	index []byte,
	kind types.LockedKind,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedStoreKey(kind.Bytes()))
	store.Set(index, types.Placeholder)
}

func (k Keeper) HasLocked(
	ctx sdk.Context,
	index []byte,
	kind types.LockedKind,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedStoreKey(kind.Bytes()))
	return store.Has(index)
}

// GetLocked returns a locked from its index
func (k Keeper) GetLocked(
	ctx sdk.Context,
	index []byte,
	kind types.LockedKind,
) (val []byte, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedStoreKey(kind.Bytes()))

	if !k.HasLocked(ctx, index, kind) {
		return []byte{}, false
	}

	b := store.Get(index)

	if b == nil {
		return val, false
	}

	return b, true
}

// RemoveLocked removes a locked from the store
func (k Keeper) RemoveLocked(
	ctx sdk.Context,
	index []byte,
	kind types.LockedKind,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedStoreKey(kind.Bytes()))
	store.Delete(index)
}

// GetAllLocked returns all locked
func (k Keeper) GetAllLocked(ctx sdk.Context, kind types.LockedKind) (list []*types.Locked) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedStoreKey(kind.Bytes()))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, &types.Locked{
			Index:  iterator.Key(),
			Locked: iterator.Value(),
			Kind:   kind.String(),
		})
	}

	return
}
