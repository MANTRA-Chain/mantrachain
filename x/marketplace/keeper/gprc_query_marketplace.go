package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Marketplace(c context.Context, req *types.QueryGetMarketplaceRequest) (*types.QueryGetMarketplaceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	creator, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)
	err = types.ValidateMarketplaceId(conf.ValidMarketplaceId, req.Id, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	index := types.GetMarketplaceIndex(creator, req.Id)

	marketplace, found := k.GetMarketplace(
		ctx,
		sdk.AccAddress(creator),
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetMarketplaceResponse{
		Id:          marketplace.Id,
		Name:        marketplace.Name,
		Description: marketplace.Description,
		Url:         marketplace.Url,
		Creator:     marketplace.Creator.String(),
		Owner:       marketplace.Owner.String(),
		Opened:      marketplace.Opened,
		Options:     marketplace.Options,
		Attributes:  marketplace.Attributes,
		Images:      marketplace.Images,
		Links:       marketplace.Links,
		Data:        marketplace.Data,
	}, nil
}

func (k Keeper) AllMarketplaces(c context.Context, req *types.QueryGetAllMarketplacesRequest) (*types.QueryGetAllMarketplacesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	marketplacesStore := prefix.NewStore(store, types.MarketplaceStoreKey(nil))

	var marketplaces []*types.Marketplace
	pageRes, err := query.Paginate(marketplacesStore, req.Pagination, func(_ []byte, value []byte) error {
		var marketplace types.Marketplace
		if err := k.cdc.Unmarshal(value, &marketplace); err != nil {
			return err
		}
		marketplaces = append(marketplaces, &marketplace)
		return nil
	})

	if err != nil {
		return nil, err
	}

	var marketplacesRes []*types.QueryGetMarketplaceResponse

	for _, marketplace := range marketplaces {
		marketplacesRes = append(marketplacesRes, &types.QueryGetMarketplaceResponse{
			Id:          marketplace.Id,
			Name:        marketplace.Name,
			Description: marketplace.Description,
			Url:         marketplace.Url,
			Creator:     marketplace.Creator.String(),
			Owner:       marketplace.Owner.String(),
			Opened:      marketplace.Opened,
			Options:     marketplace.Options,
			Attributes:  marketplace.Attributes,
			Images:      marketplace.Images,
			Links:       marketplace.Links,
			Data:        marketplace.Data,
		})
	}

	return &types.QueryGetAllMarketplacesResponse{
		Marketplaces: marketplacesRes,
		Pagination:   pageRes,
	}, nil
}

func (k Keeper) MarketplacesByCreator(c context.Context, req *types.QueryGetMarketplacesByCreatorRequest) (*types.QueryGetMarketplacesByCreatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	creator, err := sdk.AccAddressFromBech32(req.Creator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	store := ctx.KVStore(k.storeKey)
	marketplacesStore := prefix.NewStore(store, types.MarketplaceStoreKey(creator))

	var marketplaces []*types.Marketplace
	pageRes, err := query.Paginate(marketplacesStore, req.Pagination, func(_ []byte, value []byte) error {
		var marketplace types.Marketplace
		if err := k.cdc.Unmarshal(value, &marketplace); err != nil {
			return err
		}
		marketplaces = append(marketplaces, &marketplace)
		return nil
	})

	if err != nil {
		return nil, err
	}

	var marketplacesRes []*types.QueryGetMarketplaceResponse

	for _, marketplace := range marketplaces {
		marketplacesRes = append(marketplacesRes, &types.QueryGetMarketplaceResponse{
			Id:          marketplace.Id,
			Name:        marketplace.Name,
			Description: marketplace.Description,
			Url:         marketplace.Url,
			Creator:     marketplace.Creator.String(),
			Owner:       marketplace.Owner.String(),
			Opened:      marketplace.Opened,
			Options:     marketplace.Options,
			Attributes:  marketplace.Attributes,
			Images:      marketplace.Images,
			Links:       marketplace.Links,
			Data:        marketplace.Data,
		})
	}

	return &types.QueryGetMarketplacesByCreatorResponse{
		Creator:      creator.String(),
		Marketplaces: marketplacesRes,
		Pagination:   pageRes,
	}, nil
}
