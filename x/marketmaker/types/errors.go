package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/marketmaker module sentinel errors
var (
	ErrInvalidSigner           = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrAlreadyExistMarketMaker = errors.Register(ModuleName, 1101, "already exist market maker")
	ErrEmptyClaimableIncentive = errors.Register(ModuleName, 1102, "empty claimable incentives")
	ErrNotExistMarketMaker     = errors.Register(ModuleName, 1103, "not exist market maker")
	ErrInvalidPairId           = errors.Register(ModuleName, 1104, "invalid pair id")
	ErrUnregisteredPairId      = errors.Register(ModuleName, 1105, "unregistered pair id")
	ErrInvalidDeposit          = errors.Register(ModuleName, 1106, "invalid apply deposit")
	ErrInvalidInclusion        = errors.Register(ModuleName, 1107, "invalid inclusion, already eligible")
	ErrInvalidExclusion        = errors.Register(ModuleName, 1108, "invalid exclusion, not eligible")
	ErrInvalidRejection        = errors.Register(ModuleName, 1109, "invalid rejection, already eligible")
	ErrNotEligibleMarketMaker  = errors.Register(ModuleName, 1110, "invalid distribution, not eligible")
)
