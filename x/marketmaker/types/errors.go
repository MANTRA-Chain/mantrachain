package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// marketmaker module sentinel errors
var (
	ErrAlreadyExistMarketMaker = sdkerrors.Register(ModuleName, 1701, "already exist market maker")
	ErrEmptyClaimableIncentive = sdkerrors.Register(ModuleName, 1702, "empty claimable incentives")
	ErrNotExistMarketMaker     = sdkerrors.Register(ModuleName, 1703, "not exist market maker")
	ErrInvalidPairId           = sdkerrors.Register(ModuleName, 1704, "invalid pair id")
	ErrUnregisteredPairId      = sdkerrors.Register(ModuleName, 1705, "unregistered pair id")
	ErrInvalidDeposit          = sdkerrors.Register(ModuleName, 1706, "invalid apply deposit")
	ErrInvalidInclusion        = sdkerrors.Register(ModuleName, 1707, "invalid inclusion, already eligible")
	ErrInvalidExclusion        = sdkerrors.Register(ModuleName, 1708, "invalid exclusion, not eligible")
	ErrInvalidRejection        = sdkerrors.Register(ModuleName, 1709, "invalid rejection, already eligible")
	ErrNotEligibleMarketMaker  = sdkerrors.Register(ModuleName, 1710, "invalid distribution, not eligible")
)
