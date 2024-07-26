package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
	AddTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress) []sdk.AccAddress
	RemoveTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress)
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}
