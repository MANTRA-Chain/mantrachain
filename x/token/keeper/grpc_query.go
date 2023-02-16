package keeper

import (
	"github.com/LimeChain/mantrachain/x/token/types"
)

var _ types.QueryServer = Keeper{}
