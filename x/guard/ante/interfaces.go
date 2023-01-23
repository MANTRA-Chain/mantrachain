package ante

import (
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardKeeper interface {
	CheckCanTransfer(ctx sdk.Context, tokenKeeper TokenKeeper, nftKeeper NFTKeeper, addresses []sdk.AccAddress, amount sdk.Coins) (bool, error)
}

type TokenKeeper interface {
	GetNft(
		ctx sdk.Context,
		collectionIndex []byte,
		nftIndex []byte,
	) (val tokentypes.Nft, found bool)
}

type NFTKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
}
