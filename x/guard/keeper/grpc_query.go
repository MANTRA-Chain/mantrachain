package keeper

import (
	"github.com/AumegaChain/aumega/x/guard/types"
)

var _ types.QueryServer = Keeper{}
