package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

var _ types.QueryServer = Keeper{}
