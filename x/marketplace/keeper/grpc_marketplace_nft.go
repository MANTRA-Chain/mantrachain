package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MarketplaceNft(c context.Context, req *types.QueryGetMarketplaceNftRequest) (*types.QueryGetMarketplaceNftResponse, error) {
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

	if strings.TrimSpace(req.CollectionId) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if strings.TrimSpace(req.NftId) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokenExecutor := NewTokenExecutor(ctx, k.tokenKeeper)
	nft, found := tokenExecutor.GetNft(collectionCreator, req.CollectionId, req.NftId)

	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	marketplaceIndex := types.GetMarketplaceIndex(marketplaceCreator, req.MarketplaceId)
	collectionIndex := nft.CollectionIndex
	index := nft.Index

	var forSale bool
	var initiallySold bool
	var minPrice sdk.Coin
	var cw20ContractAddress sdk.AccAddress

	marketplaceNft, found := k.GetMarketplaceNft(
		ctx,
		marketplaceIndex,
		collectionIndex,
		index,
	)

	if !found {
		marketplaceCollection, found := k.GetMarketplaceCollection(
			ctx,
			marketplaceIndex,
			collectionIndex,
		)

		if !found {
			return nil, status.Error(codes.InvalidArgument, "invalid request")
		}

		forSale = marketplaceCollection.InitiallyNftCollectionOwnerNftsForSale
		initiallySold = false
		minPrice = *marketplaceCollection.InitiallyNftCollectionOwnerNftsMinPrice
		cw20ContractAddress = marketplaceCollection.Cw20ContractAddress
	} else {
		forSale = marketplaceNft.ForSale
		initiallySold = marketplaceNft.InitiallySold
		minPrice = *marketplaceNft.MinPrice
		cw20ContractAddress = marketplaceNft.Cw20ContractAddress
	}

	return &types.QueryGetMarketplaceNftResponse{
		MarketplaceCreator:  marketplaceCreator.String(),
		MarketplaceId:       req.MarketplaceId,
		CollectionCreator:   nft.CollectionCreator.String(),
		CollectionId:        nft.CollectionId,
		NftId:               nft.Id,
		ForSale:             forSale,
		InitiallySold:       initiallySold,
		MinPrice:            &minPrice,
		Cw20ContractAddress: cw20ContractAddress.String(),
		Creator:             marketplaceNft.Creator.String(),
	}, nil
}
