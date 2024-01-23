package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
)

var _ types.QueryServer = Keeper{}
