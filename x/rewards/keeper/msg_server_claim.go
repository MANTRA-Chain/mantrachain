package keeper

import (
	"context"
	"math"
	"strconv"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
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
		return nil, errors.Wrap(types.ErrProviderNotFound, "provider not found")
	}

	endClaimedSnapshotId := k.GetEndClaimedSnapshotId(ctx, msg.PairId)

	pairIdx, found := provider.PairIdToIdx[msg.PairId]
	if !found {
		return nil, errors.Wrap(types.ErrProviderPairNotFound, "provider pair not found")
	}

	providerPair := provider.Pairs[pairIdx]

	if providerPair == nil {
		return nil, errors.Wrap(types.ErrProviderPairNotFound, "provider pair not found")
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
			SnapshotId: 0,
		}
	}

	if snapshotStartId.SnapshotId > startClaimedSnapshotId {
		startClaimedSnapshotId = snapshotStartId.SnapshotId
	}

	if conf.MaxClaimedRangeLength > 0 && endClaimedSnapshotId-startClaimedSnapshotId >= conf.MaxClaimedRangeLength {
		endClaimedSnapshotId = startClaimedSnapshotId + conf.MaxClaimedRangeLength - 1
	}

	provider, err = k.CalculateRewards(ctx, msg.PairId, provider, &types.ClaimParams{
		StartClaimedSnapshotId: &startClaimedSnapshotId,
		EndClaimedSnapshotId:   &endClaimedSnapshotId,
		IsQuery:                false,
	})
	if err != nil {
		return nil, err
	}

	pair, found := k.liquidityKeeper.GetPair(ctx, msg.PairId)

	if !found {
		return nil, errors.Wrap(types.ErrPairNotFound, "pair not found")
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
		if err := k.guardKeeper.CheckCanTransferCoins(ctx, creator, rewards); err != nil {
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

	baseAmount := sdkmath.ZeroInt()
	quoteAmount := sdkmath.ZeroInt()

	if ok, _ := rewards.Find(pair.BaseCoinDenom); ok {
		baseAmount = rewards.AmountOf(pair.BaseCoinDenom)
	}
	if ok, _ := rewards.Find(pair.QuoteCoinDenom); ok {
		quoteAmount = rewards.AmountOf(pair.QuoteCoinDenom)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyProvider, provider.Index),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeySnapshotId, strconv.FormatUint(endClaimedSnapshotId, 10)),
			sdk.NewAttribute(types.AttributeKeyBaseDenom, pair.BaseCoinDenom),
			sdk.NewAttribute(types.AttributeKeyBaseAmount, baseAmount.String()),
			sdk.NewAttribute(types.AttributeKeyQuoteDenom, pair.QuoteCoinDenom),
			sdk.NewAttribute(types.AttributeKeyQuoteAmount, quoteAmount.String()),
		),
	})

	return &types.MsgClaimResponse{}, nil
}
