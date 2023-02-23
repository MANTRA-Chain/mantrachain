package types

import "cosmossdk.io/errors"

// DONTCOVER

// x/liquidity module sentinel errors
var (
	ErrInsufficientDepositAmount = errors.Register(ModuleName, 2, "insufficient deposit amount")
	ErrPairAlreadyExists         = errors.Register(ModuleName, 3, "pair already exists")
	ErrPoolAlreadyExists         = errors.Register(ModuleName, 4, "pool already exists")
	ErrWrongPoolCoinDenom        = errors.Register(ModuleName, 5, "wrong pool coin denom")
	ErrInvalidCoinDenom          = errors.Register(ModuleName, 6, "invalid coin denom")
	ErrNoLastPrice               = errors.Register(ModuleName, 8, "cannot make a market order to a pair with no last price")
	ErrInsufficientOfferCoin     = errors.Register(ModuleName, 9, "insufficient offer coin")
	ErrPriceOutOfRange           = errors.Register(ModuleName, 10, "price out of range limit")
	ErrTooLongOrderLifespan      = errors.Register(ModuleName, 11, "order lifespan is too long")
	ErrDisabledPool              = errors.Register(ModuleName, 12, "disabled pool")
	ErrWrongPair                 = errors.Register(ModuleName, 13, "wrong denom pair")
	ErrSameBatch                 = errors.Register(ModuleName, 14, "cannot cancel an order within the same batch")
	ErrAlreadyCanceled           = errors.Register(ModuleName, 15, "the order is already canceled")
	ErrDuplicatePairId           = errors.Register(ModuleName, 16, "duplicate pair id presents in the pair id list")
	ErrTooSmallOrder             = errors.Register(ModuleName, 17, "too small order")
	ErrTooLargePool              = errors.Register(ModuleName, 18, "too large pool")
	ErrTooManyPools              = errors.Register(ModuleName, 19, "too many pools in the pair")
	ErrPriceNotOnTicks           = errors.Register(ModuleName, 20, "price is not on ticks")
	ErrMaxNumMMOrdersExceeded    = errors.Register(ModuleName, 21, "number of MM orders exceeded the limit")
)
