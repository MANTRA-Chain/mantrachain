package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardKeeper interface {
	CheckCanTransfer(ctx sdk.Context, nftKeeper NFTKeeper, addresses []sdk.AccAddress, amount sdk.Coins) (bool, error)
}

type NFTKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
}
