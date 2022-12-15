package keeper

import (
	"context"
	"strings"

	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LastEpochs(c context.Context, req *types.QueryGetLastEpochsRequest) (*types.QueryGetLastEpochsResponse, error) {
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

	if params.StakingValidatorAddress == "" {
		return nil, status.Error(codes.Unavailable, "staking validator address param not set")
	}

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, stakingChain, stakingValidator)

	if !found {
		return nil, status.Error(codes.NotFound, "last epoch block not found")
	}

	lastEpoch, found := k.GetEpoch(ctx, stakingChain, stakingValidator, lastEpochBlock.BlockHeight)

	if !found {
		return nil, status.Error(codes.NotFound, "last epoch not found")
	}

	var epochs []*types.Epoch = []*types.Epoch{&lastEpoch}

	if lastEpoch.PrevEpochBlock != types.UndefinedBlockHeight {
		prevEpoch, found := k.GetEpoch(ctx, stakingChain, stakingValidator, lastEpoch.PrevEpochBlock)

		if !found {
			return nil, status.Error(codes.NotFound, "prev epoch not found")
		}

		epochs = append(epochs, &prevEpoch)
	}

	var epochsRes []*types.QueryGetEpochsResponse

	for _, epoch := range epochs {
		var rewards []*sdk.Coin = nil

		for _, reward := range epoch.Rewards {
			rewards = append(rewards, &sdk.Coin{
				Denom:  reward.Denom,
				Amount: reward.Amount,
			})
		}

		epochsRes = append(epochsRes, &types.QueryGetEpochsResponse{
			BlockStart:       epoch.BlockStart,
			BlockEnd:         epoch.BlockEnd,
			Staked:           epoch.Staked.String(),
			Rewards:          rewards,
			PrevEpochBlock:   epoch.PrevEpochBlock,
			NextEpochBlock:   epoch.NextEpochBlock,
			StartAt:          epoch.StartAt,
			EndAt:            epoch.EndAt,
			StakingChain:     epoch.StakingChain,
			StakingValidator: epoch.StakingValidator,
		})
	}

	return &types.QueryGetLastEpochsResponse{
		Epochs: epochsRes,
	}, nil
}

func (k Keeper) LastEpochBlock(c context.Context, req *types.QueryGetLastEpochBlockRequest) (*types.QueryGetLastEpochBlockResponse, error) {
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

	if params.StakingValidatorAddress == "" {
		return nil, status.Error(codes.Unavailable, "staking validator address param not set")
	}

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, stakingChain, stakingValidator)

	if !found {
		return nil, status.Error(codes.NotFound, "last epoch block not found")
	}

	return &types.QueryGetLastEpochBlockResponse{
		Creator:          lastEpochBlock.Creator,
		BlockHeight:      lastEpochBlock.BlockHeight,
		StakingChain:     lastEpochBlock.StakingChain,
		StakingValidator: lastEpochBlock.StakingValidator,
	}, nil
}
