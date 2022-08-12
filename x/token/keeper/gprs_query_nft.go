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

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err = types.ValidateNftId(conf.ValidNftId, req.Id, nil)

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

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)

	nft, found := nftExecutor.GetNft(string(collectionIndex), string(index))
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetNftResponse{
		Id:           nft.UriHash,
		Title:        meta.Title,
		Description:  meta.Description,
		Images:       meta.Images,
		Url:          meta.Url,
		Links:        meta.Links,
		Attributes:   meta.Attributes,
		Creator:      meta.Creator.String(),
		Owner:        nftExecutor.GetNftOwner(string(collectionIndex), string(index)).String(),
		Data:         nft.Data,
		CollectionId: meta.CollectionId,
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

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collectionIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	store := ctx.KVStore(k.storeKey)
	nftsMetaStore := prefix.NewStore(store, types.NftStoreKey(collectionIndex))

	var nftsMeta []*types.Nft
	pageRes, err := query.Paginate(nftsMetaStore, req.Pagination, func(_ []byte, value []byte) error {
		var nftMeta types.Nft
		if err := k.cdc.Unmarshal(value, &nftMeta); err != nil {
			return err
		}
		nftsMeta = append(nftsMeta, &nftMeta)
		return nil
	})

	if err != nil {
		return nil, err
	}

	var nftsIndexes []string

	for _, nftMeta := range nftsMeta {
		nftsIndexes = append(nftsIndexes, string(nftMeta.Index))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nfts := nftExecutor.GetNfts(string(collectionIndex), nftsIndexes)

	var nftsRes []*types.QueryGetNftResponse

	for i, nft := range nfts {
		meta := nftsMeta[i]

		if nft.Id == "" || meta.Index == nil {
			continue
		}

		nftsRes = append(nftsRes, &types.QueryGetNftResponse{
			Id:           nft.UriHash,
			Title:        meta.Title,
			Description:  meta.Description,
			Images:       meta.Images,
			Url:          meta.Url,
			Links:        meta.Links,
			Attributes:   meta.Attributes,
			Creator:      meta.Creator.String(),
			Owner:        nftExecutor.GetNftOwner(string(collectionIndex), string(meta.Index)).String(),
			Data:         nft.Data,
			CollectionId: meta.CollectionId,
		})
	}

	return &types.QueryGetAllCollectionNftsResponse{
		Id:         req.CollectionId,
		Nfts:       nftsRes,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) NftOwner(c context.Context, req *types.QueryGetNftOwnerRequest) (*types.QueryGetNftOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err = types.ValidateNftId(conf.ValidNftId, req.Id, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collectionIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	index := types.GetNftIndex(collectionIndex, req.Id)

	_, found := k.GetNft(
		ctx,
		collectionIndex,
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)

	return &types.QueryGetNftOwnerResponse{
		Id:      req.Id,
		Address: nftExecutor.GetNftOwner(string(collectionIndex), string(index)).String(),
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

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	err = types.ValidateNftId(conf.ValidNftId, req.Id, nil)

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
		if approval {
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
		Address:  operator.String(),
		Approved: k.GetIsApprovedForAllNfts(ctx, owner, operator),
	}, nil
}

func (k Keeper) NftBalance(c context.Context, req *types.QueryGetNftBalanceRequest) (*types.QueryGetNftBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	owner, err := sdk.AccAddressFromBech32(req.Owner)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collectionIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)

	return &types.QueryGetNftBalanceResponse{
		Balance: nftExecutor.GetNftBalance(string(collectionIndex), owner),
		Address: owner.String(),
	}, nil
}

func (k Keeper) CollectionNftsByOwner(goCtx context.Context, req *types.QueryGetCollectionNftsByOwnerRequest) (*types.QueryGetCollectionNftsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(req.Owner)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateNftCollectionId(conf.ValidNftCollectionId, req.CollectionId, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionIndex := types.GetNftCollectionIndex(collectionCreator, req.CollectionId)

	if !k.HasNftCollection(ctx, collectionCreator, collectionIndex) {
		return nil, status.Error(codes.InvalidArgument, "collection not exists")
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nfts := nftExecutor.GetNftsOfClassByOwner(string(collectionIndex), owner)

	var nftsIndexes [][]byte

	for _, nft := range nfts {
		nftIndex := types.GetNftIndex(collectionCreator, nft.UriHash)
		nftsIndexes = append(nftsIndexes, nftIndex)
	}

	nftsMeta := k.GetNftsByIndexes(ctx, collectionIndex, nftsIndexes)

	var nftsRes []*types.QueryGetNftResponse

	for i, nft := range nfts {
		meta := nftsMeta[i]

		if nft.Id == "" || meta.Index == nil {
			continue
		}

		nftsRes = append(nftsRes, &types.QueryGetNftResponse{
			Id:           nft.UriHash,
			Title:        meta.Title,
			Description:  meta.Description,
			Images:       meta.Images,
			Url:          meta.Url,
			Links:        meta.Links,
			Attributes:   meta.Attributes,
			Creator:      meta.Creator.String(),
			Owner:        owner.String(),
			Data:         nft.Data,
			CollectionId: meta.CollectionId,
		})
	}

	return &types.QueryGetCollectionNftsByOwnerResponse{
		Address: owner.String(),
		Nfts:    nftsRes,
	}, nil
}
