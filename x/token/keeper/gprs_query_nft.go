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

func (k Keeper) Nft(c context.Context, req *types.QueryGetNftRequest) (*types.QueryGetNftResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err = types.ValidateNftId(conf.ValidNftId, req.Id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collectionIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	index := types.GetNftIndex(collectionIndex, req.Id)

	meta, found := k.GetNft(
		ctx,
		collectionIndex,
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetNftResponse{
		Nft: &meta,
	}, nil
}

func (k Keeper) AllCollectionNfts(goCtx context.Context, req *types.QueryGetAllCollectionNftsRequest) (*types.QueryGetAllCollectionNftsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collectionIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	store := ctx.KVStore(k.storeKey)
	nftsStore := prefix.NewStore(store, types.NftStoreKey(collectionIndex))

	var nfts []*types.Nft
	pageRes, err := query.Paginate(nftsStore, req.Pagination, func(_ []byte, value []byte) error {
		var nftMeta types.Nft
		if err := k.cdc.Unmarshal(value, &nftMeta); err != nil {
			return err
		}
		nfts = append(nfts, &nftMeta)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryGetAllCollectionNftsResponse{
		CollectionCreator: collectionCreator.String(),
		CollectionId:      req.CollectionId,
		Nfts:              nfts,
		Pagination:        pageRes,
	}, nil
}

func (k Keeper) NftApproved(c context.Context, req *types.QueryGetNftApprovedRequest) (*types.QueryGetNftApprovedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err = types.ValidateNftId(conf.ValidNftId, req.Id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collectionIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	index := types.GetNftIndex(collectionIndex, req.Id)

	if !k.HasNft(
		ctx,
		collectionIndex,
		index,
	) {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	owner := nftExecutor.GetNftOwner(string(collectionIndex), string(index))

	var approved []string
	approvals := k.GetNftApproved(ctx, collectionIndex, index, owner)
	for operator, approval := range approvals {
		if approval != nil {
			approved = append(approved, operator)
		}
	}

	return &types.QueryGetNftApprovedResponse{
		Id:       req.Id,
		Approved: approved,
	}, nil
}

func (k Keeper) IsApprovedForAllNfts(c context.Context, req *types.QueryGetIsApprovedForAllNftsRequest) (*types.QueryGetIsApprovedForAllNftsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	owner, err := sdk.AccAddressFromBech32(req.Owner)

	if err != nil {
		return nil, err
	}

	operator, err := sdk.AccAddressFromBech32(req.Operator)

	if err != nil {
		return nil, err
	}

	return &types.QueryGetIsApprovedForAllNftsResponse{
		Operator: operator.String(),
		Approved: k.GetIsApprovedForAllNfts(ctx, owner, operator),
	}, nil
}
