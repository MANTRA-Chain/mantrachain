package types

import (
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NFTKeeper interface {
	SaveClass(ctx sdk.Context, class nfttypes.Class) error
	GetClass(ctx sdk.Context, classID string) (nfttypes.Class, bool)
	Mint(ctx sdk.Context, token nfttypes.NFT, receiver sdk.AccAddress) error
	MintBatch(ctx sdk.Context, tokens []nfttypes.NFT, receiver sdk.AccAddress) error
	GetNFT(ctx sdk.Context, classID, nftID string) (nfttypes.NFT, bool)
	GetNFTsByIds(ctx sdk.Context, classID string, nftIDs []string) (nfts []nfttypes.NFT)
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	Burn(ctx sdk.Context, classID string, nftID string) error
	BurnBatch(ctx sdk.Context, classID string, nftIDs []string) error
	FilterNotOwnNFTsIdsOfClass(ctx sdk.Context, classID string, nftIDs []string, owner sdk.AccAddress) (list []string)
	GetClassesByIds(ctx sdk.Context, classesIds []string) (classes []nfttypes.Class)
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
	TransferBatch(ctx sdk.Context, classID string, nftIDs []string, receiver sdk.AccAddress) error
	GetBalance(ctx sdk.Context, classID string, owner sdk.AccAddress) uint64
	GetTotalSupply(ctx sdk.Context, classID string) uint64
	GetNFTsOfClassByOwner(ctx sdk.Context, classID string, owner sdk.AccAddress) (nfts []nfttypes.NFT)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	// Methods imported from bank should be defined here
}
