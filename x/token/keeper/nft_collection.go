package keeper

import (
	"github.com/LimeChain/mantrachain/x/token/types"
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

func (k Keeper) GetAllNftCollection(ctx sdk.Context, creator sdk.AccAddress) (list []types.NftCollection) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftCollectionStoreKey(creator))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NftCollection
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
