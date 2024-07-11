package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
)

// CreatePrivatePlan defines a method to create a new private plan.
func (k msgServer) CreatePrivatePlan(goCtx context.Context, msg *types.MsgCreatePrivatePlan) (*types.MsgCreatePrivatePlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, errors.Wrap(err, "unauthorized")
	}

	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	plan, err := k.Keeper.CreatePrivatePlan(
		ctx, creatorAddr, msg.Description, msg.RewardAllocations, msg.StartTime, msg.EndTime)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreatePrivatePlanResponse{
		PlanId:             plan.Id,
		FarmingPoolAddress: plan.FarmingPoolAddress,
	}, nil
}

// TerminatePrivatePlan defines a method to terminate a private plan.
func (k msgServer) TerminatePrivatePlan(goCtx context.Context, msg *types.MsgTerminatePrivatePlan) (*types.MsgTerminatePrivatePlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	plan, found := k.GetPlan(ctx, msg.PlanId)
	if !found {
		return nil, errors.Wrapf(errorstypes.ErrNotFound, "plan not found: %d", msg.PlanId)
	}
	if !plan.IsPrivate {
		return nil, errors.Wrap(errorstypes.ErrInvalidRequest, "cannot terminate public plan")
	}
	if plan.TerminationAddress != msg.Creator {
		return nil, errors.Wrapf(
			errorstypes.ErrUnauthorized,
			"plan's termination address must be same with the sender's address")
	}

	if err := k.Keeper.TerminatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return &types.MsgTerminatePrivatePlanResponse{}, nil
}

// Farm defines a method for farming coins.
func (k msgServer) Farm(goCtx context.Context, msg *types.MsgFarm) (*types.MsgFarmResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	farmerAddr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		return nil, err
	}

	withdrawnRewards, err := k.Keeper.GetFarm(ctx, farmerAddr, msg.Coin)
	if err != nil {
		return nil, err
	}

	return &types.MsgFarmResponse{
		WithdrawnRewards: withdrawnRewards,
	}, nil
}

// Unfarm defines a method for un-farming coins.
func (k msgServer) Unfarm(goCtx context.Context, msg *types.MsgUnfarm) (*types.MsgUnfarmResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	farmerAddr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		return nil, err
	}

	withdrawnRewards, err := k.Keeper.Unfarm(ctx, farmerAddr, msg.Coin)
	if err != nil {
		return nil, err
	}

	return &types.MsgUnfarmResponse{
		WithdrawnRewards: withdrawnRewards,
	}, nil
}

// Harvest defines a method for harvesting farming rewards.
func (k msgServer) Harvest(goCtx context.Context, msg *types.MsgHarvest) (*types.MsgHarvestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	farmerAddr, err := sdk.AccAddressFromBech32(msg.Farmer)
	if err != nil {
		return nil, err
	}

	withdrawnRewards, err := k.Keeper.Harvest(ctx, farmerAddr, msg.Denom)
	if err != nil {
		return nil, err
	}

	return &types.MsgHarvestResponse{
		WithdrawnRewards: withdrawnRewards,
	}, nil
}
