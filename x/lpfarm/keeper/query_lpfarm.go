package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
)

func (k Keeper) QueryPlans(c context.Context, req *types.QueryPlansRequest) (*types.QueryPlansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PlanKeyPrefix)
	var plans []types.Plan
	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var plan types.Plan
		if err := k.cdc.Unmarshal(value, &plan); err != nil {
			return err
		}
		plans = append(plans, plan)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryPlansResponse{Plans: plans, Pagination: pageRes}, nil
}

func (k Keeper) QueryPlan(c context.Context, req *types.QueryPlanRequest) (*types.QueryPlanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	plan, found := k.GetPlan(ctx, req.PlanId)
	if !found {
		return nil, status.Error(codes.NotFound, "plan not found")
	}
	return &types.QueryPlanResponse{Plan: plan}, nil
}

func (k Keeper) QueryFarm(c context.Context, req *types.QueryFarmRequest) (*types.QueryFarmResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid denom: %v", err)
	}
	ctx := sdk.UnwrapSDKContext(c)
	farm, found := k.GetFarm(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "farm not found")
	}
	return &types.QueryFarmResponse{Farm: farm}, nil
}

func (k Keeper) QueryPositions(c context.Context, req *types.QueryPositionsRequest) (*types.QueryPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	farmerAddr, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid farmer address: %v", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	keyPrefix := types.GetPositionsByFarmerKeyPrefix(farmerAddr)
	store := prefix.NewStore(storeAdapter, keyPrefix)
	var positions []types.Position
	pageReq, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var position types.Position
		if err := k.cdc.Unmarshal(value, &position); err != nil {
			return err
		}
		positions = append(positions, position)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryPositionsResponse{Positions: positions, Pagination: pageReq}, nil
}

func (k Keeper) QueryPosition(c context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	farmerAddr, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid farmer address: %v", err)
	}
	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid denom: %v", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	position, found := k.GetPosition(ctx, farmerAddr, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "position not found")
	}

	return &types.QueryPositionResponse{Position: position}, nil
}

func (k Keeper) QueryHistoricalRewards(c context.Context, req *types.QueryHistoricalRewardsRequest) (*types.QueryHistoricalRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid denom: %v", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	keyPrefix := types.GetHistoricalRewardsByDenomKeyPrefix(req.Denom)
	store := prefix.NewStore(storeAdapter, keyPrefix)
	var hists []types.HistoricalRewardsResponse
	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		_, period := types.ParseHistoricalRewardsKey(append(keyPrefix, key...))
		var hist types.HistoricalRewards
		if err := k.cdc.Unmarshal(value, &hist); err != nil {
			return err
		}
		hists = append(hists, types.HistoricalRewardsResponse{
			Period:                period,
			CumulativeUnitRewards: hist.CumulativeUnitRewards,
			ReferenceCount:        hist.ReferenceCount,
		})
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryHistoricalRewardsResponse{HistoricalRewards: hists, Pagination: pageRes}, nil
}

func (k Keeper) QueryTotalRewards(c context.Context, req *types.QueryTotalRewardsRequest) (*types.QueryTotalRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	farmerAddr, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid farmer address: %v", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(storeAdapter, types.GetPositionsByFarmerKeyPrefix(farmerAddr))
	defer iter.Close()
	rewards := sdk.DecCoins{}
	for ; iter.Valid(); iter.Next() {
		_, denom := types.ParsePositionKey(iter.Key())
		rewards = rewards.Add(k.GetRewards(ctx, farmerAddr, denom)...)
	}

	return &types.QueryTotalRewardsResponse{Rewards: rewards}, nil
}

func (k Keeper) QueryRewards(c context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	farmerAddr, err := sdk.AccAddressFromBech32(req.Farmer)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid farmer address: %v", err)
	}
	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid denom: %v", err)
	}

	ctx := sdk.UnwrapSDKContext(c)
	_, found := k.GetFarm(ctx, req.Denom)
	if !found {
		return &types.QueryRewardsResponse{Rewards: nil}, nil
	}

	return &types.QueryRewardsResponse{
		Rewards: k.GetRewards(ctx, farmerAddr, req.Denom),
	}, nil
}
