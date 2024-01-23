package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetNftCollection(ctx sdk.Context, nftCollection types.NftCollection) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftCollectionStoreKey(nftCollection.Creator))
	b := k.cdc.MustMarshal(&nftCollection)
	store.Set(nftCollection.Index, b)
}

func (k Keeper) GetNftCollection(
	ctx sdk.Context,
	creator sdk.AccAddress,
	index []byte,
) (val types.NftCollection, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftCollectionStoreKey(creator))

	if !k.HasNftCollection(ctx, creator, index) {
		return types.NftCollection{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasNftCollection(
	ctx sdk.Context,
	creator sdk.AccAddress,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftCollectionStoreKey(creator))
	return store.Has(index)
}

func (k Keeper) GetAllNftCollection(ctx sdk.Context) (list []types.NftCollection) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftCollectionStoreKey(nil))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NftCollection
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetNftCollectionOwner(ctx sdk.Context, index []byte, owner sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftCollectionOwnerKeyPrefix))
	store.Set(index, owner)
}

func (k Keeper) GetNftCollectionOwner(
	ctx sdk.Context,
	index []byte,
) (val []byte, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftCollectionOwnerKeyPrefix))

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	return b, true
}

func (k Keeper) GetAllNftCollectionOwner(ctx sdk.Context) (list []*types.NftCollectionOwner) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NftCollectionOwnerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, &types.NftCollectionOwner{
			Index: iterator.Key(),
			Owner: iterator.Value(),
		})
	}

	return
}

func (k Keeper) SetOpenedNftsCollection(ctx sdk.Context, index []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OpenedNftsCollectionKeyPrefix))
	store.Set(index, types.Placeholder)
}

func (k Keeper) HasOpenedNftsCollection(
	ctx sdk.Context,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OpenedNftsCollectionKeyPrefix))
	return store.Has(index)
}

func (k Keeper) GetAllOpenedNftsCollection(ctx sdk.Context) (list [][]byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OpenedNftsCollectionKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Key())
	}

	return
}

func (k Keeper) SetSoulBondedNftsCollection(ctx sdk.Context, index []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	store.Set(index, types.Placeholder)
}

func (k Keeper) HasSoulBondedNftsCollection(
	ctx sdk.Context,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	return store.Has(index)
}

func (k Keeper) GetAllSoulBondedNftsCollection(ctx sdk.Context) (list [][]byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Key())
	}

	return
}

func (k Keeper) SetRestrictedNftsCollection(ctx sdk.Context, index []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RestrictedNftsCollectionKeyPrefix))
	store.Set(index, types.Placeholder)
}

func (k Keeper) HasRestrictedNftsCollection(
	ctx sdk.Context,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RestrictedNftsCollectionKeyPrefix))
	return store.Has(index)
}

func (k Keeper) GetAllRestrictedNftsCollection(ctx sdk.Context) (list [][]byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RestrictedNftsCollectionKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Key())
	}

	return
}
