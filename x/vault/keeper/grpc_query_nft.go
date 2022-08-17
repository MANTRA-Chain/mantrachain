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

	rewardsController := NewRewardsController(ctx, marketplaceCreator, req.MarketplaceId, collectionCreator, req.CollectionId).
		WithNftId(req.Id).
		WithKeeper(k)

	err = rewardsController.
		NftStakeMustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	var staked []*types.Stake = []*types.Stake{}
	for _, stake := range rewardsController.getNftStake().Staked {
		staked = append(staked, &stake)
	}

	return &types.QueryGetNftStakeResponse{
		MarketplaceCreator: marketplaceCreator.String(),
		MarketplaceId:      req.MarketplaceId,
		CollectionCreator:  collectionCreator.String(),
		CollectionId:       req.CollectionId,
		NftId:              req.Id,
		Creator:            rewardsController.getNftStake().Creator.String(),
		Staked:             staked,
	}, nil
}

func (k Keeper) NftBalance(c context.Context, req *types.QueryGetNftBalanceRequest) (*types.QueryGetNftBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetParams(ctx)

	if params.StakingValidatorAddress == "" {
		return nil, status.Error(codes.Unavailable, "staking validator address param not set")
	}

	if params.StakingValidatorDenom == "" {
		return nil, status.Error(codes.Unavailable, "staking validator denom param not set")
	}

	if params.MinRewardWithdrawAmount == 0 {
		return nil, status.Error(codes.Unavailable, "min reward withdraw amount param not set")
	}

	if ctx.ChainID() == "" {
		return nil, status.Error(codes.Unavailable, "chain id not set yet")
	}

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

	rewardsController := NewRewardsController(ctx, marketplaceCreator, req.MarketplaceId, collectionCreator, req.CollectionId).
		WithNftId(req.Id).
		WithKeeper(k).
		WithConfiguration(params)

	err = rewardsController.
		NftStakeMustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	var startAt int64
	var endAt int64
	var epochs []*types.Epoch
	var balance sdk.DecCoin = sdk.NewDecCoin(params.StakingValidatorDenom, sdk.Int(sdk.NewDec(0)))

	staked, err := rewardsController.getNativeStaked()

	if err != nil {
		return nil, err
	}

	minEpochRewardsStartBH := rewardsController.getMinEpochRewardsStartBH(staked)

	if minEpochRewardsStartBH != types.UndefinedBlockHeight {
		epochs = k.GetNextRewardsEpochsFromPrevEpochId(
			ctx,
			ctx.ChainID(),
			params.StakingValidatorAddress,
			params.StakingValidatorDenom,
			minEpochRewardsStartBH,
		)
	}

	if len(epochs) > 0 {
		startAt = epochs[0].StartAt
		endAt = epochs[len(epochs)-1].EndAt

		balance = rewardsController.calcNftBalance(epochs, staked, params.StakingValidatorDenom)
	}

	prevBalance := rewardsController.getPrevBalanceNative()
	balance.Amount = balance.Amount.Add(prevBalance.Amount)

	return &types.QueryGetNftBalanceResponse{
		MarketplaceCreator: marketplaceCreator.String(),
		MarketplaceId:      req.MarketplaceId,
		CollectionCreator:  collectionCreator.String(),
		CollectionId:       req.CollectionId,
		NftId:              req.Id,
		Balance:            &balance,
		StartAt:            startAt,
		EndAt:              endAt,
	}, nil
}
