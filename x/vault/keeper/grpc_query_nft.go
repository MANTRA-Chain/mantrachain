package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/LimeChain/mantrachain/x/vault/utils"
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

	if strings.TrimSpace(req.NftId) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rewardsController := NewRewardsController(ctx, marketplaceCreator, req.MarketplaceId, collectionCreator, req.CollectionId).
		WithNftId(req.NftId).
		WithKeeper(k)

	err = rewardsController.
		NftStakeMustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	nftStake := rewardsController.getNftStake()

	nftStaked := []*types.QueryGetNftStakeListItemResponse{}

	for _, v := range nftStake.Staked {
		nftStaked = append(nftStaked, &types.QueryGetNftStakeListItemResponse{
			Index:               v.Index,
			Amount:              v.Amount.String(),
			Denom:               v.Denom,
			Cw20ContractAddress: v.Cw20ContractAddress,
			Shares:              v.Shares,
			Validator:           v.Validator,
			Chain:               v.Chain,
			Data:                v.Data,
			Creator:             v.Creator.String(),
			BlockHeight:         v.BlockHeight,
			StakedAt:            v.StakedAt,
			UnstakedAt:          v.UnstakedAt,
			StakedEpoch:         v.StakedEpoch,
		})
	}

	nftBalances := []*types.QueryGetNftStakeBalanceResponse{}

	for _, v := range nftStake.Balances {
		nftBalances = append(nftBalances, &types.QueryGetNftStakeBalanceResponse{
			Amount:             v.Amount.String(),
			Denom:              v.Denom,
			Validator:          v.Validator,
			Chain:              v.Chain,
			LastWithdrawnAt:    v.LastWithdrawnAt,
			LastWithdrawnEpoch: v.LastWithdrawnEpoch,
		})
	}

	nftEarningsOnYieldReward := []*types.VaultEarning{}

	for _, k := range nftStake.NftEarningsOnYieldReward {
		nftEarningsOnYieldReward = append(nftEarningsOnYieldReward, &types.VaultEarning{
			Type:       k.Type,
			Address:    k.Address,
			Percentage: k.Percentage,
		})
	}

	return &types.QueryGetNftStakeResponse{
		MarketplaceCreator:       marketplaceCreator.String(),
		MarketplaceId:            req.MarketplaceId,
		CollectionCreator:        collectionCreator.String(),
		CollectionId:             req.CollectionId,
		NftId:                    req.NftId,
		Creator:                  nftStake.Creator.String(),
		Staked:                   nftStaked,
		Balances:                 nftBalances,
		NftEarningsOnYieldReward: nftEarningsOnYieldReward,
		InitiallyRewardWithdrawn: nftStake.InitiallyRewardWithdrawn,
	}, nil
}

func (k Keeper) NftBalances(c context.Context, req *types.QueryGetNftBalancesRequest) (*types.QueryGetNftBalancesResponse, error) {
	var stakingChain = ""
	var stakingValidator = ""

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	if strings.TrimSpace(req.StakingChain) != "" {
		stakingChain = req.StakingChain
		stakingValidator = req.StakingValidator
	} else {
		stakingChain = ctx.ChainID()
		stakingValidator = params.StakingValidatorAddress
	}

	if strings.TrimSpace(stakingValidator) == "" {
		return nil, status.Error(codes.Unavailable, "staking validator address param not set")
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

	if strings.TrimSpace(req.NftId) == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rewardsController := NewRewardsController(ctx, marketplaceCreator, req.MarketplaceId, collectionCreator, req.CollectionId).
		WithNftId(req.NftId).
		WithKeeper(k)

	err = rewardsController.
		NftStakeMustExist().
		Validate()

	if err != nil {
		return nil, err
	}

	var startAt int64
	var endAt int64
	var epochs []*types.Epoch
	var balances []*sdk.DecCoin = nil
	var intBalances []*sdk.Coin = nil

	staked, err := rewardsController.getStaked(stakingChain, stakingValidator)

	if err != nil {
		return nil, err
	}

	var lastEpochWithdrawn int64 = rewardsController.getLastWithdrawnEpoch(stakingChain, stakingValidator)
	minEpochRewardsStartBH := rewardsController.getMinEpochRewardsStartBH(staked, lastEpochWithdrawn)

	if minEpochRewardsStartBH != types.UndefinedBlockHeight {
		epochs = k.GetNextRewardsEpochsFromPrevEpochId(
			ctx,
			&stakingChain,
			&stakingValidator,
			minEpochRewardsStartBH,
		)
	}

	if len(epochs) > 0 {
		startAt = epochs[0].StartAt
		endAt = epochs[len(epochs)-1].EndAt

		balances = rewardsController.calcNftBalances(epochs, staked)
	}

	prevBalances := rewardsController.getBalancesCoin(stakingChain, stakingValidator)
	prevBalances = append(prevBalances, balances...)

	balances = utils.SumCoins(prevBalances)

	for _, v := range balances {
		intBalances = append(intBalances, &sdk.Coin{
			Denom:  v.Denom,
			Amount: v.Amount.TruncateInt(),
		})
	}

	return &types.QueryGetNftBalancesResponse{
		MarketplaceCreator: marketplaceCreator.String(),
		MarketplaceId:      req.MarketplaceId,
		CollectionCreator:  collectionCreator.String(),
		CollectionId:       req.CollectionId,
		NftId:              req.NftId,
		Balances:           intBalances,
		StartAt:            startAt,
		EndAt:              endAt,
		StakingChain:       stakingChain,
		StakingValidator:   stakingValidator,
	}, nil
}
