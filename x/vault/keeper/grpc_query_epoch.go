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

	lastEpochBlock, found := k.GetLastEpochBlock(ctx, ctx.ChainID(), params.StakingValidatorAddress)

	if !found {
		return nil, status.Error(codes.NotFound, "last epoch block not found")
	}

	lastEpoch, found := k.GetEpoch(ctx, ctx.ChainID(), params.StakingValidatorAddress, lastEpochBlock.BlockHeight)

	if !found {
		return nil, status.Error(codes.NotFound, "last epoch not found")
	}

	var epochs []*types.Epoch = []*types.Epoch{&lastEpoch}

	if lastEpoch.PrevEpochBlock != types.UndefinedBlockHeight {
		prevEpoch, found := k.GetEpoch(ctx, ctx.ChainID(), params.StakingValidatorAddress, lastEpoch.PrevEpochBlock)

		if !found {
			return nil, status.Error(codes.NotFound, "prev epoch not found")
		}

		epochs = append(epochs, &prevEpoch)
	}

	var epochsRes []*types.QueryGetEpochResponse

	for _, epoch := range epochs {
		rewards := make([]*sdk.Coin, len(epoch.Rewards))
		for _, reward := range epoch.Rewards {
			rewards = append(rewards, &sdk.Coin{
				Denom:  reward.Denom,
				Amount: reward.Amount,
			})
		}

		epochsRes = append(epochsRes, &types.QueryGetEpochResponse{
			BlockStart:     epoch.BlockStart,
			BlockEnd:       epoch.BlockEnd,
			Staked:         epoch.Staked.String(),
			Rewards:        rewards,
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
