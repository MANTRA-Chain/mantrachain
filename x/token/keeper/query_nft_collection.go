package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/runtime"
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

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftCollectionStoreKey(creator))

	var collections []*types.NftCollection
	pageRes, err := query.Paginate(store, req.Pagination, func(_ []byte, value []byte) error {
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

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.NftCollectionStoreKey(nil))

	var collections []*types.NftCollection
	pageRes, err := query.Paginate(store, req.Pagination, func(_ []byte, value []byte) error {
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

	isOpened := k.HasOpenedNftsCollection(
		ctx,
		req.Index,
	)

	return &types.QueryGetOpenedNftsCollectionResponse{
		Index:                req.Index,
		OpenedNftsCollection: isOpened,
	}, nil
}

func (k Keeper) RestrictedNftsCollection(goCtx context.Context, req *types.QueryGetRestrictedNftsCollectionRequest) (*types.QueryGetRestrictedNftsCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	isRestricted := k.HasRestrictedNftsCollection(
		ctx,
		req.Index,
	)

	return &types.QueryGetRestrictedNftsCollectionResponse{
		Index:                    req.Index,
		RestrictedNftsCollection: isRestricted,
	}, nil
}

func (k Keeper) SoulBondedNftsCollection(goCtx context.Context, req *types.QueryGetSoulBondedNftsCollectionRequest) (*types.QueryGetSoulBondedNftsCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	isSoulBonded := k.HasSoulBondedNftsCollection(
		ctx,
		req.Index,
	)

	return &types.QueryGetSoulBondedNftsCollectionResponse{
		Index:                    req.Index,
		SoulBondedNftsCollection: isSoulBonded,
	}, nil
}
