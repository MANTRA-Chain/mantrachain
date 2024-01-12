package keeper

import (
	"github.com/AumegaChain/aumega/x/txfees/types"
)

var _ types.QueryServer = Keeper{}
