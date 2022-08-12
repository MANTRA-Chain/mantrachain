package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
)

var _ types.QueryServer = Keeper{}
