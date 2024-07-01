package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetProvider set a specific provider in the store from its index
func (k Keeper) SetProvider(ctx sdk.Context, provider types.Provider) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProviderKeyPrefix))
	b := k.cdc.MustMarshal(&provider)
	store.Set(types.ProviderKey(
		provider.Index,
	), b)
}

// GetProvider returns a provider from its index
func (k Keeper) GetProvider(
	ctx sdk.Context,
	index string,

) (val types.Provider, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProviderKeyPrefix))

	b := store.Get(types.ProviderKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveProvider removes a provider from the store
func (k Keeper) RemoveProvider(
	ctx sdk.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProviderKeyPrefix))
	store.Delete(types.ProviderKey(
		index,
	))
}

// GetAllProvider returns all provider
func (k Keeper) GetAllProvider(ctx sdk.Context) (list []types.Provider) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ProviderKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Provider
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
