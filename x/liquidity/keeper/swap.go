package keeper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MANTRA-Finance/mantrachain/x/liquidity/amm"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
)

func CalculatePairCreatorSwapFeeAmount(ctx sdk.Context, pairCreatorSwapFeeRatio sdkmath.LegacyDec, accumulatedSwapFee math.Int) math.Int {
	return math.LegacyNewDecFromInt(accumulatedSwapFee).MulTruncate(pairCreatorSwapFeeRatio).TruncateInt()
}

func CalculateSwapFeeAmount(ctx sdk.Context, swapFeeRate sdkmath.LegacyDec, calculatedOfferCoinAmt math.Int) math.Int {
	return math.LegacyNewDecFromInt(calculatedOfferCoinAmt).MulTruncate(swapFeeRate).TruncateInt()
}

func (k Keeper) PriceLimits(ctx sdk.Context, lastPrice sdkmath.LegacyDec) (lowest, highest sdkmath.LegacyDec) {
	return types.PriceLimits(lastPrice, k.GetMaxPriceLimitRatio(ctx), int(k.GetTickPrecision(ctx)))
}

// ValidateMsgLimitOrder validates types.MsgLimitOrder with state and returns
// calculated offer coin and price that is fit into ticks.
func (k Keeper) ValidateMsgLimitOrder(ctx sdk.Context, msg *types.MsgLimitOrder, pair types.Pair) (amount math.Int, offerCoin sdk.Coin, swapFeeCoin sdk.Coin, price sdkmath.LegacyDec, err error) {
	if err := k.guardKeeper.CheckCanTransferCoins(ctx, msg.GetAccOrderer(), sdk.Coins{sdk.NewCoin(msg.OfferCoin.Denom, math.ZeroInt()), sdk.NewCoin(msg.DemandCoinDenom, math.ZeroInt())}); err != nil {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, err
	}

	spendable := k.bankKeeper.SpendableCoins(ctx, msg.GetAccOrderer())
	if spendableAmt := spendable.AmountOf(msg.OfferCoin.Denom); spendableAmt.LT(msg.OfferCoin.Amount) {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(msg.OfferCoin.Denom, spendableAmt), msg.OfferCoin)
	}

	amount = msg.Amount
	tickPrec := k.GetTickPrecision(ctx)
	maxOrderLifespan := k.GetMaxOrderLifespan(ctx)

	var swapFeeRate sdkmath.LegacyDec
	if pair.SwapFeeRate != nil {
		swapFeeRate = *pair.SwapFeeRate
	} else {
		swapFeeRate = k.GetSwapFeeRate(ctx)
	}

	if msg.OrderLifespan > maxOrderLifespan {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{},
			errors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, maxOrderLifespan)
	}

	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	var upperPriceLimit, lowerPriceLimit sdkmath.LegacyDec
	if pair.LastPrice != nil {
		lowerPriceLimit, upperPriceLimit = k.PriceLimits(ctx, *pair.LastPrice)
	} else {
		upperPriceLimit = amm.HighestTick(int(tickPrec))
		lowerPriceLimit = amm.LowestTick(int(tickPrec))
	}
	switch {
	case msg.Price.GT(upperPriceLimit):
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(types.ErrPriceOutOfRange, "%s is higher than %s", msg.Price, upperPriceLimit)
	case msg.Price.LT(lowerPriceLimit):
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(types.ErrPriceOutOfRange, "%s is lower than %s", msg.Price, lowerPriceLimit)
	}

	switch msg.Direction {
	case types.OrderDirectionBuy:
		if msg.OfferCoin.Denom != pair.QuoteCoinDenom || msg.DemandCoinDenom != pair.BaseCoinDenom {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{},
				errors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.DemandCoinDenom, msg.OfferCoin.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToDownTick(msg.Price, int(tickPrec))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, price, amount))
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, swapFeeRate, offerCoin.Amount))

		if msg.OfferCoin.IsLT(offerCoin.Add(swapFeeCoin)) {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, offerCoin.Add(swapFeeCoin))
		}
	case types.OrderDirectionSell:
		if msg.OfferCoin.Denom != pair.BaseCoinDenom || msg.DemandCoinDenom != pair.QuoteCoinDenom {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{},
				errors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.OfferCoin.Denom, msg.DemandCoinDenom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToUpTick(msg.Price, int(tickPrec))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amount)
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, swapFeeRate, offerCoin.Amount))

		if swapFeeCoin.IsPositive() {
			offerCoin.Amount = offerCoin.Amount.Sub(swapFeeCoin.Amount)
			amount = offerCoin.Amount
		}

		if msg.OfferCoin.Amount.LT(swapFeeCoin.Amount.Add(offerCoin.Amount)) {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, sdk.NewCoin(msg.OfferCoin.Denom, swapFeeCoin.Amount.Add(offerCoin.Amount)))
		}
	}
	if types.IsTooSmallOrderAmount(amount, price) {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, types.ErrTooSmallOrder
	}

	return amount, offerCoin, swapFeeCoin, price, nil
}

// LimitOrder handles types.MsgLimitOrder and stores types.Order.
func (k Keeper) LimitOrder(ctx sdk.Context, msg *types.MsgLimitOrder) (types.Order, error) {
	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return types.Order{}, errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	amount, offerCoin, swapFeeCoin, price, err := k.ValidateMsgLimitOrder(ctx, msg, pair)
	if err != nil {
		return types.Order{}, err
	}

	msg.Amount = amount
	refundedCoin := msg.OfferCoin.Sub(offerCoin.Add(swapFeeCoin))

	whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
	err = k.bankKeeper.SendCoins(ctx, msg.GetAccOrderer(), pair.GetEscrowAddress(), sdk.NewCoins(offerCoin.Add(swapFeeCoin)))
	k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
	if err != nil {
		return types.Order{}, err
	}

	requestId := k.getNextOrderIdWithUpdate(ctx, pair)
	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	order := types.NewOrderForLimitOrder(msg, requestId, pair, offerCoin, price, expireAt, ctx.BlockHeight(), swapFeeCoin)
	k.SetOrder(ctx, order)
	k.SetOrderIndex(ctx, order)

	ctx.GasMeter().ConsumeGas(k.GetOrderExtraGas(ctx), "OrderExtraGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLimitOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderDirection, msg.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, offerCoin.String()),
			sdk.NewAttribute(types.AttributeKeyDemandCoinDenom, msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeKeyPrice, price.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOrderId, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyBatchId, strconv.FormatUint(order.BatchId, 10)),
			sdk.NewAttribute(types.AttributeKeyExpireAt, order.ExpireAt.Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyRefundedCoins, refundedCoin.String()),
			sdk.NewAttribute(types.AttributeKeySwapFeeCoins, swapFeeCoin.String()),
		),
	})

	return order, nil
}

// ValidateMsgMarketOrder validates types.MsgMarketOrder with state and returns
// calculated offer coin and price.
func (k Keeper) ValidateMsgMarketOrder(ctx sdk.Context, msg *types.MsgMarketOrder, pair types.Pair) (amount math.Int, offerCoin sdk.Coin, swapFeeCoin sdk.Coin, price sdkmath.LegacyDec, err error) {
	if err := k.guardKeeper.CheckCanTransferCoins(ctx, msg.GetAccOrderer(), sdk.Coins{sdk.NewCoin(msg.OfferCoin.Denom, math.ZeroInt()), sdk.NewCoin(msg.DemandCoinDenom, math.ZeroInt())}); err != nil {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, err
	}

	spendable := k.bankKeeper.SpendableCoins(ctx, msg.GetAccOrderer())
	if spendableAmt := spendable.AmountOf(msg.OfferCoin.Denom); spendableAmt.LT(msg.OfferCoin.Amount) {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(msg.OfferCoin.Denom, spendableAmt), msg.OfferCoin)
	}

	amount = msg.Amount
	maxOrderLifespan := k.GetMaxOrderLifespan(ctx)
	maxPriceLimitRatio := k.GetMaxPriceLimitRatio(ctx)
	tickPrec := k.GetTickPrecision(ctx)

	var swapFeeRate sdkmath.LegacyDec
	if pair.SwapFeeRate != nil {
		swapFeeRate = *pair.SwapFeeRate
	} else {
		swapFeeRate = k.GetSwapFeeRate(ctx)
	}

	if msg.OrderLifespan > maxOrderLifespan {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{},
			errors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, maxOrderLifespan)
	}

	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	if pair.LastPrice == nil {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, types.ErrNoLastPrice
	}
	lastPrice := *pair.LastPrice

	switch msg.Direction {
	case types.OrderDirectionBuy:
		if msg.OfferCoin.Denom != pair.QuoteCoinDenom || msg.DemandCoinDenom != pair.BaseCoinDenom {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{},
				errors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.DemandCoinDenom, msg.OfferCoin.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToDownTick(lastPrice.Mul(math.LegacyOneDec().Add(maxPriceLimitRatio)), int(tickPrec))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, price, amount))
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, swapFeeRate, offerCoin.Amount))

		if msg.OfferCoin.IsLT(offerCoin.Add(swapFeeCoin)) {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, offerCoin.Add(swapFeeCoin))
		}
	case types.OrderDirectionSell:
		if msg.OfferCoin.Denom != pair.BaseCoinDenom || msg.DemandCoinDenom != pair.QuoteCoinDenom {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{},
				errors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.OfferCoin.Denom, msg.DemandCoinDenom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToUpTick(lastPrice.Mul(math.LegacyOneDec().Sub(maxPriceLimitRatio)), int(tickPrec))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amount)
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, swapFeeRate, offerCoin.Amount))

		if swapFeeCoin.IsPositive() {
			offerCoin.Amount = offerCoin.Amount.Sub(swapFeeCoin.Amount)
			amount = offerCoin.Amount
		}

		if msg.OfferCoin.Amount.LT(swapFeeCoin.Amount.Add(offerCoin.Amount)) {
			return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, sdk.NewCoin(msg.OfferCoin.Denom, swapFeeCoin.Amount.Add(offerCoin.Amount)))
		}
	}
	if types.IsTooSmallOrderAmount(amount, price) {
		return math.NewInt(0), sdk.Coin{}, sdk.Coin{}, sdkmath.LegacyDec{}, types.ErrTooSmallOrder
	}

	return amount, offerCoin, swapFeeCoin, price, nil
}

func (k Keeper) GetSwapAmount(ctx sdk.Context, pairId uint64, demandCoin sdk.Coin) (offerCoin sdk.Coin, price sdkmath.LegacyDec, err error) {
	maxPriceLimitRatio := k.GetMaxPriceLimitRatio(ctx)
	tickPrec := k.GetTickPrecision(ctx)

	pair, found := k.GetPair(ctx, pairId)
	if !found {
		return sdk.Coin{}, sdkmath.LegacyDec{}, errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", pairId)
	}

	var swapFeeRate sdkmath.LegacyDec
	if pair.SwapFeeRate != nil {
		swapFeeRate = *pair.SwapFeeRate
	} else {
		swapFeeRate = k.GetSwapFeeRate(ctx)
	}

	if demandCoin.Denom != pair.QuoteCoinDenom && demandCoin.Denom != pair.BaseCoinDenom {
		return sdk.Coin{}, sdkmath.LegacyDec{},
			errors.Wrapf(types.ErrWrongPair, "pair denom %s not exist, pair %d", demandCoin.Denom, pairId)
	}

	if pair.LastPrice == nil {
		return sdk.Coin{}, sdkmath.LegacyDec{}, types.ErrNoLastPrice
	}
	lastPrice := *pair.LastPrice

	if demandCoin.Denom == pair.BaseCoinDenom { // Buy
		price = amm.PriceToDownTick(lastPrice.Mul(math.LegacyOneDec().Add(maxPriceLimitRatio)), int(tickPrec))
		offerCoin = sdk.NewCoin(pair.QuoteCoinDenom, amm.OfferCoinAmount(amm.Buy, price, demandCoin.Amount))
	} else if demandCoin.Denom == pair.QuoteCoinDenom { // Sell
		price = amm.PriceToUpTick(lastPrice.Mul(math.LegacyOneDec().Sub(maxPriceLimitRatio)), int(tickPrec))
		offerCoin = sdk.NewCoin(pair.BaseCoinDenom, math.LegacyNewDecFromInt(demandCoin.Amount).Quo(price).Ceil().TruncateInt())
	}

	if swapFeeRate.IsPositive() {
		swapFeeCoin := sdk.NewCoin(offerCoin.Denom, CalculateSwapFeeAmount(ctx, swapFeeRate, offerCoin.Amount))
		offerCoin.Amount = offerCoin.Amount.Add(swapFeeCoin.Amount)
	}

	return offerCoin, price, nil
}

// MarketOrder handles types.MsgMarketOrder and stores types.Order.
func (k Keeper) MarketOrder(ctx sdk.Context, msg *types.MsgMarketOrder) (types.Order, error) {
	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return types.Order{}, errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	amount, offerCoin, swapFeeCoin, price, err := k.ValidateMsgMarketOrder(ctx, msg, pair)
	if err != nil {
		return types.Order{}, err
	}

	msg.Amount = amount
	refundedCoin := msg.OfferCoin.Sub(offerCoin.Add(swapFeeCoin))

	whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
	err = k.bankKeeper.SendCoins(ctx, msg.GetAccOrderer(), pair.GetEscrowAddress(), sdk.NewCoins(offerCoin.Add(swapFeeCoin)))
	k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
	if err != nil {
		return types.Order{}, err
	}

	requestId := k.getNextOrderIdWithUpdate(ctx, pair)
	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	order := types.NewOrderForMarketOrder(msg, requestId, pair, offerCoin, price, expireAt, ctx.BlockHeight(), swapFeeCoin)
	k.SetOrder(ctx, order)
	k.SetOrderIndex(ctx, order)

	ctx.GasMeter().ConsumeGas(k.GetOrderExtraGas(ctx), "OrderExtraGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMarketOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderDirection, msg.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, offerCoin.String()),
			sdk.NewAttribute(types.AttributeKeyDemandCoinDenom, msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeKeyPrice, price.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOrderId, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyBatchId, strconv.FormatUint(order.BatchId, 10)),
			sdk.NewAttribute(types.AttributeKeyExpireAt, order.ExpireAt.Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyRefundedCoins, refundedCoin.String()),
			sdk.NewAttribute(types.AttributeKeySwapFeeCoins, swapFeeCoin.String()),
		),
	})

	return order, nil
}

func (k Keeper) MMOrder(ctx sdk.Context, msg *types.MsgMMOrder) (orders []types.Order, err error) {
	tickPrec := int(k.GetTickPrecision(ctx))

	if msg.SellAmount.IsPositive() {
		if !amm.PriceToDownTick(msg.MinSellPrice, tickPrec).Equal(msg.MinSellPrice) {
			return nil, errors.Wrapf(types.ErrPriceNotOnTicks, "min sell price is not on ticks")
		}
		if !amm.PriceToDownTick(msg.MaxSellPrice, tickPrec).Equal(msg.MaxSellPrice) {
			return nil, errors.Wrapf(types.ErrPriceNotOnTicks, "max sell price is not on ticks")
		}
	}
	if msg.BuyAmount.IsPositive() {
		if !amm.PriceToDownTick(msg.MinBuyPrice, tickPrec).Equal(msg.MinBuyPrice) {
			return nil, errors.Wrapf(types.ErrPriceNotOnTicks, "min buy price is not on ticks")
		}
		if !amm.PriceToDownTick(msg.MaxBuyPrice, tickPrec).Equal(msg.MaxBuyPrice) {
			return nil, errors.Wrapf(types.ErrPriceNotOnTicks, "max buy price is not on ticks")
		}
	}

	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	var lowestPrice, highestPrice sdkmath.LegacyDec
	if pair.LastPrice != nil {
		lowestPrice, highestPrice = types.PriceLimits(*pair.LastPrice, k.GetMaxPriceLimitRatio(ctx), tickPrec)
	} else {
		lowestPrice = amm.LowestTick(tickPrec)
		highestPrice = amm.HighestTick(tickPrec)
	}

	if msg.SellAmount.IsPositive() {
		if msg.MinSellPrice.LT(lowestPrice) || msg.MinSellPrice.GT(highestPrice) {
			return nil, errors.Wrapf(types.ErrPriceOutOfRange, "min sell price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
		if msg.MaxSellPrice.LT(lowestPrice) || msg.MaxSellPrice.GT(highestPrice) {
			return nil, errors.Wrapf(types.ErrPriceOutOfRange, "max sell price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
	}
	if msg.BuyAmount.IsPositive() {
		if msg.MinBuyPrice.LT(lowestPrice) || msg.MinBuyPrice.GT(highestPrice) {
			return nil, errors.Wrapf(types.ErrPriceOutOfRange, "min buy price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
		if msg.MaxBuyPrice.LT(lowestPrice) || msg.MaxBuyPrice.GT(highestPrice) {
			return nil, errors.Wrapf(types.ErrPriceOutOfRange, "max buy price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
	}

	maxNumTicks := int(k.GetMaxNumMarketMakingOrderTicks(ctx))

	var buyTicks, sellTicks []types.MMOrderTick
	offerBaseCoin := sdk.NewInt64Coin(pair.BaseCoinDenom, 0)
	offerQuoteCoin := sdk.NewInt64Coin(pair.QuoteCoinDenom, 0)
	if msg.BuyAmount.IsPositive() {
		buyTicks = types.MMOrderTicks(
			types.OrderDirectionBuy, msg.MinBuyPrice, msg.MaxBuyPrice, msg.BuyAmount, maxNumTicks, tickPrec)
		for _, tick := range buyTicks {
			offerQuoteCoin = offerQuoteCoin.AddAmount(tick.OfferCoinAmount)
		}
	}
	if msg.SellAmount.IsPositive() {
		sellTicks = types.MMOrderTicks(
			types.OrderDirectionSell, msg.MinSellPrice, msg.MaxSellPrice, msg.SellAmount, maxNumTicks, tickPrec)
		for _, tick := range sellTicks {
			offerBaseCoin = offerBaseCoin.AddAmount(tick.OfferCoinAmount)
		}
	}

	orderer := msg.GetAccOrderer()
	spendable := k.bankKeeper.SpendableCoins(ctx, orderer)
	if spendableAmt := spendable.AmountOf(pair.BaseCoinDenom); spendableAmt.LT(offerBaseCoin.Amount) {
		return nil, errors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(pair.BaseCoinDenom, spendableAmt), offerBaseCoin)
	}
	if spendableAmt := spendable.AmountOf(pair.QuoteCoinDenom); spendableAmt.LT(offerQuoteCoin.Amount) {
		return nil, errors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(pair.QuoteCoinDenom, spendableAmt), offerQuoteCoin)
	}

	maxOrderLifespan := k.GetMaxOrderLifespan(ctx)
	if msg.OrderLifespan > maxOrderLifespan {
		return nil, errors.Wrapf(
			types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, maxOrderLifespan)
	}

	// First, cancel existing market making orders in the pair from the orderer.
	canceledOrderIds, err := k.cancelMMOrder(ctx, orderer, pair, true)
	if err != nil {
		return nil, err
	}

	whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
	err = k.bankKeeper.SendCoins(ctx, orderer, pair.GetEscrowAddress(), sdk.NewCoins(offerBaseCoin, offerQuoteCoin))
	k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
	if err != nil {
		return nil, err
	}

	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	lastOrderId := pair.LastOrderId

	var orderIds []uint64
	for _, tick := range buyTicks {
		lastOrderId++
		offerCoin := sdk.NewCoin(pair.QuoteCoinDenom, tick.OfferCoinAmount)
		order := types.NewOrder(
			types.OrderTypeMM, lastOrderId, pair, orderer,
			offerCoin, tick.Price, tick.Amount, expireAt, ctx.BlockHeight())
		k.SetOrder(ctx, order)
		k.SetOrderIndex(ctx, order)
		orders = append(orders, order)
		orderIds = append(orderIds, order.Id)
	}
	for _, tick := range sellTicks {
		lastOrderId++
		offerCoin := sdk.NewCoin(pair.BaseCoinDenom, tick.OfferCoinAmount)
		order := types.NewOrder(
			types.OrderTypeMM, lastOrderId, pair, orderer,
			offerCoin, tick.Price, tick.Amount, expireAt, ctx.BlockHeight())
		k.SetOrder(ctx, order)
		k.SetOrderIndex(ctx, order)
		orders = append(orders, order)
		orderIds = append(orderIds, order.Id)
	}

	pair.LastOrderId = lastOrderId
	k.SetPair(ctx, pair)

	k.SetMMOrderIndex(ctx, types.NewMMOrderIndex(orderer, pair.Id, orderIds))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMMOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyBatchId, strconv.FormatUint(pair.CurrentBatchId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderIds, types.FormatUint64s(orderIds)),
			sdk.NewAttribute(types.AttributeKeyCanceledOrderIds, types.FormatUint64s(canceledOrderIds)),
		),
	})
	return
}

// ValidateMsgCancelOrder validates types.MsgCancelOrder and returns the order.
func (k Keeper) ValidateMsgCancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) (order types.Order, err error) {
	var found bool
	order, found = k.GetOrder(ctx, msg.PairId, msg.OrderId)
	if !found {
		return types.Order{},
			errors.Wrapf(sdkerrors.ErrNotFound, "order %d not found in pair %d", msg.OrderId, msg.PairId)
	}
	if msg.Orderer != order.Orderer {
		return types.Order{}, errors.Wrap(sdkerrors.ErrUnauthorized, "mismatching orderer")
	}
	if order.Status == types.OrderStatusCanceled {
		return types.Order{}, types.ErrAlreadyCanceled
	}
	pair, _ := k.GetPair(ctx, msg.PairId)
	if order.BatchId == pair.CurrentBatchId {
		return types.Order{}, types.ErrSameBatch
	}
	return order, nil
}

// CancelOrder handles types.MsgCancelOrder and cancels an order.
func (k Keeper) CancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) error {
	order, err := k.ValidateMsgCancelOrder(ctx, msg)
	if err != nil {
		return err
	}

	if err := k.FinishOrder(ctx, order, types.OrderStatusCanceled); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderId, strconv.FormatUint(msg.OrderId, 10)),
		),
	})

	return nil
}

// CancelAllOrders handles types.MsgCancelAllOrders and cancels all orders.
func (k Keeper) CancelAllOrders(ctx sdk.Context, msg *types.MsgCancelAllOrders) error {
	orderPairCache := map[uint64]types.Pair{} // maps order's pair id to pair, to cache the result
	pairIdSet := map[uint64]struct{}{}        // set of pairs where to cancel orders
	var pairIds []string                      // needed to emit an event
	for _, pairId := range msg.PairIds {
		pair, found := k.GetPair(ctx, pairId)
		if !found { // check if the pair exists
			return errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", pairId)
		}
		pairIdSet[pairId] = struct{}{} // add pair id to the set
		pairIds = append(pairIds, strconv.FormatUint(pairId, 10))
		orderPairCache[pairId] = pair // also cache the pair to use at below
	}

	var canceledOrderIds []string
	if err := k.IterateOrdersByOrderer(ctx, msg.GetAccOrderer(), func(order types.Order) (stop bool, err error) {
		_, ok := pairIdSet[order.PairId] // is the pair included in the pair set?
		if len(pairIdSet) == 0 || ok {   // pair ids not specified(cancel all), or the pair is in the set
			pair, ok := orderPairCache[order.PairId]
			if !ok {
				pair, _ = k.GetPair(ctx, order.PairId)
				orderPairCache[order.PairId] = pair
			}
			if order.Status != types.OrderStatusCanceled && order.BatchId < pair.CurrentBatchId {
				if err := k.FinishOrder(ctx, order, types.OrderStatusCanceled); err != nil {
					return false, err
				}
				canceledOrderIds = append(canceledOrderIds, strconv.FormatUint(order.Id, 10))
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAllOrders,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairIds, strings.Join(pairIds, ",")),
			sdk.NewAttribute(types.AttributeKeyCanceledOrderIds, strings.Join(canceledOrderIds, ",")),
		),
	})

	return nil
}

func (k Keeper) cancelMMOrder(ctx sdk.Context, orderer sdk.AccAddress, pair types.Pair, skipIfNotFound bool) (canceledOrderIds []uint64, err error) {
	index, found := k.GetMMOrderIndex(ctx, orderer, pair.Id)
	if found {
		for _, orderId := range index.OrderIds {
			order, found := k.GetOrder(ctx, pair.Id, orderId)
			if !found {
				// The order has already been deleted from store.
				continue
			}
			if order.BatchId == pair.CurrentBatchId {
				return nil, errors.Wrap(types.ErrSameBatch, "couldn't cancel previously placed orders")
			}
			if order.Status.CanBeCanceled() {
				if err := k.FinishOrder(ctx, order, types.OrderStatusCanceled); err != nil {
					return nil, err
				}
				canceledOrderIds = append(canceledOrderIds, order.Id)
			}
		}
		k.DeleteMMOrderIndex(ctx, index)
	} else if !skipIfNotFound {
		return nil, errors.Wrap(sdkerrors.ErrNotFound, "previous market making orders not found")
	}
	return
}

// CancelMMOrder handles types.MsgCancelMMOrder and cancels previous market making
// orders.
func (k Keeper) CancelMMOrder(ctx sdk.Context, msg *types.MsgCancelMMOrder) (canceledOrderIds []uint64, err error) {
	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	canceledOrderIds, err = k.cancelMMOrder(ctx, msg.GetAccOrderer(), pair, false)
	if err != nil {
		return
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelMMOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(pair.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCanceledOrderIds, types.FormatUint64s(canceledOrderIds)),
		),
	})

	return canceledOrderIds, nil
}

func (k Keeper) ExecuteMatching(ctx sdk.Context, pair types.Pair) error {
	ob := amm.NewOrderBook()

	if err := k.IterateOrdersByPair(ctx, pair.Id, func(order types.Order) (stop bool, err error) {
		switch order.Status {
		case types.OrderStatusNotExecuted,
			types.OrderStatusNotMatched,
			types.OrderStatusPartiallyMatched:
			if order.Status != types.OrderStatusNotExecuted && order.ExpiredAt(ctx.BlockTime()) {
				if err := k.FinishOrder(ctx, order, types.OrderStatusExpired); err != nil {
					return false, err
				}
				return false, nil
			}
			// TODO: add orders only when price is in the range?
			ob.AddOrder(types.NewUserOrder(order))
			if order.Status == types.OrderStatusNotExecuted {
				order.SetStatus(types.OrderStatusNotMatched)
				k.SetOrder(ctx, order)
			}
		case types.OrderStatusCanceled:
		default:
			return false, fmt.Errorf("invalid order status: %s", order.Status)
		}
		return false, nil
	}); err != nil {
		return err
	}

	var pools []*types.PoolOrderer
	_ = k.IteratePoolsByPair(ctx, pair.Id, func(pool types.Pool) (stop bool, err error) {
		if pool.Disabled {
			return false, nil
		}
		rx, ry := k.getPoolBalances(ctx, pool, pair)
		ps := k.GetPoolCoinSupply(ctx, pool)
		ammPool := types.NewPoolOrderer(
			pool.AMMPool(rx.Amount, ry.Amount, ps),
			pool.Id, pool.GetReserveAddress(), pair.BaseCoinDenom, pair.QuoteCoinDenom)
		if ammPool.IsDepleted() {
			k.MarkPoolAsDisabled(ctx, pool)
			return false, nil
		}
		pools = append(pools, ammPool)
		return false, nil
	})

	matchPrice, quoteCoinDiff, matched := k.Match(ctx, ob, pools, pair.LastPrice)
	if matched {
		orders := ob.Orders()
		if err := k.ApplyMatchResult(ctx, pair, orders, quoteCoinDiff); err != nil {
			return err
		}
		pair.LastPrice = &matchPrice
	}

	pair.CurrentBatchId++
	k.SetPair(ctx, pair)

	return nil
}

func (k Keeper) Match(ctx sdk.Context, ob *amm.OrderBook, pools []*types.PoolOrderer, lastPrice *sdkmath.LegacyDec) (matchPrice sdkmath.LegacyDec, quoteCoinDiff math.Int, matched bool) {
	tickPrec := int(k.GetTickPrecision(ctx))
	if lastPrice == nil {
		ov := amm.MultipleOrderViews{ob.MakeView()}
		for _, pool := range pools {
			ov = append(ov, pool)
		}
		var found bool
		matchPrice, found = amm.FindMatchPrice(ov, tickPrec)
		if !found {
			return sdkmath.LegacyDec{}, math.Int{}, false
		}
		for _, pool := range pools {
			buyAmt := pool.BuyAmountOver(matchPrice, true)
			if buyAmt.IsPositive() {
				ob.AddOrder(pool.Order(amm.Buy, matchPrice, buyAmt))
			}
			sellAmt := pool.SellAmountUnder(matchPrice, true)
			if sellAmt.IsPositive() {
				ob.AddOrder(pool.Order(amm.Sell, matchPrice, sellAmt))
			}
		}
		quoteCoinDiff, matched = ob.MatchAtSinglePrice(matchPrice)
	} else {
		lowestPrice, highestPrice := k.PriceLimits(ctx, *lastPrice)
		for _, pool := range pools {
			poolOrders := amm.PoolOrders(pool, pool, lowestPrice, highestPrice, tickPrec)
			ob.AddOrder(poolOrders...)
		}
		matchPrice, quoteCoinDiff, matched = ob.Match(*lastPrice)
	}
	return
}

func (k Keeper) ApplyMatchResult(ctx sdk.Context, pair types.Pair, orders []amm.Order, quoteCoinDiff math.Int) error {
	whitelisted := make([]string, 0)
	bulkOp := types.NewBulkSendCoinsOperation()
	for _, order := range orders { // TODO: need optimization to filter matched orders only
		order, ok := order.(*types.PoolOrder)
		if !ok {
			continue
		}
		if !order.IsMatched() {
			continue
		}
		paidCoin := sdk.NewCoin(order.OfferCoinDenom, order.PaidOfferCoinAmount)

		whitelisted = append(whitelisted, k.guardKeeper.AddTransferAccAddressesWhitelist([]string{order.ReserveAddress.String()})...)
		bulkOp.QueueSendCoins(order.ReserveAddress, pair.GetEscrowAddress(), sdk.NewCoins(paidCoin))
	}
	err := bulkOp.Run(ctx, k.bankKeeper)
	k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
	if err != nil {
		return err
	}
	bulkOp = types.NewBulkSendCoinsOperation()
	type PoolMatchResult struct {
		PoolId         uint64
		OrderDirection types.OrderDirection
		PaidCoin       sdk.Coin
		ReceivedCoin   sdk.Coin
		MatchedAmount  math.Int
	}
	poolMatchResultById := map[uint64]*PoolMatchResult{}
	var poolMatchResults []*PoolMatchResult
	for _, order := range orders {
		if !order.IsMatched() {
			continue
		}

		matchedAmt := order.GetAmount().Sub(order.GetOpenAmount())

		switch order := order.(type) {
		case *types.UserOrder:
			paidCoin := sdk.NewCoin(order.OfferCoinDenom, order.PaidOfferCoinAmount)
			receivedCoin := sdk.NewCoin(order.DemandCoinDenom, order.ReceivedDemandCoinAmount)

			o, _ := k.GetOrder(ctx, pair.Id, order.OrderId)
			o.OpenAmount = o.OpenAmount.Sub(matchedAmt)
			o.RemainingOfferCoin = o.RemainingOfferCoin.Sub(paidCoin)
			o.ReceivedCoin = o.ReceivedCoin.Add(receivedCoin)

			if o.OpenAmount.IsZero() {
				if err := k.FinishOrder(ctx, o, types.OrderStatusCompleted); err != nil {
					return err
				}
			} else {
				o.SetStatus(types.OrderStatusPartiallyMatched)
				k.SetOrder(ctx, o)
			}
			bulkOp.QueueSendCoins(pair.GetEscrowAddress(), order.Orderer, sdk.NewCoins(receivedCoin))

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeUserOrderMatched,
					sdk.NewAttribute(types.AttributeKeyOrderDirection, types.OrderDirectionFromAMM(order.Direction).String()),
					sdk.NewAttribute(types.AttributeKeyOrderer, order.Orderer.String()),
					sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(pair.Id, 10)),
					sdk.NewAttribute(types.AttributeKeyOrderId, strconv.FormatUint(order.OrderId, 10)),
					sdk.NewAttribute(types.AttributeKeyMatchedAmount, matchedAmt.String()),
					sdk.NewAttribute(types.AttributeKeyPaidCoin, paidCoin.String()),
					sdk.NewAttribute(types.AttributeKeyReceivedCoin, receivedCoin.String()),
				),
			})
		case *types.PoolOrder:
			paidCoin := sdk.NewCoin(order.OfferCoinDenom, order.PaidOfferCoinAmount)
			receivedCoin := sdk.NewCoin(order.DemandCoinDenom, order.ReceivedDemandCoinAmount)

			bulkOp.QueueSendCoins(pair.GetEscrowAddress(), order.ReserveAddress, sdk.NewCoins(receivedCoin))

			r, ok := poolMatchResultById[order.PoolId]
			if !ok {
				r = &PoolMatchResult{
					PoolId:         order.PoolId,
					OrderDirection: types.OrderDirectionFromAMM(order.Direction),
					PaidCoin:       sdk.NewCoin(paidCoin.Denom, math.ZeroInt()),
					ReceivedCoin:   sdk.NewCoin(receivedCoin.Denom, math.ZeroInt()),
					MatchedAmount:  math.ZeroInt(),
				}
				poolMatchResultById[order.PoolId] = r
				poolMatchResults = append(poolMatchResults, r)
			}
			dir := types.OrderDirectionFromAMM(order.Direction)
			if r.OrderDirection != dir {
				panic(fmt.Errorf("wrong order direction: %s != %s", dir, r.OrderDirection))
			}
			r.PaidCoin = r.PaidCoin.Add(paidCoin)
			r.ReceivedCoin = r.ReceivedCoin.Add(receivedCoin)
			r.MatchedAmount = r.MatchedAmount.Add(matchedAmt)
		default:
			panic(fmt.Errorf("invalid order type: %T", order))
		}
	}

	whitelisted = k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
	bulkOp.QueueSendCoins(pair.GetEscrowAddress(), k.GetDustCollector(ctx), sdk.NewCoins(sdk.NewCoin(pair.QuoteCoinDenom, quoteCoinDiff)))
	err = bulkOp.Run(ctx, k.bankKeeper)
	k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
	if err != nil {
		return err
	}
	for _, r := range poolMatchResults {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypePoolOrderMatched,
				sdk.NewAttribute(types.AttributeKeyOrderDirection, r.OrderDirection.String()),
				sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(pair.Id, 10)),
				sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(r.PoolId, 10)),
				sdk.NewAttribute(types.AttributeKeyMatchedAmount, r.MatchedAmount.String()),
				sdk.NewAttribute(types.AttributeKeyPaidCoin, r.PaidCoin.String()),
				sdk.NewAttribute(types.AttributeKeyReceivedCoin, r.ReceivedCoin.String()),
			),
		})
	}
	return nil
}

func (k Keeper) FinishOrder(ctx sdk.Context, order types.Order, status types.OrderStatus) error {
	if order.Type == types.OrderTypeMM {
		return k.FinishMMOrder(ctx, order, status)
	}
	if order.Status == types.OrderStatusCompleted || order.Status.IsCanceledOrExpired() { // sanity check
		return nil
	}

	pair, _ := k.GetPair(ctx, order.PairId)

	var swapFeeRate sdkmath.LegacyDec
	if pair.SwapFeeRate != nil {
		swapFeeRate = *pair.SwapFeeRate
	} else {
		swapFeeRate = k.GetSwapFeeRate(ctx)
	}

	accumulatedSwapFee := sdk.NewCoin(order.OfferCoin.Denom, math.NewInt(0))
	collectedSwapFeeAmountFromOrderer := order.SwapFeeCoin.Amount

	pairCreatorSwapFeeCoin := sdk.NewCoin(accumulatedSwapFee.Denom, math.NewInt(0))
	rewardsSwapFeeCoin := sdk.NewCoin(accumulatedSwapFee.Denom, math.NewInt(0))
	refundCoin := sdk.NewCoin(order.RemainingOfferCoin.Denom, math.NewInt(0))

	if order.RemainingOfferCoin.IsPositive() {
		refundCoin = sdk.NewCoin(order.RemainingOfferCoin.Denom, order.RemainingOfferCoin.Amount)

		if refundCoin.IsEqual(order.OfferCoin) {
			// refund full swap fees back to orderer
			refundCoin.Amount = refundCoin.Amount.Add(collectedSwapFeeAmountFromOrderer)
		} else {
			// refund partial swap fees back to orderer and transfer remaining to swap fee collector address
			swappedCoin := order.OfferCoin.Sub(refundCoin)
			swapFeeAmt := CalculateSwapFeeAmount(ctx, swapFeeRate, swappedCoin.Amount)

			refundableSwapFeeAmt := collectedSwapFeeAmountFromOrderer.Sub(swapFeeAmt)
			if refundableSwapFeeAmt.IsNegative() { // refundable swap fee amount shouldn't be negative, edge case
				refundableSwapFeeAmt = math.ZeroInt()
				swapFeeAmt = collectedSwapFeeAmountFromOrderer
			}
			accumulatedSwapFee.Amount = accumulatedSwapFee.Amount.Add(swapFeeAmt)
			refundCoin.Amount = refundCoin.Amount.Add(refundableSwapFeeAmt)
		}

		if refundCoin.IsPositive() {

			whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
			err := k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), order.GetOrderer(), sdk.NewCoins(refundCoin))
			k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
			if err != nil {
				return err
			}
		}
	} else {
		accumulatedSwapFee.Amount = accumulatedSwapFee.Amount.Add(collectedSwapFeeAmountFromOrderer)
	}

	totalSwapFeeCoin := sdk.NewCoin(accumulatedSwapFee.Denom, accumulatedSwapFee.Amount)

	if accumulatedSwapFee.IsPositive() {
		var pairCreatorSwapFeeRatio sdkmath.LegacyDec
		if pair.PairCreatorSwapFeeRatio != nil {
			pairCreatorSwapFeeRatio = *pair.PairCreatorSwapFeeRatio
		} else {
			pairCreatorSwapFeeRatio = k.GetPairCreatorSwapFeeRatio(ctx)
		}

		if pairCreatorSwapFeeRatio.IsPositive() {
			pairCreatorSwapFeeAmt := CalculatePairCreatorSwapFeeAmount(ctx, pairCreatorSwapFeeRatio, accumulatedSwapFee.Amount)
			accumulatedSwapFee.Amount = accumulatedSwapFee.Amount.Sub(pairCreatorSwapFeeAmt)
			pairCreatorSwapFeeCoin = sdk.NewCoin(accumulatedSwapFee.Denom, pairCreatorSwapFeeAmt)

			if pairCreatorSwapFeeCoin.IsPositive() {
				pairCreator, err := sdk.AccAddressFromBech32(pair.Creator)
				if err != nil {
					k.logger.Error("fail to parse pair creator address", "err", err)
					return nil
				}

				whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
				err = k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), pairCreator, sdk.NewCoins(pairCreatorSwapFeeCoin))
				k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
				if err != nil {
					k.logger.Error("fail to send coins to pair creator", "err", err)
					return nil
				}
			}
		}

		if accumulatedSwapFee.IsPositive() {

			whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
			rewardsSwapFeeCoin = sdk.NewCoin(accumulatedSwapFee.Denom, accumulatedSwapFee.Amount)
			err := k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), pair.GetSwapFeeCollectorAddress(), sdk.NewCoins(rewardsSwapFeeCoin))
			k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
			if err != nil {
				k.logger.Error("fail to send coins to swap fee collector", "err", err)
				return nil
			}
		}
	}

	order.SetStatus(status)
	k.SetOrder(ctx, order)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOrderResult,
			sdk.NewAttribute(types.AttributeKeyOrderDirection, order.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOrderer, order.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(order.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderId, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, order.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOpenAmount, order.OpenAmount.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, order.OfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingOfferCoin, order.RemainingOfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyReceivedCoin, order.ReceivedCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRefundedCoins, refundCoin.String()),
			sdk.NewAttribute(types.AttributeKeySwapFeeCoins, totalSwapFeeCoin.String()),
			sdk.NewAttribute(types.AttributeKeyPairCreatorSwapFeeCoins, pairCreatorSwapFeeCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRewardsSwapFeeCoins, rewardsSwapFeeCoin.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, order.Status.String()),
		),
	})

	return nil
}

func (k Keeper) FinishMMOrder(ctx sdk.Context, order types.Order, status types.OrderStatus) error {
	if order.Status == types.OrderStatusCompleted || order.Status.IsCanceledOrExpired() { // sanity check
		return nil
	}

	if order.RemainingOfferCoin.IsPositive() {
		pair, _ := k.GetPair(ctx, order.PairId)

		whitelisted := k.guardKeeper.AddTransferAccAddressesWhitelist([]string{pair.GetEscrowAddress().String()})
		err := k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), order.GetOrderer(), sdk.NewCoins(order.RemainingOfferCoin))
		k.guardKeeper.RemoveTransferAccAddressesWhitelist(whitelisted)
		if err != nil {
			return err
		}
	}

	order.SetStatus(status)
	k.SetOrder(ctx, order)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOrderResult,
			sdk.NewAttribute(types.AttributeKeyOrderDirection, order.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOrderer, order.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairId, strconv.FormatUint(order.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderId, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, order.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOpenAmount, order.OpenAmount.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, order.OfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingOfferCoin, order.RemainingOfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyReceivedCoin, order.ReceivedCoin.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, order.Status.String()),
		),
	})

	return nil
}
