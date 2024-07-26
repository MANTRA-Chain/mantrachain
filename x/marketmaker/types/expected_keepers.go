package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected bank send keeper
type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// MintCoins is used only for simulation test codes
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
}

type GuardKeeper interface {
	AddTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress) []sdk.AccAddress
	RemoveTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress)
}
