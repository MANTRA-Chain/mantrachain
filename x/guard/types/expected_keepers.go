package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
}

type NFTKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
}

type CoinFactoryKeeper interface {
	HasAdmin(ctx sdk.Context, denom string) bool
	GetAdmin(ctx sdk.Context, denom string) (sdk.AccAddress, bool)
}
