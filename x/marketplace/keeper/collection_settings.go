package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetCollectionSettings(ctx sdk.Context, collectionSettings types.CollectionSettings) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CollectionSettingsStoreKey(collectionSettings.MarketplaceIndex))
	b := k.cdc.MustMarshal(&collectionSettings)
	store.Set(collectionSettings.Index, b)
}

func (k Keeper) GetCollectionSettings(
	ctx sdk.Context,
	marketplaceIndex []byte,
	index []byte,
) (val types.CollectionSettings, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CollectionSettingsStoreKey(marketplaceIndex))

	if !k.HasCollectionSettings(ctx, marketplaceIndex, index) {
		return types.CollectionSettings{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasCollectionSettings(
	ctx sdk.Context,
	marketplaceIndex []byte,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CollectionSettingsStoreKey(marketplaceIndex))
	return store.Has(index)
}

func (k Keeper) GetAllCollectionSettings(ctx sdk.Context, marketplaceIndex []byte) (list []types.CollectionSettings) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CollectionSettingsStoreKey(marketplaceIndex))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CollectionSettings
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
