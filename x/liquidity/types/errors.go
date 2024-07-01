package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/liquidity module sentinel errors
var (
	ErrInvalidSigner             = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInsufficientDepositAmount = errors.Register(ModuleName, 1101, "insufficient deposit amount")
	ErrPairAlreadyExists         = errors.Register(ModuleName, 1102, "pair already exists")
	ErrPoolAlreadyExists         = errors.Register(ModuleName, 1103, "pool already exists")
	ErrWrongPoolCoinDenom        = errors.Register(ModuleName, 1104, "wrong pool coin denom")
	ErrInvalidCoinDenom          = errors.Register(ModuleName, 1105, "invalid coin denom")
	ErrNoLastPrice               = errors.Register(ModuleName, 1106, "cannot make a market order to a pair with no last price")
	ErrInsufficientOfferCoin     = errors.Register(ModuleName, 1107, "insufficient offer coin")
	ErrPriceOutOfRange           = errors.Register(ModuleName, 1108, "price out of range limit")
	ErrTooLongOrderLifespan      = errors.Register(ModuleName, 1109, "order lifespan is too long")
	ErrDisabledPool              = errors.Register(ModuleName, 1110, "disabled pool")
	ErrWrongPair                 = errors.Register(ModuleName, 1111, "wrong denom pair")
	ErrSameBatch                 = errors.Register(ModuleName, 1112, "cannot cancel an order within the same batch")
	ErrAlreadyCanceled           = errors.Register(ModuleName, 1113, "the order is already canceled")
	ErrDuplicatePairId           = errors.Register(ModuleName, 1114, "duplicate pair id presents in the pair id list")
	ErrTooSmallOrder             = errors.Register(ModuleName, 1115, "too small order")
	ErrTooLargePool              = errors.Register(ModuleName, 1116, "too large pool")
	ErrTooManyPools              = errors.Register(ModuleName, 1117, "too many pools in the pair")
	ErrPriceNotOnTicks           = errors.Register(ModuleName, 1118, "price is not on ticks")
)
