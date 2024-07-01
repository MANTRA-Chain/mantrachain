package types

import (
	"context"
	"time"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
}

type AuthzKeeper interface {
	SaveGrant(ctx context.Context, grantee, granter sdk.AccAddress, authorization authz.Authorization, expiration *time.Time) error
	DeleteGrant(ctx context.Context, grantee sdk.AccAddress, granter sdk.AccAddress, msgType string) error
}

type NFTKeeper interface {
	GetOwner(ctx context.Context, classID string, nftID string) sdk.AccAddress
}

type TokenKeeper interface {
	HasRestrictedNftsCollection(ctx sdk.Context, index []byte) bool
}

type CoinFactoryKeeper interface {
	GetAdmin(ctx sdk.Context, denom string) (sdk.AccAddress, bool)
}
