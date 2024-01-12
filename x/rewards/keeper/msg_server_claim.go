package keeper

import (
	"context"
	"math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/AumegaChain/aumega/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.GetCreator())
	if err != nil {
		return nil, err
	}

	provider, found := k.GetProvider(ctx, msg.GetCreator())

	if !found {
		return nil, sdkerrors.Wrap(types.ErrProviderNotFound, "provider not found")
	}

	endClaimedSnapshotId := k.GetEndClaimedSnapshotId(ctx, msg.PairId)

	pairIdx, found := provider.PairIdToIdx[msg.PairId]
	if !found {
		return nil, sdkerrors.Wrap(types.ErrProviderPairNotFound, "provider pair not found")
	}

	providerPair := provider.Pairs[pairIdx]

	if providerPair == nil {
		return nil, sdkerrors.Wrap(types.ErrProviderPairNotFound, "provider pair not found")
	}

	startClaimedSnapshotId := providerPair.LastClaimedSnapshotId

	if startClaimedSnapshotId == uint64(math.MaxUint64) {
		startClaimedSnapshotId = 0
	} else {
		startClaimedSnapshotId++
	}

	snapshotStartId, found := k.GetSnapshotStartId(ctx, msg.PairId)

	if !found {
		snapshotStartId = types.SnapshotStartId{
			PairId:     msg.PairId,
			SnapshotId: 0,
		}
	}

	if snapshotStartId.SnapshotId > startClaimedSnapshotId {
		startClaimedSnapshotId = snapshotStartId.SnapshotId
	}

	if conf.MaxClaimedRangeLength > 0 && endClaimedSnapshotId-startClaimedSnapshotId >= conf.MaxClaimedRangeLength {
		endClaimedSnapshotId = startClaimedSnapshotId + conf.MaxClaimedRangeLength - 1
	}

	provider = k.CalculateRewards(ctx, msg.PairId, provider, &types.ClaimParams{
		StartClaimedSnapshotId: &startClaimedSnapshotId,
		EndClaimedSnapshotId:   &endClaimedSnapshotId,
		IsQuery:                false,
	})

	pair, found := k.liquidityKeeper.GetPair(ctx, msg.PairId)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrPairNotFound, "pair not found")
	}

	rewards := sdk.NewCoins()

	for _, reward := range providerPair.OwedRewards {
		if reward.Amount.IsZero() {
			continue
		}

		if reward.Denom == pair.BaseCoinDenom || reward.Denom == pair.QuoteCoinDenom {
			rewards = rewards.Add(sdk.NewCoin(reward.Denom, reward.Amount.TruncateInt()))
		}
	}

	if !rewards.IsZero() {
		if err := k.gk.CheckCanTransferCoins(ctx, creator, rewards); err != nil {
			return nil, err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, rewards)
		if err != nil {
			return nil, err
		}
	}

	// Update the provider pair
	providerPair.OwedRewards = providerPair.OwedRewards.Sub(sdk.NewDecCoinsFromCoins(rewards...))

	// Update the provider
	k.SetProvider(ctx, provider)

	return &types.MsgClaimResponse{}, nil
}
