package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) NftStake(c context.Context, req *types.QueryGetNftStakeRequest) (*types.QueryGetNftStakeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	marketplaceCreator, err := sdk.AccAddressFromBech32(req.MarketplaceCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	collectionCreator, err := sdk.AccAddressFromBech32(req.CollectionCreator)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if strings.TrimSpace(req.MarketplaceId) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if strings.TrimSpace(req.CollectionId) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if strings.TrimSpace(req.Id) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	marketplaceResolver := NewMarketplaceResolver()
	tokenResolver := NewTokenResolver()

	marketplaceIndex := marketplaceResolver.GetMarketplaceIndex(marketplaceCreator, req.MarketplaceId)
	collectionIndex := tokenResolver.GetCollectionIndex(collectionCreator, req.CollectionId)
	index := tokenResolver.GetNftIndex(collectionCreator, req.CollectionId, req.Id)

	stake, found := k.GetNftStake(
		ctx,
		marketplaceIndex,
		collectionIndex,
		index,
	)

	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	var staked []*types.Stake = []*types.Stake{}

	for _, stake := range stake.Staked {
		staked = append(staked, &types.Stake{
			Amount:     stake.Amount,
			Staked:     stake.Staked,
			Shares:     stake.Shares,
			Validator:  stake.Validator,
			Chain:      stake.Chain,
			StakedAt:   stake.StakedAt,
			UnstakedAt: stake.UnstakedAt,
			Creator:    stake.Creator,
			Owner:      stake.Owner,
		})
	}

	return &types.QueryGetNftStakeResponse{
		MarketplaceCreator: marketplaceCreator.String(),
		MarketplaceId:      req.MarketplaceId,
		CollectionCreator:  collectionCreator.String(),
		CollectionId:       req.CollectionId,
		NftId:              req.Id,
		Staked:             staked,
	}, nil
}
