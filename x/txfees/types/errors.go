package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/txfees module sentinel errors
var (
	ErrSample                                 = sdkerrors.Register(ModuleName, 2, "sample error")
	ErrInvalidBaseDenomParam                  = sdkerrors.Register(ModuleName, 3, "invalid base denom param")
	ErrTooManyFeeCoins                        = sdkerrors.Register(ModuleName, 4, "too many fee coins")
	ErrInvalidFeeDenom                        = sdkerrors.Register(ModuleName, 5, "invalid fee denom")
	ErrLiquidityPoolPairNotFound              = sdkerrors.Register(ModuleName, 6, "Liquidity pool pair not found")
	ErrLiquidityPoolPairLastPriceNotAvailable = sdkerrors.Register(ModuleName, 7, "Liquidity pool pair last price not available")
	ErrLiquidityPoolPairFeeDenomNotMatch      = sdkerrors.Register(ModuleName, 8, "Liquidity pool pair fee denom not match")
	ErrTooManyGasPricesCoins                  = sdkerrors.Register(ModuleName, 9, "too many gas prices coins")
	ErrZeroFee                                = sdkerrors.Register(ModuleName, 10, "zero fee")
	ErrInvalidAmount                          = sdkerrors.Register(ModuleName, 11, "invalid amount")
)
