package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/nft/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Mint defines a method for minting a new nft
func (k Keeper) Mint(ctx sdk.Context, token types.NFT, receiver sdk.AccAddress) error {
	if !k.HasClass(ctx, token.ClassId) {
		return sdkerrors.Wrap(types.ErrClassNotExists, token.ClassId)
	}

	if k.HasNFT(ctx, token.ClassId, token.Id) {
		return sdkerrors.Wrap(types.ErrNFTExists, token.Id)
	}

	k.setNFT(ctx, token)
	k.setOwner(ctx, token.ClassId, token.Id, receiver)
	k.incrTotalSupply(ctx, token.ClassId)

	ctx.EventManager().EmitTypedEvent(&types.EventMint{
		ClassId: token.ClassId,
		Id:      token.Id,
		Owner:   receiver.String(),
	})
	return nil
}

// Burn defines a method for burning a nft from a specific account.
// Note: When the upper module uses this method, it needs to authenticate nft
func (k Keeper) Burn(ctx sdk.Context, classID string, nftID string) error {
	if !k.HasClass(ctx, classID) {
		return sdkerrors.Wrap(types.ErrClassNotExists, classID)
	}

	if !k.HasNFT(ctx, classID, nftID) {
		return sdkerrors.Wrap(types.ErrNFTNotExists, nftID)
	}

	owner := k.GetOwner(ctx, classID, nftID)
	nftStore := k.getNFTStore(ctx, classID)
	nftStore.Delete([]byte(nftID))

	k.deleteOwner(ctx, classID, nftID, owner)
	k.decrTotalSupply(ctx, classID)
	ctx.EventManager().EmitTypedEvent(&types.EventBurn{
		ClassId: classID,
		Id:      nftID,
		Owner:   owner.String(),
	})
	return nil
}

// Update defines a method for updating an exist nft
// Note: When the upper module uses this method, it needs to authenticate nft
func (k Keeper) Update(ctx sdk.Context, token types.NFT) error {
	if !k.HasClass(ctx, token.ClassId) {
		return sdkerrors.Wrap(types.ErrClassNotExists, token.ClassId)
	}

	if !k.HasNFT(ctx, token.ClassId, token.Id) {
		return sdkerrors.Wrap(types.ErrNFTNotExists, token.Id)
	}
	k.setNFT(ctx, token)
	return nil
}

// Transfer defines a method for sending a nft from one account to another account.
// Note: When the upper module uses this method, it needs to authenticate nft
func (k Keeper) Transfer(ctx sdk.Context,
	classID string,
	nftID string,
	receiver sdk.AccAddress,
) error {
	if !k.HasClass(ctx, classID) {
		return sdkerrors.Wrap(types.ErrClassNotExists, classID)
	}

	if !k.HasNFT(ctx, classID, nftID) {
		return sdkerrors.Wrap(types.ErrNFTNotExists, nftID)
	}

	owner := k.GetOwner(ctx, classID, nftID)
	k.deleteOwner(ctx, classID, nftID, owner)
	k.setOwner(ctx, classID, nftID, receiver)
	return nil
}

// GetNFT returns the nft information of the specified classID and nftID
func (k Keeper) GetNFT(ctx sdk.Context, classID, nftID string) (types.NFT, bool) {
	store := k.getNFTStore(ctx, classID)
	bz := store.Get([]byte(nftID))
	if len(bz) == 0 {
		return types.NFT{}, false
	}
	var nft types.NFT
	k.cdc.MustUnmarshal(bz, &nft)
	return nft, true
}

// GetNFTsOfClassByOwner returns all nft information of the specified classID under the specified owner
func (k Keeper) GetNFTsOfClassByOwner(ctx sdk.Context, classID string, owner sdk.AccAddress) (nfts []types.NFT) {
	ownerStore := k.getClassStoreByOwner(ctx, owner, classID)
	iterator := ownerStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nft, has := k.GetNFT(ctx, classID, string(iterator.Key()))
		if has {
			nfts = append(nfts, nft)
		}
	}
	return nfts
}

// GetNFTsOfClass returns all nft information under the specified classID
func (k Keeper) GetNFTsOfClass(ctx sdk.Context, classID string) (nfts []types.NFT) {
	nftStore := k.getNFTStore(ctx, classID)
	iterator := nftStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var nft types.NFT
		k.cdc.MustUnmarshal(iterator.Value(), &nft)
		nfts = append(nfts, nft)
	}
	return nfts
}

// GetOwner returns the owner information of the specified nft
func (k Keeper) GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(ownerStoreKey(classID, nftID))
	return sdk.AccAddress(bz)
}

// GetBalance returns the specified account, the number of all nfts under the specified classID
func (k Keeper) GetBalance(ctx sdk.Context, classID string, owner sdk.AccAddress) uint64 {
	nfts := k.GetNFTsOfClassByOwner(ctx, classID, owner)
	return uint64(len(nfts))
}

// GetTotalSupply returns the number of all nfts under the specified classID
func (k Keeper) GetTotalSupply(ctx sdk.Context, classID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(classTotalSupply(classID))
	return sdk.BigEndianToUint64(bz)
}

// HasNFT determines whether the specified classID and nftID exist
func (k Keeper) HasNFT(ctx sdk.Context, classID, id string) bool {
	store := k.getNFTStore(ctx, classID)
	return store.Has([]byte(id))
}

func (k Keeper) setNFT(ctx sdk.Context, token types.NFT) {
	nftStore := k.getNFTStore(ctx, token.ClassId)
	bz := k.cdc.MustMarshal(&token)
	nftStore.Set([]byte(token.Id), bz)
}

func (k Keeper) setOwner(ctx sdk.Context, classID, nftID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(ownerStoreKey(classID, nftID), owner.Bytes())

	ownerStore := k.getClassStoreByOwner(ctx, owner, classID)
	ownerStore.Set([]byte(nftID), Placeholder)
}

func (k Keeper) deleteOwner(ctx sdk.Context, classID, nftID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(ownerStoreKey(classID, nftID))

	ownerStore := k.getClassStoreByOwner(ctx, owner, classID)
	ownerStore.Delete([]byte(nftID))
}

func (k Keeper) getNFTStore(ctx sdk.Context, classID string) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, nftStoreKey(classID))
}

func (k Keeper) getClassStoreByOwner(ctx sdk.Context, owner sdk.AccAddress, classID string) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	key := nftOfClassByOwnerStoreKey(owner, classID)
	return prefix.NewStore(store, key)
}

func (k Keeper) prefixStoreNftOfClassByOwner(ctx sdk.Context, owner sdk.AccAddress) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	key := prefixNftOfClassByOwnerStoreKey(owner)
	return prefix.NewStore(store, key)
}

func (k Keeper) incrTotalSupply(ctx sdk.Context, classID string) {
	supply := k.GetTotalSupply(ctx, classID) + 1
	k.updateTotalSupply(ctx, classID, supply)
}

func (k Keeper) decrTotalSupply(ctx sdk.Context, classID string) {
	supply := k.GetTotalSupply(ctx, classID) - 1
	k.updateTotalSupply(ctx, classID, supply)
}

func (k Keeper) updateTotalSupply(ctx sdk.Context, classID string, supply uint64) {
	store := ctx.KVStore(k.storeKey)
	supplyKey := classTotalSupply(classID)
	store.Set(supplyKey, sdk.Uint64ToBigEndian(supply))
}

func (k Keeper) FilterNotOwnNFTsIdsOfClass(ctx sdk.Context, classID string, nftIDs []string, owner sdk.AccAddress) (list []string) {
	for _, id := range nftIDs {
		if owner.Equals(k.GetOwner(ctx, classID, id)) {
			list = append(list, id)
		}
	}

	return
}

func (k Keeper) MintBatch(ctx sdk.Context, tokens []types.NFT, receiver sdk.AccAddress) error {
	for _, token := range tokens {
		err := k.Mint(ctx, token, receiver)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) BurnBatch(ctx sdk.Context, classID string, nftIDs []string) error {
	for _, id := range nftIDs {
		err := k.Burn(ctx, classID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) TransferBatch(ctx sdk.Context, classID string, nftIDs []string, receiver sdk.AccAddress) error {
	for _, id := range nftIDs {
		err := k.Transfer(ctx, classID, id, receiver)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GetNFTsByIds(ctx sdk.Context, classID string, nftIDs []string) (nfts []types.NFT) {
	store := k.getNFTStore(ctx, classID)
	for _, nftID := range nftIDs {
		bz := store.Get([]byte(nftID))

		var nft types.NFT
		if len(bz) != 0 {
			k.cdc.MustUnmarshal(bz, &nft)
		}
		nfts = append(nfts, nft)
	}

	return
}
