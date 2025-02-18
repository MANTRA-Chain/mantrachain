package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/sanction module sentinel errors
var (
	ErrAccountBlacklisted = errorsmod.Register(ModuleName, 500, "Account blacklisted")
)
