package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
)

var _ types.QueryServer = Keeper{}
