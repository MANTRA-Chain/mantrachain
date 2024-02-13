package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/txfees module sentinel errors
var (
	ErrSample                                 = sdkerrors.Register(ModuleName, 2000, "sample error")
	ErrInvalidBaseDenomParam                  = sdkerrors.Register(ModuleName, 2001, "invalid base denom param")
	ErrTooManyFeeCoins                        = sdkerrors.Register(ModuleName, 2002, "too many fee coins")
	ErrInvalidFeeDenom                        = sdkerrors.Register(ModuleName, 2003, "invalid fee denom")
	ErrLiquidityPoolPairNotFound              = sdkerrors.Register(ModuleName, 2004, "Liquidity pool pair not found")
	ErrLiquidityPoolPairLastPriceNotAvailable = sdkerrors.Register(ModuleName, 2005, "Liquidity pool pair last price not available")
	ErrLiquidityPoolPairFeeDenomNotMatch      = sdkerrors.Register(ModuleName, 2006, "Liquidity pool pair fee denom not match")
	ErrTooManyGasPricesCoins                  = sdkerrors.Register(ModuleName, 2007, "too many gas prices coins")
	ErrZeroFee                                = sdkerrors.Register(ModuleName, 2008, "zero fee")
	ErrInvalidAmount                          = sdkerrors.Register(ModuleName, 2009, "invalid amount")
)
