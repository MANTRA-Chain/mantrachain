package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
)

// getNextPairIdWithUpdate increments pair id by one and set it.
func (k Keeper) getNextPairIdWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetLastPairId(ctx) + 1
	k.SetLastPairId(ctx, id)
	return id
}

// getNextOrderIdWithUpdate increments the pair's last order id and returns it.
func (k Keeper) getNextOrderIdWithUpdate(ctx sdk.Context, pair types.Pair) uint64 {
	id := pair.LastOrderId + 1
	pair.LastOrderId = id
	k.SetPair(ctx, pair)
	return id
}

// ValidateMsgCreatePair validates types.MsgCreatePair.
func (k Keeper) ValidateMsgCreatePair(ctx sdk.Context, msg *types.MsgCreatePair) error {
	if _, found := k.GetPairByDenoms(ctx, msg.BaseCoinDenom, msg.QuoteCoinDenom); found {
		return types.ErrPairAlreadyExists
	}
	if msg.SwapFeeRate == nil && msg.PairCreatorSwapFeeRatio != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "swap fee rate must not be nil when pair creator swap fee ratio is not nil")
	}
	if msg.SwapFeeRate != nil && msg.PairCreatorSwapFeeRatio == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pair creator swap fee ratio must not be nil when swap fee rate is not nil")
	}
	return nil
}

// ValidateMsgUpdatePairSwapFee validates types.MsgUpdatePairSwapFee.
func (k Keeper) ValidateMsgUpdatePairSwapFee(ctx sdk.Context, msg *types.MsgUpdatePairSwapFee) error {
	_, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}
	if msg.SwapFeeRate == nil && msg.PairCreatorSwapFeeRatio != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "swap fee rate must not be nil when pair creator swap fee ratio is not nil")
	}
	if msg.SwapFeeRate != nil && msg.PairCreatorSwapFeeRatio == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pair creator swap fee ratio must not be nil when swap fee rate is not nil")
	}
	return nil
}

// CreatePair handles types.MsgCreatePair and creates a pair.
func (k Keeper) CreatePair(ctx sdk.Context, msg *types.MsgCreatePair) (types.Pair, error) {
	if err := k.ValidateMsgCreatePair(ctx, msg); err != nil {
		return types.Pair{}, err
	}

	feeCollector := k.GetFeeCollector(ctx)
	pairCreationFee := k.GetPairCreationFee(ctx)

	// Send the pair creation fee to the fee collector.
	if err := k.bankKeeper.SendCoins(ctx, msg.GetCreator(), feeCollector, pairCreationFee); err != nil {
		return types.Pair{}, sdkerrors.Wrap(err, "insufficient pair creation fee")
	}

	id := k.getNextPairIdWithUpdate(ctx)
	pair := types.NewPair(id, msg.BaseCoinDenom, msg.QuoteCoinDenom, msg.Creator, msg.SwapFeeRate, msg.PairCreatorSwapFeeRatio)
	k.SetPair(ctx, pair)
	k.SetPairIndex(ctx, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
	k.SetPairLookupIndex(ctx, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
	k.SetPairLookupIndex(ctx, pair.QuoteCoinDenom, pair.BaseCoinDenom, pair.Id)

	var swapFeeRate sdk.Dec
	if msg.SwapFeeRate != nil {
		swapFeeRate = *pair.SwapFeeRate
	} else {
		swapFeeRate = k.GetSwapFeeRate(ctx)
	}

	var pairCreatorSwapFeeRatio sdk.Dec
	if msg.PairCreatorSwapFeeRatio != nil {
		pairCreatorSwapFeeRatio = *pair.PairCreatorSwapFeeRatio
	} else {
		pairCreatorSwapFeeRatio = k.GetPairCreatorSwapFeeRatio(ctx)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePair,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyBaseCoinDenom, msg.BaseCoinDenom),
			sdk.NewAttribute(types.AttributeKeyQuoteCoinDenom, msg.QuoteCoinDenom),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(pair.Id, 10)),
			sdk.NewAttribute(types.AttributeKeySwapFeeRate, swapFeeRate.String()),
			sdk.NewAttribute(types.AttributeKeyPairCreatorSwapFeeRatio, pairCreatorSwapFeeRatio.String()),
			sdk.NewAttribute(types.AttributeKeyEscrowAddress, pair.EscrowAddress),
		),
	})

	return pair, nil
}

// UpdatePairSwapFee handles types.MsgUpdatePairSwapFee and creates a pair.
func (k Keeper) UpdatePairSwapFee(ctx sdk.Context, msg *types.MsgUpdatePairSwapFee) (types.Pair, error) {
	if err := k.ValidateMsgUpdatePairSwapFee(ctx, msg); err != nil {
		return types.Pair{}, err
	}

	pair, _ := k.GetPair(ctx, msg.PairId)

	var swapFeeRate sdk.Dec
	if msg.SwapFeeRate != nil {
		swapFeeRate = *msg.SwapFeeRate
	} else {
		swapFeeRate = sdk.Dec{}
	}

	var pairCreatorSwapFeeRatio sdk.Dec
	if msg.PairCreatorSwapFeeRatio != nil {
		pairCreatorSwapFeeRatio = *msg.PairCreatorSwapFeeRatio
	} else {
		pairCreatorSwapFeeRatio = sdk.Dec{}
	}

	pair.SwapFeeRate = msg.SwapFeeRate
	pair.PairCreatorSwapFeeRatio = msg.PairCreatorSwapFeeRatio

	k.SetPair(ctx, pair)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdatePairSwapFee,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(pair.Id, 10)),
			sdk.NewAttribute(types.AttributeKeySwapFeeRate, swapFeeRate.String()),
			sdk.NewAttribute(types.AttributeKeyPairCreatorSwapFeeRatio, pairCreatorSwapFeeRatio.String()),
		),
	})

	return pair, nil
}
