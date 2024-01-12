package types

import (
	liquiditytypes "github.com/MANTRA-Finance/aumega/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type LiquidityKeeper interface {
	GetPair(ctx sdk.Context, id uint64) (pair liquiditytypes.Pair, found bool)
	GetAllPairsIds(ctx sdk.Context) (pairsIds []uint64)
}

type GuardKeeper interface {
	GetAdmin(ctx sdk.Context) sdk.AccAddress
	WhitelistTransferAccAddresses(addresses []string, isWhitelisted bool) []string
	CheckCanTransferCoins(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error
}
