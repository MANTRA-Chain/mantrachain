package keeper

import (
	"github.com/AumegaChain/aumega/x/rewards/types"
)

var _ types.QueryServer = Keeper{}
