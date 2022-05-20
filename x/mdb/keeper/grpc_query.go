package keeper

import (
	"github.com/LimeChain/mantrachain/x/mdb/types"
)

var _ types.QueryServer = Keeper{}
