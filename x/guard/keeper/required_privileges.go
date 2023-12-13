package keeper

import (
	"github.com/MANTRA-Finance/aumega/x/guard/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetRequiredPrivileges(
	ctx sdk.Context,
	index []byte,
	kind types.RequiredPrivilegesKind,
	privileges []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequiredPrivilegesStoreKey(kind.Bytes()))
	store.Set(index, privileges)
}

func (k Keeper) HasRequiredPrivileges(
	ctx sdk.Context,
	index []byte,
	kind types.RequiredPrivilegesKind,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequiredPrivilegesStoreKey(kind.Bytes()))
	return store.Has(index)
}

func (k Keeper) GetRequiredPrivileges(
	ctx sdk.Context,
	index []byte,
	kind types.RequiredPrivilegesKind,
) (val []byte, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequiredPrivilegesStoreKey(kind.Bytes()))

	if !k.HasRequiredPrivileges(ctx, index, kind) {
		return []byte{}, false
	}

	b := store.Get(index)

	if b == nil {
		return val, false
	}

	return b, true
}

func (k Keeper) GetRequiredPrivilegesMany(
	ctx sdk.Context,
	indexes [][]byte,
	kind types.RequiredPrivilegesKind,
) (list [][]byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequiredPrivilegesStoreKey(kind.Bytes()))

	for _, index := range indexes {
		bz := store.Get(index)

		list = append(list, bz)
	}

	return
}

func (k Keeper) RemoveRequiredPrivileges(
	ctx sdk.Context,
	index []byte,
	kind types.RequiredPrivilegesKind,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequiredPrivilegesStoreKey(kind.Bytes()))
	store.Delete(index)
}

func (k Keeper) GetAllRequiredPrivileges(ctx sdk.Context, kind types.RequiredPrivilegesKind) (list []*types.RequiredPrivileges) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RequiredPrivilegesStoreKey(kind.Bytes()))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, &types.RequiredPrivileges{
			Index:      iterator.Key(),
			Privileges: iterator.Value(),
			Kind:       kind.String(),
		})
	}

	return
}
