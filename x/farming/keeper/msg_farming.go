package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/farming/types"
)

// CreateFixedAmountPlan defines a method for creating fixed amount farming plan.
func (k msgServer) CreateFixedAmountPlan(goCtx context.Context, msg *types.MsgCreateFixedAmountPlan) (*types.MsgCreateFixedAmountPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, errors.Wrap(err, "unauthorized")
	}

	poolAcc, err := k.DerivePrivatePlanFarmingPoolAcc(ctx, msg.Name)
	if err != nil {
		return nil, err
	}

	if _, err := k.Keeper.CreateFixedAmountPlan(ctx, msg, poolAcc, msg.GetAccCreator(), types.PlanTypePrivate); err != nil {
		return nil, err
	}

	return &types.MsgCreateFixedAmountPlanResponse{}, nil
}

// CreateRatioPlan defines a method for creating ratio farming plan.
func (k msgServer) CreateRatioPlan(goCtx context.Context, msg *types.MsgCreateRatioPlan) (*types.MsgCreateRatioPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !EnableRatioPlan {
		return nil, types.ErrRatioPlanDisabled
	}

	if err := k.Keeper.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, errors.Wrap(err, "unauthorized")
	}

	poolAcc, err := k.DerivePrivatePlanFarmingPoolAcc(ctx, msg.Name)
	if err != nil {
		return nil, err
	}

	if _, err := k.Keeper.CreateRatioPlan(ctx, msg, poolAcc, msg.GetAccCreator(), types.PlanTypePrivate); err != nil {
		return nil, err
	}

	plans := k.GetPlans(ctx)
	if err := types.ValidateTotalEpochRatio(plans); err != nil {
		return nil, err
	}

	return &types.MsgCreateRatioPlanResponse{}, nil
}

// Stake defines a method for staking coins to the farming plan.
func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.Stake(ctx, msg.GetAccFarmer(), msg.StakingCoins); err != nil {
		return nil, err
	}

	return &types.MsgStakeResponse{}, nil
}

// Unstake defines a method for unstaking coins from the farming plan.
func (k msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.Unstake(ctx, msg.GetAccFarmer(), msg.UnstakingCoins); err != nil {
		return nil, err
	}

	return &types.MsgUnstakeResponse{}, nil
}

// Harvest defines a method for claiming farming rewards from the farming plan.
func (k msgServer) Harvest(goCtx context.Context, msg *types.MsgHarvest) (*types.MsgHarvestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.Harvest(ctx, msg.GetAccFarmer(), msg.StakingCoinDenoms); err != nil {
		return nil, err
	}

	return &types.MsgHarvestResponse{}, nil
}

func (k msgServer) RemovePlan(goCtx context.Context, msg *types.MsgRemovePlan) (*types.MsgRemovePlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.RemovePlan(ctx, msg.GetAccCreator(), msg.PlanId); err != nil {
		return nil, err
	}

	return &types.MsgRemovePlanResponse{}, nil
}

// AdvanceEpoch defines a method for advancing epoch by one, just for testing purpose
// and shouldn't be used in real world.
func (k msgServer) AdvanceEpoch(goCtx context.Context, msg *types.MsgAdvanceEpoch) (*types.MsgAdvanceEpochResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.guardKeeper.CheckIsAdmin(ctx, msg.Requester); err != nil {
		return nil, errors.Wrap(err, "unauthorized")
	}

	if EnableAdvanceEpoch {
		currentEpochDays := k.GetCurrentEpochDays(ctx)
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Duration(currentEpochDays) * types.Day))
		if err := k.Keeper.AdvanceEpoch(ctx); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("AdvanceEpoch is disabled")
	}

	return &types.MsgAdvanceEpochResponse{}, nil
}
