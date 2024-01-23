package keeper

import (
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Rewards(goCtx context.Context, req *types.QueryGetRewardsRequest) (*types.QueryGetRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(req.Provider)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	provider, found := k.GetProvider(ctx, creator.String())

	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	endClaimedSnapshotId := k.GetEndClaimedSnapshotId(ctx, req.PairId)

	pairIdx, found := provider.PairIdToIdx[req.PairId]
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	providerPair := provider.Pairs[pairIdx]

	if providerPair == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	k.CalculateRewards(ctx, creator.String(), req.PairId, provider, &types.ClaimParams{EndClaimedSnapshotId: &endClaimedSnapshotId, IsQuery: true})

	pair, found := k.liquidityKeeper.GetPair(ctx, req.PairId)

	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rewards := sdk.NewCoins()

	for _, reward := range providerPair.OwedRewards {
		if reward.Amount.IsZero() {
			continue
		}

		if reward.Denom == pair.BaseCoinDenom || reward.Denom == pair.QuoteCoinDenom {
			rewards = rewards.Add(sdk.NewCoin(reward.Denom, reward.Amount.TruncateInt()))
			reward.Amount = reward.Amount.Sub(sdk.NewDecFromInt(reward.Amount.TruncateInt()))
		}
	}

	result := []*sdk.Coin{}

	if !rewards.IsZero() {
		for _, reward := range rewards {
			result = append(result, &sdk.Coin{Denom: reward.Denom, Amount: reward.Amount})
		}
	}

	return &types.QueryGetRewardsResponse{PairId: req.PairId, Rewards: result}, nil
}
