package types

import (
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NFTKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
}

type TokenKeeper interface {
	GetNftCollection(ctx sdk.Context, creator sdk.AccAddress, index []byte) (val tokentypes.NftCollection, found bool)
	HasNftCollection(ctx sdk.Context, creator sdk.AccAddress, index []byte) (exists bool)
	GetNft(ctx sdk.Context, collectionIndex []byte, nftIndex []byte) (val tokentypes.Nft, found bool)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type VaultKeeper interface {
	UpsertNftStake(ctx sdk.Context, marketplaceIndex []byte, collectionIndex []byte, index []byte, creator sdk.AccAddress, amount sdk.Coin, delegate bool) error
}
