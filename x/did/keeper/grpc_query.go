package keeper

import (
	"github.com/LimeChain/mantrachain/x/did/types"
)

var _ types.QueryServer = Keeper{}
