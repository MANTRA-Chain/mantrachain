package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetNft(ctx sdk.Context, nft types.Nft) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(nft.CollectionIndex))
	b := k.cdc.MustMarshal(&nft)
	store.Set(nft.Index, b)
}

func (k Keeper) SetNfts(ctx sdk.Context, nfts []types.Nft) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	for _, nft := range nfts {
		store := prefix.NewStore(storeAdapter, types.NftStoreKey(nft.CollectionIndex))
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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftApprovedStoreKey(collectionIndex))
	bz := store.Get(nftIndex)
	var approvedAddressesList types.ApprovedAddressesList
	noOp := true

	if bz == nil && approved {
		approvedAddressesList = types.ApprovedAddressesList{List: map[string]*types.ApprovedAddresses{owner.String(): {Addresses: map[string][]byte{receiver.String(): types.Placeholder}}}}

		noOp = false
	} else if bz != nil {
		k.cdc.MustUnmarshal(bz, &approvedAddressesList)
		approvedAddresses := approvedAddressesList.List[owner.String()]

		if approvedAddresses == nil && approved {
			approvedAddressesList.List[owner.String()] = &types.ApprovedAddresses{Addresses: map[string][]byte{receiver.String(): types.Placeholder}}

			noOp = false
		} else if approved && approvedAddresses.Addresses[receiver.String()] == nil {
			approvedAddresses.Addresses[receiver.String()] = types.Placeholder

			noOp = false
		} else if !approved && approvedAddresses.Addresses[receiver.String()] != nil {
			delete(approvedAddresses.Addresses, receiver.String())

			if len(approvedAddressesList.List[owner.String()].Addresses) == 0 {
				delete(approvedAddressesList.List, owner.String())
			}

			noOp = false
		}
	}

	if !noOp {
		if len(approvedAddressesList.List) == 0 {
			store.Delete(nftIndex)
		} else {
			store.Set(nftIndex, k.cdc.MustMarshal(&approvedAddressesList))
		}
	}
}

func (k Keeper) IsApproved(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
	owner sdk.AccAddress,
	operator sdk.AccAddress,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftApprovedAllStoreKey())
	index := types.GetNftApprovedAllIndex(owner)

	bz := store.Get(index)
	var approvedAddresses types.ApprovedAddresses
	k.cdc.MustUnmarshal(bz, &approvedAddresses)

	if approvedAddresses.Addresses[operator.String()] != nil {
		return true
	}

	store = prefix.NewStore(storeAdapter, types.NftApprovedStoreKey(collectionIndex))
	bz = store.Get(nftIndex)

	var approvedAddressesList types.ApprovedAddressesList
	k.cdc.MustUnmarshal(bz, &approvedAddressesList)

	return approvedAddressesList.List[owner.String()] != nil &&
		approvedAddressesList.List[owner.String()].Addresses[operator.String()] != nil
}

func (k Keeper) GetNftApproved(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
	owner sdk.AccAddress,
) map[string][]byte {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftApprovedStoreKey(collectionIndex))
	bz := store.Get(nftIndex)
	var approvedAddressesList types.ApprovedAddressesList
	k.cdc.MustUnmarshal(bz, &approvedAddressesList)
	if approvedAddressesList.List[owner.String()] != nil {
		return approvedAddressesList.List[owner.String()].Addresses
	}
	return nil
}

func (k Keeper) GetIsApprovedForAllNfts(
	ctx sdk.Context,
	owner sdk.AccAddress,
	operator sdk.AccAddress,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftApprovedAllStoreKey())
	index := types.GetNftApprovedAllIndex(owner)
	bz := store.Get(index)
	var approvedAddresses types.ApprovedAddresses
	k.cdc.MustUnmarshal(bz, &approvedAddresses)

	return approvedAddresses.Addresses[operator.String()] != nil
}

func (k Keeper) DeleteApprovedNft(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftApprovedStoreKey(collectionIndex))
	store.Delete(nftIndex)
}

func (k Keeper) DeleteApprovedNfts(
	ctx sdk.Context,
	collectionIndex []byte,
	nftsIndexes [][]byte,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftApprovedStoreKey(collectionIndex))
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
	for _, nftIndex := range nftsIndexes {
		k.SetApprovedNft(ctx, collectionIndex, nftIndex, owner, receiver, approved)
	}
}

func (k Keeper) SetApprovedAllNfts(
	ctx sdk.Context,
	owner sdk.AccAddress,
	receiver sdk.AccAddress,
	approved bool,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftApprovedAllStoreKey())
	index := types.GetNftApprovedAllIndex(owner)
	bz := store.Get(index)
	var approvedAddresses types.ApprovedAddresses
	noOp := true

	if bz == nil && approved {
		approvedAddresses = types.ApprovedAddresses{Addresses: map[string][]byte{receiver.String(): types.Placeholder}}

		noOp = false
	} else if bz != nil {
		k.cdc.MustUnmarshal(bz, &approvedAddresses)

		if approvedAddresses.Addresses[receiver.String()] == nil && approved {
			approvedAddresses.Addresses[receiver.String()] = types.Placeholder

			noOp = false
		} else if !approved && approvedAddresses.Addresses[receiver.String()] != nil {
			delete(approvedAddresses.Addresses, receiver.String())

			noOp = false
		}
	}

	if !noOp {
		if len(approvedAddresses.Addresses) == 0 {
			store.Delete(index)
		} else {
			store.Set(index, k.cdc.MustMarshal(&approvedAddresses))
		}
	}
}

func (k Keeper) GetNft(
	ctx sdk.Context,
	collectionIndex []byte,
	nftIndex []byte,
) (val types.Nft, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(collectionIndex))

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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(collectionIndex))
	return store.Has(nftIndex)
}

func (k Keeper) FilterNotExist(ctx sdk.Context, collectionIndex []byte, nftsIndexes [][]byte) (list [][]byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		if store.Has(nftIndex) {
			list = append(list, nftIndex)
		}
	}

	return
}

func (k Keeper) FilterExist(ctx sdk.Context, collectionIndex []byte, nftsIndexes [][]byte) (list [][]byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		if !store.Has(nftIndex) {
			list = append(list, nftIndex)
		}
	}

	return
}

func (k Keeper) DeleteNft(ctx sdk.Context, collectionIndex []byte, nftIndex []byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(collectionIndex))
	store.Delete(nftIndex)
}

func (k Keeper) DeleteNfts(ctx sdk.Context, collectionIndex []byte, nftsIndexes [][]byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		store.Delete(nftIndex)
	}
}

func (k Keeper) GetAllNft(ctx sdk.Context) (list []types.Nft) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(nil))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Nft
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetNftsByIndexes(ctx sdk.Context, collectionIndex []byte, nftsIndexes [][]byte) (list []types.Nft) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftStoreKey(collectionIndex))
	for _, nftIndex := range nftsIndexes {
		bz := store.Get(nftIndex)

		var nft types.Nft
		if len(bz) != 0 {
			k.cdc.MustUnmarshal(bz, &nft)
		}
		list = append(list, nft)
	}

	return
}

func (k Keeper) TransferNft(
	ctx sdk.Context,
	operator sdk.AccAddress,
	owner sdk.AccAddress,
	receiver sdk.AccAddress,
	collectionIndex []byte,
	index []byte,
) error {
	if !owner.Equals(operator) && !k.IsApproved(ctx, collectionIndex, index, owner, operator) {
		return errors.Wrap(types.ErrInvalidNft, "operator not approved to transfer nft")
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	err := nftExecutor.TransferNft(string(collectionIndex), string(index), receiver)
	if err != nil {
		return err
	}

	k.DeleteApprovedNft(ctx, collectionIndex, index)

	return nil
}
