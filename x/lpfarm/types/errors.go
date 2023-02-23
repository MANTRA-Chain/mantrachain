package types

import "cosmossdk.io/errors"

var (
	ErrPlanAlreadyTerminated = errors.Register(ModuleName, 2, "plan is already terminated")
)
