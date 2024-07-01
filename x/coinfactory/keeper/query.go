package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

var _ types.QueryServer = Keeper{}
