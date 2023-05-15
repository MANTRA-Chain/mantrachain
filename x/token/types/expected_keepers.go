package types

import (
	nfttypes "github.com/MANTRA-Finance/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NFTKeeper interface {
	SaveClass(ctx sdk.Context, class nfttypes.Class) error
	Mint(ctx sdk.Context, token nfttypes.NFT, receiver sdk.AccAddress) error
	BatchMint(ctx sdk.Context, tokens []nfttypes.NFT, receiver sdk.AccAddress) error
	GetNFT(ctx sdk.Context, classID, nftID string) (nfttypes.NFT, bool)
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	Burn(ctx sdk.Context, classID string, nftID string) error
	BatchBurn(ctx sdk.Context, classID string, nftIDs []string) error
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
	BatchTransfer(ctx sdk.Context, classID string, nftIDs []string, receiver sdk.AccAddress) error
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	// Methods imported from bank should be defined here
}
type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
	CheckNewRestrictedNftsCollection(ctx sdk.Context, restrictedNftsCollection bool, address string) error
	CheckRestrictedNftsCollection(ctx sdk.Context, collectionCreator string, collectionId string, address string) error
	GetAccountPrivilegesTokenCollectionCreatorAndCollectionId(ctx sdk.Context) (string, string)
}
