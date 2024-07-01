package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/did/types"
)

var _ types.QueryServer = Keeper{}
