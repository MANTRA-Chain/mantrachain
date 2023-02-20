package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	nft "github.com/cosmos/cosmos-sdk/x/nft"
)

type NFTKeeper interface {
	SaveClass(ctx sdk.Context, class nft.Class) error
	GetClass(ctx sdk.Context, classID string) (nft.Class, bool)
	Mint(ctx sdk.Context, token nft.NFT, receiver sdk.AccAddress) error
	BatchMint(ctx sdk.Context, tokens []nft.NFT, receiver sdk.AccAddress) error
	GetNFT(ctx sdk.Context, classID, nftID string) (nft.NFT, bool)
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	Burn(ctx sdk.Context, classID string, nftID string) error
	BatchBurn(ctx sdk.Context, classID string, nftIDs []string) error
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
	BatchTransfer(ctx sdk.Context, classID string, nftIDs []string, receiver sdk.AccAddress) error
	GetBalance(ctx sdk.Context, classID string, owner sdk.AccAddress) uint64
	GetTotalSupply(ctx sdk.Context, classID string) uint64
	GetNFTsOfClassByOwner(ctx sdk.Context, classID string, owner sdk.AccAddress) (nfts []nft.NFT)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	// Methods imported from bank should be defined here
}
