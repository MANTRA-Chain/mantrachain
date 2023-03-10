package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) NftCollection(c context.Context, req *types.QueryGetNftCollectionRequest) (*types.QueryGetNftCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	creator, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)
	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.Id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	index := types.GetNftCollectionIndex(creator, req.Id)

	meta, found := k.GetNftCollection(
		ctx,
		sdk.AccAddress(creator),
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetNftCollectionResponse{
		NftCollection: &meta,
	}, nil
}

func (k Keeper) NftCollectionsByCreator(c context.Context, req *types.QueryGetNftCollectionsByCreatorRequest) (*types.QueryGetNftCollectionsByCreatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	creator, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	store := ctx.KVStore(k.storeKey)
	collectionsStore := prefix.NewStore(store, types.NftCollectionStoreKey(creator))

	var collections []*types.NftCollection
	pageRes, err := query.Paginate(collectionsStore, req.Pagination, func(_ []byte, value []byte) error {
		var collection types.NftCollection
		if err := k.cdc.Unmarshal(value, &collection); err != nil {
			return err
		}
		collections = append(collections, &collection)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryGetNftCollectionsByCreatorResponse{
		Creator:        creator.String(),
		NftCollections: collections,
		Pagination:     pageRes,
	}, nil
}

func (k Keeper) AllNftCollections(c context.Context, req *types.QueryGetAllNftCollectionsRequest) (*types.QueryGetAllNftCollectionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	collectionsStore := prefix.NewStore(store, types.NftCollectionStoreKey(nil))

	var collections []*types.NftCollection
	pageRes, err := query.Paginate(collectionsStore, req.Pagination, func(_ []byte, value []byte) error {
		var collection types.NftCollection
		if err := k.cdc.Unmarshal(value, &collection); err != nil {
			return err
		}
		collections = append(collections, &collection)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryGetAllNftCollectionsResponse{
		NftCollections: collections,
		Pagination:     pageRes,
	}, nil
}

func (k Keeper) NftCollectionOwner(goCtx context.Context, req *types.QueryGetNftCollectionOwnerRequest) (*types.QueryGetNftCollectionOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetNftCollectionOwner(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetNftCollectionOwnerResponse{NftCollectionOwner: string(val)}, nil
}

func (k Keeper) OpenedNftsCollection(goCtx context.Context, req *types.QueryGetOpenedNftsCollectionRequest) (*types.QueryGetOpenedNftsCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryGetOpenedNftsCollectionResponse{OpenedNftsCollection: k.HasOpenedNftsCollection(
		ctx,
		req.Index,
	)}, nil
}

func (k Keeper) RestrictedNftsCollection(goCtx context.Context, req *types.QueryGetRestrictedNftsCollectionRequest) (*types.QueryGetRestrictedNftsCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryGetRestrictedNftsCollectionResponse{RestrictedNftsCollection: k.HasRestrictedNftsCollection(
		ctx,
		req.Index,
	)}, nil
}

func (k Keeper) SoulBondedNftsCollection(goCtx context.Context, req *types.QueryGetSoulBondedNftsCollectionRequest) (*types.QueryGetSoulBondedNftsCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryGetSoulBondedNftsCollectionResponse{SoulBondedNftsCollection: k.HasSoulBondedNftsCollection(
		ctx,
		req.Index,
	)}, nil
}
