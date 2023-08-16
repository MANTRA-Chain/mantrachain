package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
)

var _ types.QueryServer = Keeper{}
