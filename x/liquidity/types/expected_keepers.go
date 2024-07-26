package types

import (
	"context"

	guardtypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper is the expected account keeper
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper is the expected bank keeper
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx context.Context, denom string) sdk.Coin
	IterateTotalSupply(ctx context.Context, cb func(sdk.Coin) bool)
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
	CheckCanTransferCoins(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error
	AddTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress) []sdk.AccAddress
	RemoveTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress)
	SetWhitelistTransfersAccAddrs(ctx sdk.Context, whitelistTransfersAccAddrs guardtypes.WhitelistTransfersAccAddrs)
}

type LiquidityHooks interface {
	OnProvideLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) error
	OnWithdrawLiquidity(ctx sdk.Context, receiver sdk.Address, pairId uint64, poolId uint64, poolCoin sdk.Coin) error
}
