package keeper

import (
	"github.com/LimeChain/mantrachain/x/bridge/types"
)

var _ types.QueryServer = Keeper{}
