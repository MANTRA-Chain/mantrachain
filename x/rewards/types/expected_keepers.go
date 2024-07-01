package types

import (
	"context"

	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Acc

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type LiquidityKeeper interface {
	GetPair(ctx sdk.Context, id uint64) (pair liquiditytypes.Pair, found bool)
	GetAllPairsIds(ctx sdk.Context) (pairsIds []uint64)
	GetPool(ctx sdk.Context, id uint64) (pool liquiditytypes.Pool, found bool)
	SetHooks(gh liquiditytypes.LiquidityHooks)
}

type GuardKeeper interface {
	GetAdmin(ctx sdk.Context) sdk.AccAddress
	AddTransferAccAddressesWhitelist(addresses []string) []string
	RemoveTransferAccAddressesWhitelist(addresses []string)
	CheckCanTransferCoins(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error
}
