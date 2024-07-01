package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
)

var _ types.QueryServer = Keeper{}
