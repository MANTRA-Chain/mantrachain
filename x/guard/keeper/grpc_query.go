package keeper

import (
	"github.com/LimeChain/mantrachain/x/guard/types"
)

var _ types.QueryServer = Keeper{}
