package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DONTCOVER

// x/liquidity module sentinel errors
var (
	ErrInsufficientDepositAmount = sdkerrors.Register(ModuleName, 1501, "insufficient deposit amount")
	ErrPairAlreadyExists         = sdkerrors.Register(ModuleName, 1502, "pair already exists")
	ErrPoolAlreadyExists         = sdkerrors.Register(ModuleName, 1503, "pool already exists")
	ErrWrongPoolCoinDenom        = sdkerrors.Register(ModuleName, 1504, "wrong pool coin denom")
	ErrInvalidCoinDenom          = sdkerrors.Register(ModuleName, 1505, "invalid coin denom")
	ErrNoLastPrice               = sdkerrors.Register(ModuleName, 1506, "cannot make a market order to a pair with no last price")
	ErrInsufficientOfferCoin     = sdkerrors.Register(ModuleName, 1507, "insufficient offer coin")
	ErrPriceOutOfRange           = sdkerrors.Register(ModuleName, 1508, "price out of range limit")
	ErrTooLongOrderLifespan      = sdkerrors.Register(ModuleName, 1509, "order lifespan is too long")
	ErrDisabledPool              = sdkerrors.Register(ModuleName, 1510, "disabled pool")
	ErrWrongPair                 = sdkerrors.Register(ModuleName, 1511, "wrong denom pair")
	ErrSameBatch                 = sdkerrors.Register(ModuleName, 1512, "cannot cancel an order within the same batch")
	ErrAlreadyCanceled           = sdkerrors.Register(ModuleName, 1513, "the order is already canceled")
	ErrDuplicatePairId           = sdkerrors.Register(ModuleName, 1514, "duplicate pair id presents in the pair id list")
	ErrTooSmallOrder             = sdkerrors.Register(ModuleName, 1515, "too small order")
	ErrTooLargePool              = sdkerrors.Register(ModuleName, 1516, "too large pool")
	ErrTooManyPools              = sdkerrors.Register(ModuleName, 1517, "too many pools in the pair")
	ErrPriceNotOnTicks           = sdkerrors.Register(ModuleName, 1518, "price is not on ticks")
)
