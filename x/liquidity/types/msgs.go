package types

import (
	"time"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	"github.com/MANTRA-Finance/mantrachain/x/liquidity/amm"
)

var (
	_ sdk.Msg = &MsgCreatePair{}
	_ sdk.Msg = &MsgCreatePool{}
	_ sdk.Msg = &MsgCreateRangedPool{}
	_ sdk.Msg = &MsgDeposit{}
	_ sdk.Msg = &MsgWithdraw{}
	_ sdk.Msg = &MsgLimitOrder{}
	_ sdk.Msg = &MsgMarketOrder{}
	_ sdk.Msg = &MsgMMOrder{}
	_ sdk.Msg = &MsgCancelOrder{}
	_ sdk.Msg = &MsgCancelAllOrders{}
	_ sdk.Msg = &MsgCancelMMOrder{}

	_ legacytx.LegacyMsg = &MsgCreatePair{}
	_ legacytx.LegacyMsg = &MsgCreatePool{}
	_ legacytx.LegacyMsg = &MsgCreateRangedPool{}
	_ legacytx.LegacyMsg = &MsgDeposit{}
	_ legacytx.LegacyMsg = &MsgLimitOrder{}
	_ legacytx.LegacyMsg = &MsgMarketOrder{}
	_ legacytx.LegacyMsg = &MsgMMOrder{}
	_ legacytx.LegacyMsg = &MsgCancelOrder{}
	_ legacytx.LegacyMsg = &MsgCancelAllOrders{}
	_ legacytx.LegacyMsg = &MsgCancelMMOrder{}
)

// Message types for the liquidity module
const (
	TypeMsgCreatePair        = "create_pair"
	TypeMsgUpdatePairSwapFee = "update_pair_swap_fee"
	TypeMsgCreatePool        = "create_pool"
	TypeMsgCreateRangedPool  = "create_ranged_pool"
	TypeMsgDeposit           = "deposit"
	TypeMsgWithdraw          = "withdraw"
	TypeMsgLimitOrder        = "limit_order"
	TypeMsgMarketOrder       = "market_order"
	TypeMsgMMOrder           = "mm_order"
	TypeMsgCancelOrder       = "cancel_order"
	TypeMsgCancelAllOrders   = "cancel_all_orders"
	TypeMsgCancelMMOrder     = "cancel_mm_order"
)

// NewMsgCreatePair returns a new MsgCreatePair.
func NewMsgCreatePair(creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string, swapFeeRate *sdkmath.LegacyDec, pairCreatorSwapFeeRatio *sdkmath.LegacyDec) *MsgCreatePair {
	if swapFeeRate.IsNil() && !pairCreatorSwapFeeRatio.IsNil() {
		panic("swap fee rate must not be nil when pair creator swap fee ratio is not nil")
	}
	if !swapFeeRate.IsNil() && pairCreatorSwapFeeRatio.IsNil() {
		panic("pair creator swap fee ratio must not be nil when swap fee rate is not nil")
	}
	if swapFeeRate.IsNil() || pairCreatorSwapFeeRatio.IsNil() {
		return &MsgCreatePair{
			Creator:        creator.String(),
			BaseCoinDenom:  baseCoinDenom,
			QuoteCoinDenom: quoteCoinDenom,
		}
	}
	return &MsgCreatePair{
		Creator:                 creator.String(),
		BaseCoinDenom:           baseCoinDenom,
		QuoteCoinDenom:          quoteCoinDenom,
		SwapFeeRate:             swapFeeRate,
		PairCreatorSwapFeeRatio: pairCreatorSwapFeeRatio,
	}
}

func (msg MsgCreatePair) Route() string { return RouterKey }

func (msg MsgCreatePair) Type() string { return TypeMsgCreatePair }

func (msg MsgCreatePair) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if err := sdk.ValidateDenom(msg.BaseCoinDenom); err != nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, err.Error())
	}
	if err := sdk.ValidateDenom(msg.QuoteCoinDenom); err != nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, err.Error())
	}
	if msg.BaseCoinDenom == msg.QuoteCoinDenom {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "cannot use same denom for both base coin and quote coin")
	}
	if msg.SwapFeeRate != nil && msg.SwapFeeRate.IsNegative() {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "swap fee rate must not be negative")
	}
	if msg.PairCreatorSwapFeeRatio != nil && msg.PairCreatorSwapFeeRatio.IsNegative() {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair creator swap fee ratio must not be negative")
	}
	if msg.SwapFeeRate == nil && msg.PairCreatorSwapFeeRatio != nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "swap fee rate must not be nil when pair creator swap fee ratio is not nil")
	}
	if msg.SwapFeeRate != nil && msg.PairCreatorSwapFeeRatio == nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair creator swap fee ratio must not be nil when swap fee rate is not nil")
	}
	return nil
}

func (msg *MsgCreatePair) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreatePair) GetAccCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgUpdatePairSwapFee returns a new MsgUpdatePairSwapFee.
func NewMsgUpdatePairSwapFee(creator sdk.AccAddress, pairId uint64, swapFeeRate *sdkmath.LegacyDec, pairCreatorSwapFeeRatio *sdkmath.LegacyDec) *MsgUpdatePairSwapFee {
	if swapFeeRate.IsNil() && !pairCreatorSwapFeeRatio.IsNil() {
		panic("swap fee rate must not be nil when pair creator swap fee ratio is not nil")
	}
	if !swapFeeRate.IsNil() && pairCreatorSwapFeeRatio.IsNil() {
		panic("pair creator swap fee ratio must not be nil when swap fee rate is not nil")
	}
	if swapFeeRate.IsNil() || pairCreatorSwapFeeRatio.IsNil() {
		return &MsgUpdatePairSwapFee{
			Creator: creator.String(),
			PairId:  pairId,
		}
	}
	return &MsgUpdatePairSwapFee{
		Creator:                 creator.String(),
		PairId:                  pairId,
		SwapFeeRate:             swapFeeRate,
		PairCreatorSwapFeeRatio: pairCreatorSwapFeeRatio,
	}
}

func (msg MsgUpdatePairSwapFee) Route() string { return RouterKey }

func (msg MsgUpdatePairSwapFee) Type() string { return TypeMsgUpdatePairSwapFee }

func (msg MsgUpdatePairSwapFee) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if msg.SwapFeeRate != nil && msg.SwapFeeRate.IsNegative() {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "swap fee rate must not be negative")
	}
	if msg.PairCreatorSwapFeeRatio != nil && msg.PairCreatorSwapFeeRatio.IsNegative() {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair creator swap fee ratio must not be negative")
	}
	if msg.SwapFeeRate == nil && msg.PairCreatorSwapFeeRatio != nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "swap fee rate must not be nil when pair creator swap fee ratio is not nil")
	}
	if msg.SwapFeeRate != nil && msg.PairCreatorSwapFeeRatio == nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair creator swap fee ratio must not be nil when swap fee rate is not nil")
	}
	return nil
}

func (msg *MsgUpdatePairSwapFee) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgUpdatePairSwapFee) GetAccCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCreatePool creates a new MsgCreatePool.
func NewMsgCreatePool(
	creator sdk.AccAddress,
	pairId uint64,
	depositCoins sdk.Coins,
) *MsgCreatePool {
	return &MsgCreatePool{
		Creator:      creator.String(),
		PairId:       pairId,
		DepositCoins: depositCoins,
	}
}

func (msg MsgCreatePool) Route() string { return RouterKey }

func (msg MsgCreatePool) Type() string { return TypeMsgCreatePool }

func (msg MsgCreatePool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if msg.PairId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) != 2 {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.GT(amm.MaxCoinAmount) {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "deposit coin %s is bigger than the max amount %s", coin, amm.MaxCoinAmount)
		}
	}
	return nil
}

func (msg *MsgCreatePool) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreatePool) GetAccCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCreateRangedPool creates a new MsgCreateRangedPool.
func NewMsgCreateRangedPool(
	creator sdk.AccAddress,
	pairId uint64,
	depositCoins sdk.Coins,
	minPrice sdkmath.LegacyDec,
	maxPrice sdkmath.LegacyDec,
	initialPrice sdkmath.LegacyDec,
) *MsgCreateRangedPool {
	return &MsgCreateRangedPool{
		Creator:      creator.String(),
		PairId:       pairId,
		DepositCoins: depositCoins,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		InitialPrice: initialPrice,
	}
}

func (msg MsgCreateRangedPool) Route() string { return RouterKey }

func (msg MsgCreateRangedPool) Type() string { return TypeMsgCreateRangedPool }

func (msg MsgCreateRangedPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if msg.PairId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) == 0 || len(msg.DepositCoins) > 2 {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	for _, coin := range msg.DepositCoins {
		if coin.Amount.GT(amm.MaxCoinAmount) {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "deposit coin %s is bigger than the max amount %s", coin, amm.MaxCoinAmount)
		}
	}
	if err := amm.ValidateRangedPoolParams(msg.MinPrice, msg.MaxPrice, msg.InitialPrice); err != nil {
		return errors.Wrap(errorstypes.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (msg *MsgCreateRangedPool) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateRangedPool) GetAccCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgDeposit creates a new MsgDeposit.
func NewMsgDeposit(
	depositor sdk.AccAddress,
	poolId uint64,
	depositCoins sdk.Coins,
) *MsgDeposit {
	return &MsgDeposit{
		Depositor:    depositor.String(),
		PoolId:       poolId,
		DepositCoins: depositCoins,
	}
}

func (msg MsgDeposit) Route() string { return RouterKey }

func (msg MsgDeposit) Type() string { return TypeMsgDeposit }

func (msg MsgDeposit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Depositor); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid depositor address: %v", err)
	}
	if msg.PoolId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if len(msg.DepositCoins) == 0 || len(msg.DepositCoins) > 2 {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "wrong number of deposit coins: %d", len(msg.DepositCoins))
	}
	return nil
}

func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDeposit) GetAccDepositor() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgWithdraw creates a new MsgWithdraw.
func NewMsgWithdraw(
	withdrawer sdk.AccAddress,
	poolId uint64,
	poolCoin sdk.Coin,
) *MsgWithdraw {
	return &MsgWithdraw{
		Withdrawer: withdrawer.String(),
		PoolId:     poolId,
		PoolCoin:   poolCoin,
	}
}

func (msg MsgWithdraw) Route() string { return RouterKey }

func (msg MsgWithdraw) Type() string { return TypeMsgWithdraw }

func (msg MsgWithdraw) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Withdrawer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid withdrawer address: %v", err)
	}
	if msg.PoolId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pool id must not be 0")
	}
	if err := msg.PoolCoin.Validate(); err != nil {
		return err
	}
	if !msg.PoolCoin.IsPositive() {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pool coin must be positive")
	}
	return nil
}

func (msg *MsgWithdraw) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgWithdraw) GetAccWithdrawer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgLimitOrder creates a new MsgLimitOrder.
func NewMsgLimitOrder(
	orderer sdk.AccAddress,
	pairId uint64,
	dir OrderDirection,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	price sdkmath.LegacyDec,
	amt sdkmath.Int,
	orderLifespan time.Duration,
) *MsgLimitOrder {
	return &MsgLimitOrder{
		Orderer:         orderer.String(),
		PairId:          pairId,
		Direction:       dir,
		OfferCoin:       offerCoin,
		DemandCoinDenom: demandCoinDenom,
		Price:           price,
		Amount:          amt,
		OrderLifespan:   orderLifespan,
	}
}

func (msg MsgLimitOrder) Route() string { return RouterKey }

func (msg MsgLimitOrder) Type() string { return TypeMsgLimitOrder }

func (msg MsgLimitOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.Direction != OrderDirectionBuy && msg.Direction != OrderDirectionSell {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "invalid order direction: %s", msg.Direction)
	}
	if err := sdk.ValidateDenom(msg.DemandCoinDenom); err != nil {
		return errors.Wrap(err, "invalid demand coin denom")
	}
	if !msg.Price.IsPositive() {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "price must be positive")
	}
	if err := msg.OfferCoin.Validate(); err != nil {
		return errors.Wrap(err, "invalid offer coin")
	}
	if msg.OfferCoin.Amount.LT(amm.MinCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "offer coin %s is smaller than the min amount %s", msg.OfferCoin, amm.MinCoinAmount)
	}
	if msg.OfferCoin.Amount.GT(amm.MaxCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "offer coin %s is bigger than the max amount %s", msg.OfferCoin, amm.MaxCoinAmount)
	}
	if msg.Amount.LT(amm.MinCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "order amount %s is smaller than the min amount %s", msg.Amount, amm.MinCoinAmount)
	}
	if msg.Amount.GT(amm.MaxCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "order amount %s is bigger than the max amount %s", msg.Amount, amm.MaxCoinAmount)
	}
	var minOfferCoin sdk.Coin
	switch msg.Direction {
	case OrderDirectionBuy:
		minOfferCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, msg.Price, msg.Amount))
	case OrderDirectionSell:
		minOfferCoin = sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount)
	}
	if msg.OfferCoin.IsLT(minOfferCoin) {
		return errors.Wrapf(ErrInsufficientOfferCoin, "%s is less than %s", msg.OfferCoin, minOfferCoin)
	}
	if msg.OfferCoin.Denom == msg.DemandCoinDenom {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "offer coin denom and demand coin denom must not be same")
	}
	if msg.OrderLifespan < 0 {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
	}
	return nil
}

func (msg *MsgLimitOrder) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgLimitOrder) GetAccOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgMarketOrder creates a new MsgMarketOrder.
func NewMsgMarketOrder(
	orderer sdk.AccAddress,
	pairId uint64,
	dir OrderDirection,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	amt sdkmath.Int,
	orderLifespan time.Duration,
) *MsgMarketOrder {
	return &MsgMarketOrder{
		Orderer:         orderer.String(),
		PairId:          pairId,
		Direction:       dir,
		OfferCoin:       offerCoin,
		DemandCoinDenom: demandCoinDenom,
		Amount:          amt,
		OrderLifespan:   orderLifespan,
	}
}

func (msg MsgMarketOrder) Route() string { return RouterKey }

func (msg MsgMarketOrder) Type() string { return TypeMsgMarketOrder }

func (msg MsgMarketOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.Direction != OrderDirectionBuy && msg.Direction != OrderDirectionSell {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "invalid order direction: %s", msg.Direction)
	}
	if err := sdk.ValidateDenom(msg.DemandCoinDenom); err != nil {
		return errors.Wrap(err, "invalid demand coin denom")
	}
	if err := msg.OfferCoin.Validate(); err != nil {
		return errors.Wrap(err, "invalid offer coin")
	}
	if msg.OfferCoin.Amount.LT(amm.MinCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "offer coin %s is smaller than the min amount %s", msg.OfferCoin, amm.MinCoinAmount)
	}
	if msg.OfferCoin.Amount.GT(amm.MaxCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "offer coin %s is bigger than the max amount %s", msg.OfferCoin, amm.MaxCoinAmount)
	}
	if msg.Amount.LT(amm.MinCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "order amount %s is smaller than the min amount %s", msg.Amount, amm.MinCoinAmount)
	}
	if msg.Amount.GT(amm.MaxCoinAmount) {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "order amount %s is bigger than the max amount %s", msg.Amount, amm.MaxCoinAmount)
	}
	if msg.OfferCoin.Denom == msg.DemandCoinDenom {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "offer coin denom and demand coin denom must not be same")
	}
	if msg.OrderLifespan < 0 {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
	}
	return nil
}

func (msg *MsgMarketOrder) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMarketOrder) GetAccOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgMMOrder creates a new MsgMMOrder.
func NewMsgMMOrder(
	orderer sdk.AccAddress,
	pairId uint64,
	maxSellPrice, minSellPrice sdkmath.LegacyDec, sellAmt sdkmath.Int,
	maxBuyPrice, minBuyPrice sdkmath.LegacyDec, buyAmt sdkmath.Int,
	orderLifespan time.Duration,
) *MsgMMOrder {
	return &MsgMMOrder{
		Orderer:       orderer.String(),
		PairId:        pairId,
		MaxSellPrice:  maxSellPrice,
		MinSellPrice:  minSellPrice,
		SellAmount:    sellAmt,
		MaxBuyPrice:   maxBuyPrice,
		MinBuyPrice:   minBuyPrice,
		BuyAmount:     buyAmt,
		OrderLifespan: orderLifespan,
	}
}

func (msg MsgMMOrder) Route() string { return RouterKey }

func (msg MsgMMOrder) Type() string { return TypeMsgMMOrder }

func (msg MsgMMOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.SellAmount.IsZero() && msg.BuyAmount.IsZero() {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "sell amount and buy amount must not be zero at the same time")
	}
	if !msg.SellAmount.IsZero() {
		if msg.SellAmount.LT(amm.MinCoinAmount) {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "sell amount %s is smaller than the min amount %s", msg.SellAmount, amm.MinCoinAmount)
		}
		if !msg.MaxSellPrice.IsPositive() {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "max sell price must be positive: %s", msg.MaxSellPrice)
		}
		if !msg.MinSellPrice.IsPositive() {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "min sell price must be positive: %s", msg.MinSellPrice)
		}
		if msg.MinSellPrice.GT(msg.MaxSellPrice) {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "max sell price must not be lower than min sell price")
		}
	}
	if !msg.BuyAmount.IsZero() {
		if msg.BuyAmount.LT(amm.MinCoinAmount) {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "buy amount %s is smaller than the min amount %s", msg.BuyAmount, amm.MinCoinAmount)
		}
		if !msg.MinBuyPrice.IsPositive() {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "min buy price must be positive: %s", msg.MinBuyPrice)
		}
		if !msg.MaxBuyPrice.IsPositive() {
			return errors.Wrapf(errorstypes.ErrInvalidRequest, "max buy price must be positive: %s", msg.MaxBuyPrice)
		}
		if msg.MinBuyPrice.GT(msg.MaxBuyPrice) {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "max buy price must not be lower than min buy price")
		}
	}
	if msg.OrderLifespan < 0 {
		return errors.Wrapf(errorstypes.ErrInvalidRequest, "order lifespan must not be negative: %s", msg.OrderLifespan)
	}
	return nil
}

func (msg *MsgMMOrder) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgMMOrder) GetAccOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCancelOrder creates a new MsgCancelOrder.
func NewMsgCancelOrder(
	orderer sdk.AccAddress,
	pairId uint64,
	orderId uint64,
) *MsgCancelOrder {
	return &MsgCancelOrder{
		OrderId: orderId,
		PairId:  pairId,
		Orderer: orderer.String(),
	}
}

func (msg MsgCancelOrder) Route() string { return RouterKey }

func (msg MsgCancelOrder) Type() string { return TypeMsgCancelOrder }

func (msg MsgCancelOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
	}
	if msg.OrderId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "order id must not be 0")
	}
	return nil
}

func (msg *MsgCancelOrder) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCancelOrder) GetAccOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCancelAllOrders creates a new MsgCancelAllOrders.
func NewMsgCancelAllOrders(
	orderer sdk.AccAddress,
	pairIds []uint64,
) *MsgCancelAllOrders {
	return &MsgCancelAllOrders{
		Orderer: orderer.String(),
		PairIds: pairIds,
	}
}

func (msg MsgCancelAllOrders) Route() string { return RouterKey }

func (msg MsgCancelAllOrders) Type() string { return TypeMsgCancelAllOrders }

func (msg MsgCancelAllOrders) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	pairIdSet := map[uint64]struct{}{}
	for _, pairId := range msg.PairIds {
		if pairId == 0 {
			return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
		}
		if _, ok := pairIdSet[pairId]; ok {
			return ErrDuplicatePairId
		}
		pairIdSet[pairId] = struct{}{}
	}
	return nil
}

func (msg *MsgCancelAllOrders) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCancelAllOrders) GetAccOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgCancelMMOrder creates a new MsgCancelMMOrder.
func NewMsgCancelMMOrder(
	orderer sdk.AccAddress,
	pairId uint64,
) *MsgCancelMMOrder {
	return &MsgCancelMMOrder{
		Orderer: orderer.String(),
		PairId:  pairId,
	}
}

func (msg MsgCancelMMOrder) Route() string { return RouterKey }

func (msg MsgCancelMMOrder) Type() string { return TypeMsgCancelMMOrder }

func (msg MsgCancelMMOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Orderer); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid orderer address: %v", err)
	}
	if msg.PairId == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair id must not be 0")
	}
	return nil
}

func (msg *MsgCancelMMOrder) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCancelMMOrder) GetAccOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}
