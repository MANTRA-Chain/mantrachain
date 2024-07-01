package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
)

var _ types.QueryServer = Keeper{}
