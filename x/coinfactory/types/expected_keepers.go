package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type BankKeeper interface {
	// Methods imported from bank should be defined here
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)

	HasSupply(ctx context.Context, denom string) bool
	GetSupply(ctx context.Context, denom string) sdk.Coin

	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error

	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	SetModuleAccount(ctx context.Context, macc sdk.ModuleAccountI)
}

type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
	GetAdmin(ctx sdk.Context) sdk.AccAddress
	AddTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress) []sdk.AccAddress
	RemoveTransferAccAddressesWhitelist(ctx sdk.Context, addresses []sdk.AccAddress)
	CheckHasAuthz(ctx sdk.Context, address string, authz string) error
	CheckCanTransferCoins(ctx sdk.Context, address sdk.AccAddress, coins sdk.Coins) error
}

// CommunityPoolKeeper defines the contract needed to be fulfilled for community pool interactions.
type CommunityPoolKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
