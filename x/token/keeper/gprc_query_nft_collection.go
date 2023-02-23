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

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nftColl, found := nftExecutor.GetClass(string(index))
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetNftCollectionResponse{
		Id:          meta.Id,
		Name:        nftColl.Name,
		Symbol:      nftColl.Symbol,
		Description: nftColl.Description,
		Images:      meta.Images,
		Url:         meta.Url,
		Links:       meta.Links,
		Category:    meta.Category,
		Options:     meta.Options,
		Creator:     meta.Creator.String(),
		Owner:       meta.Owner.String(),
		Opened:      meta.Opened,
		Data:        nftColl.Data,
		// TODO: add is collection for soul bonded nfts field here
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

	var collectionsIndexes []string

	for _, collection := range collections {
		collectionsIndexes = append(collectionsIndexes, string(collection.Index))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nftColls := nftExecutor.GetClasses(collectionsIndexes)

	var nftCollections []*types.QueryGetNftCollectionResponse

	for i, nftColl := range nftColls {
		meta := collections[i]
		nftCollections = append(nftCollections, &types.QueryGetNftCollectionResponse{
			Id:          nftColl.UriHash,
			Name:        nftColl.Name,
			Symbol:      nftColl.Symbol,
			Description: nftColl.Description,
			Images:      meta.Images,
			Url:         meta.Url,
			Links:       meta.Links,
			Category:    meta.Category,
			Options:     meta.Options,
			Creator:     meta.Creator.String(),
			Owner:       meta.Owner.String(),
			Opened:      meta.Opened,
			Data:        nftColl.Data,
			// TODO: add is collection for soul bonded nfts field here
		})
	}

	return &types.QueryGetNftCollectionsByCreatorResponse{
		Creator:        creator.String(),
		NftCollections: nftCollections,
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

	var collectionsIndexes []string

	for _, collection := range collections {
		collectionsIndexes = append(collectionsIndexes, string(collection.Index))
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)
	nftColls := nftExecutor.GetClasses(collectionsIndexes)

	var nftCollections []*types.QueryGetNftCollectionResponse

	for i, nftColl := range nftColls {
		meta := collections[i]
		nftCollections = append(nftCollections, &types.QueryGetNftCollectionResponse{
			Id:          nftColl.UriHash,
			Name:        nftColl.Name,
			Symbol:      nftColl.Symbol,
			Description: nftColl.Description,
			Images:      meta.Images,
			Url:         meta.Url,
			Links:       meta.Links,
			Category:    meta.Category,
			Options:     meta.Options,
			Creator:     meta.Creator.String(),
			Owner:       meta.Owner.String(),
			Opened:      meta.Opened,
			Data:        nftColl.Data,
			// TODO: add is collection for soul bonded nfts field here
		})
	}

	return &types.QueryGetAllNftCollectionsResponse{
		NftCollections: nftCollections,
		Pagination:     pageRes,
	}, nil
}

func (k Keeper) NftCollectionSupply(c context.Context, req *types.QueryGetNftCollectionSupplyRequest) (*types.QueryGetNftCollectionSupplyResponse, error) {
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

	if !k.HasNftCollection(
		ctx,
		sdk.AccAddress(creator),
		index,
	) {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	nftExecutor := NewNftExecutor(ctx, k.nftKeeper)

	return &types.QueryGetNftCollectionSupplyResponse{
		Supply:  nftExecutor.GetClassSupply(string(index)),
		Creator: creator.String(),
		Id:      req.Id,
	}, nil
}
