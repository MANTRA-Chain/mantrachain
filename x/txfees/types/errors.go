package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/txfees module sentinel errors
var (
	ErrInvalidSigner                          = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidBaseDenomParam                  = errors.Register(ModuleName, 1101, "invalid base denom param")
	ErrTooManyFeeCoins                        = errors.Register(ModuleName, 1102, "too many fee coins")
	ErrInvalidFeeDenom                        = errors.Register(ModuleName, 1103, "invalid fee denom")
	ErrLiquidityPoolPairNotFound              = errors.Register(ModuleName, 1104, "Liquidity pool pair not found")
	ErrLiquidityPoolPairLastPriceNotAvailable = errors.Register(ModuleName, 1105, "Liquidity pool pair last price not available")
	ErrLiquidityPoolPairFeeDenomNotMatch      = errors.Register(ModuleName, 1106, "Liquidity pool pair fee denom not match")
	ErrTooManyGasPricesCoins                  = errors.Register(ModuleName, 1107, "too many gas prices coins")
	ErrZeroFee                                = errors.Register(ModuleName, 1108, "zero fee")
	ErrInvalidAmount                          = errors.Register(ModuleName, 1109, "invalid amount")
)
