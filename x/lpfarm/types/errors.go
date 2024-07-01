package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/lpfarm module sentinel errors
var (
	ErrInvalidSigner         = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrPlanAlreadyTerminated = errors.Register(ModuleName, 1101, "plan is already terminated")
)
