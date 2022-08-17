package keeper

import (
	"context"

	"github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LastEpochs(c context.Context, req *types.QueryGetLastEpochsRequest) (*types.QueryGetLastEpochsResponse, error) {
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

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, ctx.ChainID(), params.StakingValidatorAddress, params.StakingValidatorDenom)

	if !found {
		return nil, status.Error(codes.NotFound, "last epoch block not found")
	}

	lastEpoch, found := k.GetEpoch(ctx, ctx.ChainID(), params.StakingValidatorAddress, params.StakingValidatorDenom, lastEpochBlock.BlockHeight)

	if !found {
		return nil, status.Error(codes.NotFound, "last epoch not found")
	}

	var epochs []*types.Epoch = []*types.Epoch{&lastEpoch}

	if lastEpoch.PrevEpochBlock != types.UndefinedBlockHeight {
		prevEpoch, found := k.GetEpoch(ctx, ctx.ChainID(), params.StakingValidatorAddress, params.StakingValidatorDenom, lastEpoch.PrevEpochBlock)

		if !found {
			return nil, status.Error(codes.NotFound, "prev epoch not found")
		}

		epochs = append(epochs, &prevEpoch)
	}

	var epochsRes []*types.QueryGetEpochResponse

	for _, epoch := range epochs {
		epochsRes = append(epochsRes, &types.QueryGetEpochResponse{
			BlockStart: epoch.BlockStart,
			BlockEnd:   epoch.BlockEnd,
			Staked:     epoch.Staked.String(),
			Rewards: &sdk.Coin{
				Denom:  epoch.Rewards.Denom,
				Amount: epoch.Rewards.Amount,
			},
			PrevEpochBlock: epoch.PrevEpochBlock,
			NextEpochBlock: epoch.NextEpochBlock,
			StartAt:        epoch.StartAt,
			EndAt:          epoch.EndAt,
		})
	}

	return &types.QueryGetLastEpochsResponse{
		Epochs: epochsRes,
	}, nil
}
