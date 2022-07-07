package keeper

import (
	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetNft(ctx sdk.Context, nft types.Nft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(nft.CollectionIndex))
	b := k.cdc.MustMarshal(&nft)
	store.Set(nft.Index, b)
}

func (k Keeper) SetNfts(ctx sdk.Context, nfts []types.Nft) {
	for _, nft := range nfts {
		store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(nft.CollectionIndex))
		b := k.cdc.MustMarshal(&nft)
		store.Set(nft.Index, b)
	}
}

func (k Keeper) SetApprovedNft(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
	owner sdk.AccAddress,
	receiver sdk.AccAddress,
	approved bool,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftApprovedStoreKey(collectionIndex))
	bz := store.Get(nftIndex)
	var approvedAddresses types.ApprovedAddresses
	k.cdc.MustUnmarshal(bz, &approvedAddresses)
	approvedAddressesData := approvedAddresses.ApprovedAddresses[owner.String()]
	approvedAddressesData.ApprovedAddressesData[receiver.String()] = approved
	bz = k.cdc.MustMarshal(&approvedAddresses)
	store.Set(nftIndex, bz)
}

func (k Keeper) IsApproved(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
	owner sdk.AccAddress,
	operator sdk.AccAddress,
) bool {
	index := types.GetNftApprovedAllIndex(owner)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftApprovedAllStoreKey())
	bz := store.Get(index)
	var approvedAddressesData types.ApprovedAddressesData
	k.cdc.MustUnmarshal(bz, &approvedAddressesData)

	if approvedAddressesData.ApprovedAddressesData[operator.String()] {
		return true
	}

	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.NftApprovedStoreKey(collectionIndex))
	bz = store.Get(index)
	var approvedAddresses types.ApprovedAddresses
	k.cdc.MustUnmarshal(bz, &approvedAddresses)
	return approvedAddressesData.ApprovedAddressesData[operator.String()]
}

func (k Keeper) DeleteApprovedNft(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
	owner sdk.AccAddress,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftApprovedStoreKey(collectionIndex))
	store.Delete(nftIndex)
}

func (k Keeper) DeleteApprovedNfts(
	ctx sdk.Context,
	collectionIndex []byte,
	nftsIndexes [][]byte,
	owner sdk.AccAddress,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftApprovedStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		store.Delete(nftIndex)
	}
}

func (k Keeper) SetApprovedNfts(
	ctx sdk.Context,
	collectionIndex []byte,
	nftsIndexes [][]byte,
	owner sdk.AccAddress,
	receiver sdk.AccAddress,
	approved bool,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftApprovedStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		bz := store.Get(nftIndex)
		var approvedAddresses types.ApprovedAddresses
		k.cdc.MustUnmarshal(bz, &approvedAddresses)
		approvedAddressesData := approvedAddresses.ApprovedAddresses[owner.String()]
		approvedAddressesData.ApprovedAddressesData[receiver.String()] = approved
		bz = k.cdc.MustMarshal(&approvedAddresses)
		store.Set(nftIndex, bz)
	}
}

func (k Keeper) SetApprovedAllNfts(
	ctx sdk.Context,
	owner sdk.AccAddress,
	receiver sdk.AccAddress,
	approved bool,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftApprovedAllStoreKey())
	index := types.GetNftApprovedAllIndex(owner)
	bz := store.Get(index)
	var approvedAddressesData types.ApprovedAddressesData
	k.cdc.MustUnmarshal(bz, &approvedAddressesData)
	approvedAddressesData.ApprovedAddressesData[receiver.String()] = approved
	bz = k.cdc.MustMarshal(&approvedAddressesData)
	store.Set(index, bz)
}

func (k Keeper) GetNft(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
) (val types.Nft, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collectionIndex))

	if !k.HasNft(ctx, collectionIndex, nftIndex) {
		return types.Nft{}, false
	}

	b := store.Get(nftIndex)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasNft(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collectionIndex))
	return store.Has(nftIndex)
}

func (k Keeper) FilterNotExists(ctx sdk.Context, collectionIndex []byte, nftsIndexes [][]byte) (list [][]byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		if store.Has(nftIndex) {
			list = append(list, nftIndex)
		}
	}

	return
}

func (k Keeper) DeleteNft(ctx sdk.Context, collectionIndex []byte, nftIndex []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collectionIndex))
	store.Delete(nftIndex)
}

func (k Keeper) DeleteNfts(ctx sdk.Context, collectionIndex []byte, nftsIndexes [][]byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		store.Delete(nftIndex)
	}
}

func (k Keeper) GetAllNft(ctx sdk.Context, collectionIndex []byte) (list []types.Nft) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStoreKey(collectionIndex))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Nft
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
