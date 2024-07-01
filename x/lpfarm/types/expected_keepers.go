package types

import (
	"context"

	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected keeper interface of the bank module.
// Some methods are used only in simulation tests.
type BankKeeper interface {
	HasSupply(ctx context.Context, denom string) bool
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error

	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// LiquidityKeeper defines the expected keeper interface of the liquidity module.
type LiquidityKeeper interface {
	GetPair(ctx sdk.Context, id uint64) (pair liquiditytypes.Pair, found bool)
	GetAllPairs(ctx sdk.Context) (pairs []liquiditytypes.Pair)
	IteratePoolsByPair(ctx sdk.Context, pairId uint64, cb func(pool liquiditytypes.Pool) (stop bool, err error)) error
}

type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
	AddTransferAccAddressesWhitelist(addresses []string) []string
	RemoveTransferAccAddressesWhitelist(addresses []string)
}
