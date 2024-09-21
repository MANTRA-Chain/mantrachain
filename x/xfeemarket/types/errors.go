package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/xfeemarket module sentinel errors
var (
	ErrDenomNotFound = errorsmod.Register(ModuleName, 1100, "denom not found:")
)
