package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetMarketplaceNft(ctx sdk.Context, nft types.MarketplaceNft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceNftStoreKey(nft.MarketplaceIndex, nft.CollectionIndex))
	b := k.cdc.MustMarshal(&nft)
	store.Set(nft.Index, b)
}

func (k Keeper) GetMarketplaceNft(
	ctx sdk.Context,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
) (val types.MarketplaceNft, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceNftStoreKey(marketplaceIndex, collectionIndex))

	if !k.HasMarketplaceNft(ctx, marketplaceIndex, collectionIndex, index) {
		return types.MarketplaceNft{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasMarketplaceNft(
	ctx sdk.Context,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceNftStoreKey(marketplaceIndex, collectionIndex))
	return store.Has(index)
}

func (k Keeper) GetAllMarketplaceNft(ctx sdk.Context, marketplaceIndex []byte, collectionIndex []byte) (list []types.MarketplaceNft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceNftStoreKey(marketplaceIndex, collectionIndex))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MarketplaceNft
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
