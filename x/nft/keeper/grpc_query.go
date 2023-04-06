package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/nft/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = Keeper{}

// Balance return the number of NFTs of a given class owned by the owner, same as balanceOf in ERC721
func (k Keeper) Balance(goCtx context.Context, r *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if err := types.ValidateClassID(r.ClassId); err != nil {
		return nil, err
	}

	owner, err := sdk.AccAddressFromBech32(r.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	balance := k.GetBalance(ctx, r.ClassId, owner)
	return &types.QueryBalanceResponse{Amount: balance}, nil
}

// Owner return the owner of the NFT based on its class and id, same as ownerOf in ERC721
func (k Keeper) Owner(goCtx context.Context, r *types.QueryOwnerRequest) (*types.QueryOwnerResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if err := types.ValidateClassID(r.ClassId); err != nil {
		return nil, err
	}

	if err := types.ValidateNFTID(r.Id); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := k.GetOwner(ctx, r.ClassId, r.Id)
	return &types.QueryOwnerResponse{Owner: owner.String()}, nil
}

// Supply return the number of NFTs from the given class, same as totalSupply of ERC721.
func (k Keeper) Supply(goCtx context.Context, r *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if err := types.ValidateClassID(r.ClassId); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	supply := k.GetTotalSupply(ctx, r.ClassId)
	return &types.QuerySupplyResponse{Amount: supply}, nil
}

// NFTs queries all NFTs of a given class or owner (at least one must be provided), similar to tokenByIndex in ERC721Enumerable
func (k Keeper) NFTs(goCtx context.Context, r *types.QueryNFTsRequest) (*types.QueryNFTsResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	var err error
	var owner sdk.AccAddress
	if len(r.ClassId) > 0 {
		if err := types.ValidateClassID(r.ClassId); err != nil {
			return nil, err
		}
	}

	if len(r.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(r.Owner)
		if err != nil {
			return nil, err
		}
	}

	var nfts []*types.NFT
	var pageRes *query.PageResponse
	ctx := sdk.UnwrapSDKContext(goCtx)

	switch {
	case len(r.ClassId) > 0 && len(r.Owner) > 0:
		if pageRes, err = query.Paginate(k.getClassStoreByOwner(ctx, owner, r.ClassId), r.Pagination, func(key []byte, _ []byte) error {
			nft, has := k.GetNFT(ctx, r.ClassId, string(key))
			if has {
				nfts = append(nfts, &nft)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	case len(r.ClassId) > 0 && len(r.Owner) == 0:
		nftStore := k.getNFTStore(ctx, r.ClassId)
		if pageRes, err = query.Paginate(nftStore, r.Pagination, func(_ []byte, value []byte) error {
			var nft types.NFT
			if err := k.cdc.Unmarshal(value, &nft); err != nil {
				return err
			}
			nfts = append(nfts, &nft)
			return nil
		}); err != nil {
			return nil, err
		}
	case len(r.ClassId) == 0 && len(r.Owner) > 0:
		if pageRes, err = query.Paginate(k.prefixStoreNftOfClassByOwner(ctx, owner), r.Pagination, func(key []byte, value []byte) error {
			classID, nftID := parseNftOfClassByOwnerStoreKey(key)
			if n, has := k.GetNFT(ctx, classID, nftID); has {
				nfts = append(nfts, &n)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.ErrInvalidRequest.Wrap("must provide at least one of classID or owner")
	}
	return &types.QueryNFTsResponse{
		Nfts:       nfts,
		Pagination: pageRes,
	}, nil
}

// NFT return an NFT based on its class and id.
func (k Keeper) NFT(goCtx context.Context, r *types.QueryNFTRequest) (*types.QueryNFTResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if err := types.ValidateClassID(r.ClassId); err != nil {
		return nil, err
	}
	if err := types.ValidateNFTID(r.Id); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	n, has := k.GetNFT(ctx, r.ClassId, r.Id)
	if !has {
		return nil, types.ErrNFTNotExists.Wrapf("not found nft, class: %s, id: %s", r.ClassId, r.Id)
	}
	return &types.QueryNFTResponse{Nft: &n}, nil
}

// Class return an NFT class based on its id
func (k Keeper) Class(goCtx context.Context, r *types.QueryClassRequest) (*types.QueryClassResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	if err := types.ValidateClassID(r.ClassId); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	class, has := k.GetClass(ctx, r.ClassId)
	if !has {
		return nil, types.ErrClassNotExists.Wrapf("not found class: %s", r.ClassId)
	}
	return &types.QueryClassResponse{Class: &class}, nil
}

// Classes return all NFT classes
func (k Keeper) Classes(goCtx context.Context, r *types.QueryClassesRequest) (*types.QueryClassesResponse, error) {
	if r == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	classStore := prefix.NewStore(store, ClassKey)

	var classes []*types.Class
	pageRes, err := query.Paginate(classStore, r.Pagination, func(_ []byte, value []byte) error {
		var class types.Class
		if err := k.cdc.Unmarshal(value, &class); err != nil {
			return err
		}
		classes = append(classes, &class)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryClassesResponse{
		Classes:    classes,
		Pagination: pageRes,
	}, nil
}
