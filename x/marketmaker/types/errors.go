package types

import "cosmossdk.io/errors"

// marketmaker module sentinel errors
var (
	ErrAlreadyExistMarketMaker = errors.Register(ModuleName, 2, "already exist market maker")
	ErrEmptyClaimableIncentive = errors.Register(ModuleName, 3, "empty claimable incentives")
	ErrNotExistMarketMaker     = errors.Register(ModuleName, 4, "not exist market maker")
	ErrInvalidPairId           = errors.Register(ModuleName, 5, "invalid pair id")
	ErrUnregisteredPairId      = errors.Register(ModuleName, 6, "unregistered pair id")
	ErrInvalidDeposit          = errors.Register(ModuleName, 7, "invalid apply deposit")
	ErrInvalidInclusion        = errors.Register(ModuleName, 8, "invalid inclusion, already eligible")
	ErrInvalidExclusion        = errors.Register(ModuleName, 9, "invalid exclusion, not eligible")
	ErrInvalidRejection        = errors.Register(ModuleName, 10, "invalid rejection, already eligible")
	ErrNotEligibleMarketMaker  = errors.Register(ModuleName, 11, "invalid distribution, not eligible")
)
