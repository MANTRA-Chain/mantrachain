package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CollectionSettings(c context.Context, req *types.QueryGetCollectionSettingsRequest) (*types.QueryGetCollectionSettingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	marketplaceCreator, err := sdk.AccAddressFromBech32(req.MarketplaceCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	conf := k.GetParams(ctx)

	err = types.ValidateMarketplaceId(conf.ValidMarketplaceId, req.MarketplaceId, nil)

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
	collection, found := tokenExecutor.GetCollection(collectionCreator, req.CollectionId)

	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	marketplaceIndex := types.GetMarketplaceIndex(marketplaceCreator, req.MarketplaceId)
	index := collection.Index

	settings, found := k.GetCollectionSettings(
		ctx,
		marketplaceIndex,
		index,
	)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetCollectionSettingsResponse{
		MarketplaceCreator:                   marketplaceCreator.String(),
		MarketplaceId:                        req.MarketplaceId,
		CollectionCreator:                    collectionCreator.String(),
		CollectionId:                         req.CollectionId,
		InitiallyCollectionOwnerNftsForSale:  settings.InitiallyCollectionOwnerNftsForSale,
		InitiallyCollectionOwnerNftsMinPrice: settings.InitiallyCollectionOwnerNftsMinPrice.String(),
		NftsEarningsOnSale:                   settings.NftsEarningsOnSale,
		NftsEarningsOnYieldReward:            settings.NftsEarningsOnYieldReward,
		InitiallyNftsVaultLockPercentage:     settings.InitiallyNftsVaultLockPercentage.String(),
		Creator:                              settings.Creator.String(),
	}, nil
}
