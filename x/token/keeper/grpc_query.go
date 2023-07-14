package keeper

import (
	"mantrachain/x/token/types"
)

var _ types.QueryServer = Keeper{}

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}
