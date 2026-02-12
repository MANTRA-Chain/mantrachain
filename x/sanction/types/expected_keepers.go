package types

import (
	"context"

	"cosmossdk.io/x/feegrant"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

type AuthzKeeper interface {
	GranterGrants(ctx context.Context, req *authz.QueryGranterGrantsRequest) (*authz.QueryGranterGrantsResponse, error)
	DeleteGrant(ctx context.Context, grantee, granter sdk.AccAddress, msgType string) error
}

type FeegrantKeeper interface {
	AllowancesByGranter(ctx context.Context, req *feegrant.QueryAllowancesByGranterRequest) (*feegrant.QueryAllowancesByGranterResponse, error)
	RevokeAllowance(ctx context.Context, msg *feegrant.MsgRevokeAllowance) (*feegrant.MsgRevokeAllowanceResponse, error)
}
