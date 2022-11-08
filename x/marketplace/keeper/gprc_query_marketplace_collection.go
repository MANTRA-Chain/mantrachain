package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MarketplaceCollection(c context.Context, req *types.QueryGetMarketplaceCollectionRequest) (*types.QueryGetMarketplaceCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	marketplaceCreator, err := sdk.AccAddressFromBech32(req.MarketplaceCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateMarketplaceId(conf.ValidMarketplaceId, req.MarketplaceId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Add correct validation for collection id
	if strings.TrimSpace(req.CollectionId) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokenExecutor := NewTokenExecutor(ctx, k.tokenKeeper)
	nftCollection, found := tokenExecutor.GetNftCollection(collectionCreator, req.CollectionId)

	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	marketplaceIndex := types.GetMarketplaceIndex(marketplaceCreator, req.MarketplaceId)
	index := nftCollection.Index

	collection, found := k.GetMarketplaceCollection(
		ctx,
		marketplaceIndex,
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetMarketplaceCollectionResponse{
		MarketplaceCreator:                      marketplaceCreator.String(),
		MarketplaceId:                           req.MarketplaceId,
		CollectionCreator:                       collection.CollectionCreator,
		CollectionId:                            collection.CollectionId,
		InitiallyNftCollectionOwnerNftsForSale:  collection.InitiallyNftCollectionOwnerNftsForSale,
		InitiallyNftCollectionOwnerNftsMinPrice: collection.InitiallyNftCollectionOwnerNftsMinPrice,
		Cw20ContractAddress:                     collection.Cw20ContractAddress.String(),
		NftsEarningsOnSale:                      collection.NftsEarningsOnSale,
		NftsEarningsOnYieldReward:               collection.NftsEarningsOnYieldReward,
		InitiallyNftsVaultLockPercentage:        collection.InitiallyNftsVaultLockPercentage.String(),
		Creator:                                 collection.Creator.String(),
	}, nil
}

func (k Keeper) AllMarketplaceCollections(goCtx context.Context, req *types.QueryGetAllMarketplaceCollectionsRequest) (*types.QueryGetAllMarketplaceCollectionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	marketplaceCreator, err := sdk.AccAddressFromBech32(req.MarketplaceCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateMarketplaceId(conf.ValidMarketplaceId, req.MarketplaceId)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	marketplaceIndex := types.GetMarketplaceIndex(marketplaceCreator, req.MarketplaceId)

	if !k.HasMarketplace(ctx, marketplaceCreator, marketplaceIndex) {
		return nil, status.Error(codes.InvalidArgument, "marketplace not exists")
	}

	store := ctx.KVStore(k.storeKey)
	marketplaceCollectionsStore := prefix.NewStore(store, types.MarketplaceCollectionStoreKey(marketplaceIndex))

	var marketplaceCollections []*types.MarketplaceCollection
	pageRes, err := query.Paginate(marketplaceCollectionsStore, req.Pagination, func(_ []byte, value []byte) error {
		var nftMeta types.MarketplaceCollection
		if err := k.cdc.Unmarshal(value, &nftMeta); err != nil {
			return err
		}
		marketplaceCollections = append(marketplaceCollections, &nftMeta)
		return nil
	})

	if err != nil {
		return nil, err
	}

	var marketplaceCollectionsRes []*types.QueryGetMarketplaceCollectionResponse

	for _, collection := range marketplaceCollections {
		marketplaceCollectionsRes = append(marketplaceCollectionsRes, &types.QueryGetMarketplaceCollectionResponse{
			MarketplaceCreator:                      marketplaceCreator.String(),
			MarketplaceId:                           req.MarketplaceId,
			CollectionCreator:                       collection.CollectionCreator,
			CollectionId:                            collection.CollectionId,
			InitiallyNftCollectionOwnerNftsForSale:  collection.InitiallyNftCollectionOwnerNftsForSale,
			InitiallyNftCollectionOwnerNftsMinPrice: collection.InitiallyNftCollectionOwnerNftsMinPrice,
			Cw20ContractAddress:                     collection.Cw20ContractAddress.String(),
			NftsEarningsOnSale:                      collection.NftsEarningsOnSale,
			NftsEarningsOnYieldReward:               collection.NftsEarningsOnYieldReward,
			InitiallyNftsVaultLockPercentage:        collection.InitiallyNftsVaultLockPercentage.String(),
			Creator:                                 collection.Creator.String(),
		})
	}

	return &types.QueryGetAllMarketplaceCollectionsResponse{
		MarketplaceId:      req.MarketplaceId,
		MarketplaceCreator: marketplaceCreator.String(),

		Collections: marketplaceCollectionsRes,
		Pagination:  pageRes,
	}, nil
}
