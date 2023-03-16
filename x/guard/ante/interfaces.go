package ante

import sdk "github.com/cosmos/cosmos-sdk/types"

type GuardKeeper interface {
	CheckHasPerm(ctx sdk.Context, address string) error
}
