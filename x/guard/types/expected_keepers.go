package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool
}

type AuthzKeeper interface {
	SaveGrant(ctx sdk.Context, grantee, granter sdk.AccAddress, authorization authz.Authorization, expiration *time.Time) error
	DeleteGrant(ctx sdk.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType string) error
}

type NFTKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
}

type TokenKeeper interface {
	HasRestrictedNftsCollection(ctx sdk.Context, index []byte) bool
}

type CoinFactoryKeeper interface {
	HasAdmin(ctx sdk.Context, denom string) bool
	GetAdmin(ctx sdk.Context, denom string) (sdk.AccAddress, bool)
}