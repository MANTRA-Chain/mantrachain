package types

import (
	"context"

	nft "cosmossdk.io/x/nft"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type NFTKeeper interface {
	SaveClass(ctx context.Context, class nft.Class) error
	Mint(ctx context.Context, token nft.NFT, receiver sdk.AccAddress) error
	BatchMint(ctx context.Context, tokens []nft.NFT, receiver sdk.AccAddress) error
	GetNFT(ctx context.Context, classID, nftID string) (nft.NFT, bool)
	GetOwner(ctx context.Context, classID string, nftID string) sdk.AccAddress
	Burn(ctx context.Context, classID string, nftID string) error
	BatchBurn(ctx context.Context, classID string, nftIDs []string) error
	Transfer(ctx context.Context, classID string, nftID string, receiver sdk.AccAddress) error
	BatchTransfer(ctx context.Context, classID string, nftIDs []string, receiver sdk.AccAddress) error
}

type DidKeeper interface {
	CreateNewDidDocument(ctx sdk.Context, id string, controller string) (string, error)
	ForceRemoveDidDocumentIfExists(ctx sdk.Context, id string) (bool, error)
}

type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
	CheckNewRestrictedNftsCollection(ctx sdk.Context, restrictedNftsCollection bool, address string) error
	CheckRestrictedNftsCollection(ctx sdk.Context, collectionCreator string, collectionId string, address string) error
	GetAccountPrivilegesTokenCollectionCreatorAndCollectionId(ctx sdk.Context) (string, string)
}
