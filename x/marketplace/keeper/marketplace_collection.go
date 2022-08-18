package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetMarketplaceCollection(ctx sdk.Context, marketplaceCollection types.MarketplaceCollection) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceCollectionStoreKey(marketplaceCollection.MarketplaceIndex))
	b := k.cdc.MustMarshal(&marketplaceCollection)
	store.Set(marketplaceCollection.Index, b)
}

func (k Keeper) GetMarketplaceCollection(
	ctx sdk.Context,
	marketplaceIndex []byte,
	index []byte,
) (val types.MarketplaceCollection, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceCollectionStoreKey(marketplaceIndex))

	if !k.HasMarketplaceCollection(ctx, marketplaceIndex, index) {
		return types.MarketplaceCollection{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasMarketplaceCollection(
	ctx sdk.Context,
	marketplaceIndex []byte,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceCollectionStoreKey(marketplaceIndex))
	return store.Has(index)
}

func (k Keeper) GetAllMarketplaceCollection(ctx sdk.Context, marketplaceIndex []byte) (list []types.MarketplaceCollection) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceCollectionStoreKey(marketplaceIndex))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MarketplaceCollection
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
