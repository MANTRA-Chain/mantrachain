package keeper

import (
	"mantrachain/x/txfees/types"
)

var _ types.QueryServer = Keeper{}
