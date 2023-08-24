package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/txfees module sentinel errors
var (
	ErrSample                                 = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidBaseDenomParam                  = sdkerrors.Register(ModuleName, 1101, "invalid base denom param")
	ErrTooManyFeeCoins                        = sdkerrors.Register(ModuleName, 1102, "too many fee coins")
	ErrInvalidFeeDenom                        = sdkerrors.Register(ModuleName, 1103, "invalid fee denom")
	ErrLiquidityPoolPairNotFound              = sdkerrors.Register(ModuleName, 1104, "Liquidity pool pair not found")
	ErrLiquidityPoolPairLastPriceNotAvailable = sdkerrors.Register(ModuleName, 1105, "Liquidity pool pair last price not available")
	ErrLiquidityPoolPairFeeDenomNotMatch      = sdkerrors.Register(ModuleName, 1106, "Liquidity pool pair fee denom not match")
	ErrTooManyGasPricesCoins                  = sdkerrors.Register(ModuleName, 1107, "too many gas prices coins")
	ErrZeroFee                                = sdkerrors.Register(ModuleName, 1108, "zero fee")
)
