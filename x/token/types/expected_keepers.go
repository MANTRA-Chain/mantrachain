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
