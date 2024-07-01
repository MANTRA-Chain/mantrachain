package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/farming/types"
)

var _ types.QueryServer = Keeper{}
