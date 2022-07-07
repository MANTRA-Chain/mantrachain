package types

import (
	nfttypes "github.com/LimeChain/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
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
	FilterNotOwn(ctx sdk.Context, classID string, nftIDs []string) (list []string)
	GetClassesByIds(ctx sdk.Context, classesIds []string) (classes []nfttypes.Class)
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
	TransferBatch(ctx sdk.Context, classID string, nftIDs []string, receiver sdk.AccAddress) error
}

type MnsKeeper interface {
	// Methods imported from mns should be defined here
}

type DidKeeper interface {
	SetNewDidDocument(ctx sdk.Context, id string, signer sdk.Address, pubKeyHex string, pubKeyType string) (string, error)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}
