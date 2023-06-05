package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardKeeper interface {
	CheckIsAdmin(ctx sdk.Context, address string) error
}
