package keeper

import (
	"mantrachain/x/guard/types"
)

var _ types.QueryServer = Keeper{}
