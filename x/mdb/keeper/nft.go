package keeper

import (
	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetNft set a specific nft in the store from its index
func (k Keeper) SetNft(ctx sdk.Context, nft types.Nft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(nft.CollectionIndex))
	b := k.cdc.MustMarshal(&nft)
	store.Set(nft.Index, b)
}

// GetNft returns a nft from its nft
func (k Keeper) GetNft(
	ctx sdk.Context,
	collIndex []byte,
	index []byte,
) (val types.Nft, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collIndex))

	if !k.HasNft(ctx, collIndex, index) {
		return types.Nft{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// HasNft checks if the nft exists in the store
func (k Keeper) HasNft(
	ctx sdk.Context,
	collIndex []byte,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collIndex))
	return store.Has(index)
}

// GetAllNft returns all nft
func (k Keeper) GetAllNft(ctx sdk.Context, collIndex []byte) (list []types.Nft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collIndex))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Nft
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
