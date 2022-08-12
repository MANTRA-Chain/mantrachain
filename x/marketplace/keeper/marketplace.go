package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetMarketplace(ctx sdk.Context, marketplace types.Marketplace) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceStoreKey(marketplace.Creator))
	b := k.cdc.MustMarshal(&marketplace)
	store.Set(marketplace.Index, b)
}

func (k Keeper) GetMarketplace(
	ctx sdk.Context,
	creator sdk.AccAddress,
	index []byte,
) (val types.Marketplace, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceStoreKey(creator))

	if !k.HasMarketplace(ctx, creator, index) {
		return types.Marketplace{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasMarketplace(
	ctx sdk.Context,
	creator sdk.AccAddress,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceStoreKey(creator))
	return store.Has(index)
}

func (k Keeper) GetAllMarketplace(ctx sdk.Context, creator sdk.AccAddress) (list []types.Marketplace) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MarketplaceStoreKey(creator))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Marketplace
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
