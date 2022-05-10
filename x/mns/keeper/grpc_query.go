package keeper

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
)

var _ types.QueryServer = Keeper{}
