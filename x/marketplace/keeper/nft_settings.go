package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetNftSettings(ctx sdk.Context, nftSettings types.NftSettings) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftSettingsStoreKey(nftSettings.MarketplaceIndex, nftSettings.CollectionIndex))
	b := k.cdc.MustMarshal(&nftSettings)
	store.Set(nftSettings.Index, b)
}

func (k Keeper) GetNftSettings(
	ctx sdk.Context,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
) (val types.NftSettings, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftSettingsStoreKey(marketplaceIndex, collectionIndex))

	if !k.HasNftSettings(ctx, marketplaceIndex, collectionIndex, index) {
		return types.NftSettings{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasNftSettings(
	ctx sdk.Context,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftSettingsStoreKey(marketplaceIndex, collectionIndex))
	return store.Has(index)
}

func (k Keeper) GetAllNftSettings(ctx sdk.Context, marketplaceIndex []byte, collectionIndex []byte) (list []types.NftSettings) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftSettingsStoreKey(marketplaceIndex, collectionIndex))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NftSettings
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
