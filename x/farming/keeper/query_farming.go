package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/MANTRA-Finance/mantrachain/x/farming/types"
)

// Plans queries all plans.
func (k Keeper) Plans(c context.Context, req *types.QueryPlansRequest) (*types.QueryPlansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Type != "" && !(req.Type == types.PlanTypePublic.String() || req.Type == types.PlanTypePrivate.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid plan type %s", req.Type)
	}

	if req.FarmingPoolAddress != "" {
		if _, err := sdk.AccAddressFromBech32(req.FarmingPoolAddress); err != nil {
			return nil, err
		}
	}

	if req.TerminationAddress != "" {
		if _, err := sdk.AccAddressFromBech32(req.TerminationAddress); err != nil {
			return nil, err
		}
	}

	if req.StakingCoinDenom != "" {
		if err := sdk.ValidateDenom(req.StakingCoinDenom); err != nil {
			return nil, err
		}
	}

	var terminated bool
	if req.Terminated != "" {
		var err error
		terminated, err = strconv.ParseBool(req.Terminated)
		if err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PlanKeyPrefix)

	var plans []*codectypes.Any
	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		plan, err := k.UnmarshalPlan(value)
		if err != nil {
			return false, err
		}
		planAny, err := types.PackPlan(plan)
		if err != nil {
			return false, err
		}

		if req.Type != "" && plan.GetType().String() != req.Type {
			return false, nil
		}

		if req.FarmingPoolAddress != "" && plan.GetFarmingPoolAddress().String() != req.FarmingPoolAddress {
			return false, nil
		}

		if req.TerminationAddress != "" && plan.GetTerminationAddress().String() != req.TerminationAddress {
			return false, nil
		}

		if req.StakingCoinDenom != "" {
			found := false
			for _, coin := range plan.GetStakingCoinWeights() {
				if coin.Denom == req.StakingCoinDenom {
					found = true
					break
				}
			}
			if !found {
				return false, nil
			}
		}

		if req.Terminated != "" {
			if plan.IsTerminated() != terminated {
				return false, nil
			}
		}

		if accumulate {
			plans = append(plans, planAny)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPlansResponse{Plans: plans, Pagination: pageRes}, nil
}

// Plan queries a specific plan.
func (k Keeper) Plan(c context.Context, req *types.QueryPlanRequest) (*types.QueryPlanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	plan, found := k.GetPlan(ctx, req.PlanId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "plan %d not found", req.PlanId)
	}

	planAny, err := types.PackPlan(plan)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPlanResponse{Plan: planAny}, nil
}

// Position queries farming position for a farmer.
func (k Keeper) Position(c context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	farmerAcc, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, err
	}

	if req.StakingCoinDenom != "" {
		if err := sdk.ValidateDenom(req.StakingCoinDenom); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)

	resp := &types.QueryPositionResponse{
		StakedCoins: sdk.Coins{},
		QueuedCoins: sdk.Coins{},
		Rewards:     sdk.Coins{},
	}
	if req.StakingCoinDenom == "" {
		resp.StakedCoins = k.GetAllStakedCoinsByFarmer(ctx, farmerAcc)
		resp.QueuedCoins = k.GetAllQueuedCoinsByFarmer(ctx, farmerAcc)
		resp.Rewards = k.AllRewards(ctx, farmerAcc).Add(k.AllUnharvestedRewards(ctx, farmerAcc)...)
	} else {
		staking, found := k.GetStaking(ctx, req.StakingCoinDenom, farmerAcc)
		if found {
			resp.StakedCoins = resp.StakedCoins.Add(sdk.NewCoin(req.StakingCoinDenom, staking.Amount))
		}
		queuedStakingAmt := k.GetAllQueuedStakingAmountByFarmerAndDenom(ctx, farmerAcc, req.StakingCoinDenom)
		if queuedStakingAmt.IsPositive() {
			resp.QueuedCoins = resp.QueuedCoins.Add(sdk.NewCoin(req.StakingCoinDenom, queuedStakingAmt))
		}
		unharvested, _ := k.GetUnharvestedRewards(ctx, farmerAcc, req.StakingCoinDenom)
		resp.Rewards = k.GetRewards(ctx, farmerAcc, req.StakingCoinDenom).Add(unharvested.Rewards...)
	}

	return resp, nil
}

// Stakings queries all stakings of the farmer.
func (k Keeper) Stakings(c context.Context, req *types.QueryStakingsRequest) (*types.QueryStakingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	farmerAcc, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	keyPrefix := types.GetStakingsByFarmerPrefix(farmerAcc)
	store := prefix.NewStore(storeAdapter, keyPrefix)
	var stakings []types.StakingResponse

	pageRes, _ := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		_, stakingCoinDenom := types.ParseStakingIndexKey(append(keyPrefix, key...))

		if req.StakingCoinDenom != "" && stakingCoinDenom != req.StakingCoinDenom {
			return false, nil
		}

		staking, _ := k.GetStaking(ctx, stakingCoinDenom, farmerAcc)

		if accumulate {
			stakings = append(stakings, types.StakingResponse{
				StakingCoinDenom: stakingCoinDenom,
				Amount:           staking.Amount,
				StartingEpoch:    staking.StartingEpoch,
			})
		}

		return true, nil
	})

	return &types.QueryStakingsResponse{Stakings: stakings, Pagination: pageRes}, nil
}

// QueuedStakings queries all queued stakings of the farmer.
func (k Keeper) QueuedStakings(c context.Context, req *types.QueryQueuedStakingsRequest) (*types.QueryQueuedStakingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	farmerAcc, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var keyPrefix []byte
	if req.StakingCoinDenom != "" {
		keyPrefix = types.GetQueuedStakingsByFarmerAndDenomPrefix(farmerAcc, req.StakingCoinDenom)
	} else {
		keyPrefix = types.GetQueuedStakingsByFarmerPrefix(farmerAcc)
	}
	store := prefix.NewStore(storeAdapter, keyPrefix)
	var queuedStakings []types.QueuedStakingResponse

	pageRes, _ := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		_, stakingCoinDenom, endTime := types.ParseQueuedStakingIndexKey(append(keyPrefix, key...))

		queuedStaking, _ := k.GetQueuedStaking(ctx, endTime, stakingCoinDenom, farmerAcc)

		if accumulate {
			queuedStakings = append(queuedStakings, types.QueuedStakingResponse{
				StakingCoinDenom: stakingCoinDenom,
				Amount:           queuedStaking.Amount,
				EndTime:          endTime,
			})
		}

		return true, nil
	})

	return &types.QueryQueuedStakingsResponse{QueuedStakings: queuedStakings, Pagination: pageRes}, nil
}

// TotalStakings queries total staking coin amount for a specific staking coin denom.
func (k Keeper) TotalStakings(c context.Context, req *types.QueryTotalStakingsRequest) (*types.QueryTotalStakingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := sdk.ValidateDenom(req.StakingCoinDenom); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	totalStakings, found := k.GetTotalStakings(ctx, req.StakingCoinDenom)
	if !found {
		totalStakings.Amount = math.ZeroInt()
	}

	return &types.QueryTotalStakingsResponse{
		Amount: totalStakings.Amount,
	}, nil
}

// Rewards queries all accumulated rewards for a farmer.
func (k Keeper) Rewards(c context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	farmerAcc, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, err
	}

	if req.StakingCoinDenom != "" {
		if err := sdk.ValidateDenom(req.StakingCoinDenom); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var keyPrefix []byte
	if req.StakingCoinDenom != "" {
		keyPrefix = types.GetStakingIndexKey(farmerAcc, req.StakingCoinDenom)
	} else {
		keyPrefix = types.GetStakingsByFarmerPrefix(farmerAcc)
	}
	store := prefix.NewStore(storeAdapter, keyPrefix)
	var rewards []types.RewardsResponse

	pageRes, _ := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		_, stakingCoinDenom := types.ParseStakingIndexKey(append(keyPrefix, key...))

		r := k.GetRewards(ctx, farmerAcc, stakingCoinDenom)

		if accumulate {
			rewards = append(rewards, types.RewardsResponse{
				StakingCoinDenom: stakingCoinDenom,
				Rewards:          r,
			})
		}

		return true, nil
	})

	return &types.QueryRewardsResponse{Rewards: rewards, Pagination: pageRes}, nil
}

// UnharvestedRewards queries all unharvested rewards for the farmer.
func (k Keeper) UnharvestedRewards(c context.Context, req *types.QueryUnharvestedRewardsRequest) (*types.QueryUnharvestedRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	farmerAcc, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var keyPrefix []byte
	if req.StakingCoinDenom != "" {
		keyPrefix = types.GetUnharvestedRewardsKey(farmerAcc, req.StakingCoinDenom)
	} else {
		keyPrefix = types.GetUnharvestedRewardsPrefix(farmerAcc)
	}
	store := prefix.NewStore(storeAdapter, keyPrefix)
	var unharvestedRewards []types.UnharvestedRewardsResponse

	pageRes, _ := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		_, stakingCoinDenom := types.ParseUnharvestedRewardsKey(append(keyPrefix, key...))

		unharvested, _ := k.GetUnharvestedRewards(ctx, farmerAcc, stakingCoinDenom)

		if accumulate {
			unharvestedRewards = append(unharvestedRewards, types.UnharvestedRewardsResponse{
				StakingCoinDenom: stakingCoinDenom,
				Rewards:          unharvested.Rewards,
			})
		}

		return true, nil
	})

	return &types.QueryUnharvestedRewardsResponse{UnharvestedRewards: unharvestedRewards, Pagination: pageRes}, nil
}

// CurrentEpochDays queries current epoch days.
func (k Keeper) CurrentEpochDays(c context.Context, req *types.QueryCurrentEpochDaysRequest) (*types.QueryCurrentEpochDaysResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	currentEpochDays := k.GetCurrentEpochDays(ctx)

	return &types.QueryCurrentEpochDaysResponse{CurrentEpochDays: currentEpochDays}, nil
}

// HistoricalRewards queries HistoricalRewards records for a staking coin denom.
func (k Keeper) HistoricalRewards(c context.Context, req *types.QueryHistoricalRewardsRequest) (*types.QueryHistoricalRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := sdk.ValidateDenom(req.StakingCoinDenom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid staking coin denom: %v", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	keyPrefix := types.GetHistoricalRewardsPrefix(req.StakingCoinDenom)
	store := prefix.NewStore(storeAdapter, keyPrefix)
	var historicalRewards []types.HistoricalRewardsResponse

	pageRes, _ := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		_, epoch := types.ParseHistoricalRewardsKey(append(keyPrefix, key...))

		var rewards types.HistoricalRewards
		k.cdc.MustUnmarshal(value, &rewards)

		if accumulate {
			historicalRewards = append(historicalRewards, types.HistoricalRewardsResponse{
				Epoch:                 epoch,
				CumulativeUnitRewards: rewards.CumulativeUnitRewards,
			})
		}

		return true, nil
	})

	return &types.QueryHistoricalRewardsResponse{HistoricalRewards: historicalRewards, Pagination: pageRes}, nil

}
