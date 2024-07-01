package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetNftCollection(ctx sdk.Context, nftCollection types.NftCollection) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftCollectionStoreKey(nftCollection.Creator))
	b := k.cdc.MustMarshal(&nftCollection)
	store.Set(nftCollection.Index, b)
}

func (k Keeper) GetNftCollection(
	ctx sdk.Context,
	creator sdk.AccAddress,
	index []byte,
) (val types.NftCollection, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftCollectionStoreKey(creator))

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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftCollectionStoreKey(creator))
	return store.Has(index)
}

func (k Keeper) GetAllNftCollection(ctx sdk.Context) (list []types.NftCollection) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftCollectionStoreKey(nil))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NftCollection
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetNftCollectionOwner(ctx sdk.Context, index []byte, owner sdk.AccAddress) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.NftCollectionOwnerKeyPrefix))
	store.Set(index, owner)
}

func (k Keeper) GetNftCollectionOwner(
	ctx sdk.Context,
	index []byte,
) (val []byte, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.NftCollectionOwnerKeyPrefix))

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	return b, true
}

func (k Keeper) GetAllNftCollectionOwner(ctx sdk.Context) (list []*types.NftCollectionOwner) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.NftCollectionOwnerKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OpenedNftsCollectionKeyPrefix))
	store.Set(index, types.Placeholder)
}

func (k Keeper) HasOpenedNftsCollection(
	ctx sdk.Context,
	index []byte,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OpenedNftsCollectionKeyPrefix))
	return store.Has(index)
}

func (k Keeper) GetAllOpenedNftsCollection(ctx sdk.Context) (list [][]byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OpenedNftsCollectionKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Key())
	}

	return
}

func (k Keeper) SetSoulBondedNftsCollection(ctx sdk.Context, index []byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	store.Set(index, types.Placeholder)
}

func (k Keeper) HasSoulBondedNftsCollection(
	ctx sdk.Context,
	index []byte,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	return store.Has(index)
}

func (k Keeper) GetAllSoulBondedNftsCollection(ctx sdk.Context) (list [][]byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Key())
	}

	return
}

func (k Keeper) SetRestrictedNftsCollection(ctx sdk.Context, index []byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.RestrictedNftsCollectionKeyPrefix))
	store.Set(index, types.Placeholder)
}

func (k Keeper) HasRestrictedNftsCollection(
	ctx sdk.Context,
	index []byte,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.RestrictedNftsCollectionKeyPrefix))
	return store.Has(index)
}

func (k Keeper) GetAllRestrictedNftsCollection(ctx sdk.Context) (list [][]byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.RestrictedNftsCollectionKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Key())
	}

	return
}
