package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the contract needed for AccountKeeper related APIs.
// Interface provides support to use non-sdk AccountKeeper for AnteHandler's decorators.
type AccountKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// FeegrantKeeper defines the expected feegrant keeper.
type FeegrantKeeper interface {
	UseGrantedFees(ctx sdk.Context, granter, grantee sdk.AccAddress, fee sdk.Coins, msgs []sdk.Msg) error
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	// Methods imported from bank should be defined here
}

// FeegrantKeeper defines the expected feegrant keeper.
type GuardKeeper interface {
	GetAdmin(ctx sdk.Context) sdk.AccAddress
	CheckIsAdmin(ctx sdk.Context, address string) error
}

type LiquidityKeeper interface {
	GetSwapAmount(ctx sdk.Context, pairId uint64, swapCoin sdk.Coin) (offerCoin sdk.Coin, price sdk.Dec, err error)
}

type TxfeesKeeper interface {
	GetParams(ctx sdk.Context) (params Params)
	GetFeeToken(
		ctx sdk.Context,
		denom string,
	) (val FeeToken, found bool)
}
