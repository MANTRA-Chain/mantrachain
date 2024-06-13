package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
)

var _ types.QueryServer = Keeper{}
