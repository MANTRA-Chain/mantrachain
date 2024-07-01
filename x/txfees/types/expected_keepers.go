package types

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the contract needed for AccountKeeper related APIs.
// Interface provides support to use non-sdk AccountKeeper for AnteHandler's decorators.
type AccountKeeper interface {
	GetParams(ctx context.Context) (params types.Params)
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// FeegrantKeeper defines the expected feegrant keeper.
type FeegrantKeeper interface {
	UseGrantedFees(ctx context.Context, granter, grantee sdk.AccAddress, fee sdk.Coins, msgs []sdk.Msg) error
}

// FeegrantKeeper defines the expected feegrant keeper.
type GuardKeeper interface {
	GetAdmin(ctx sdk.Context) sdk.AccAddress
	CheckIsAdmin(ctx sdk.Context, address string) error
}

type LiquidityKeeper interface {
	GetSwapAmount(ctx sdk.Context, pairId uint64, swapCoin sdk.Coin) (offerCoin sdk.Coin, price sdkmath.LegacyDec, err error)
}

type TxfeesKeeper interface {
	GetParams(ctx sdk.Context) (params Params)
	GetFeeToken(
		ctx sdk.Context,
		denom string,
	) (val FeeToken, found bool)
}
