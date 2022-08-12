package keeper

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
)

var _ types.QueryServer = Keeper{}
