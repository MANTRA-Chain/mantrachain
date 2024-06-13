package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/rewards module sentinel errors
var (
	ErrProviderNotFound     = sdkerrors.Register(ModuleName, 2, "provider not found")
	ErrProviderPairNotFound = sdkerrors.Register(ModuleName, 3, "provider pair not found")
	ErrProviderPoolNotFound = sdkerrors.Register(ModuleName, 4, "provider pool not found")
	ErrPairNotFound         = sdkerrors.Register(ModuleName, 5, "pair not found")
	ErrInvalidPairId        = sdkerrors.Register(ModuleName, 6, "invalid pair id")
	ErrSnapshotNotFound     = sdkerrors.Register(ModuleName, 7, "snapshot not found")
	ErrBalanceMismatch      = sdkerrors.Register(ModuleName, 8, "balance mismatch")
	ErrSnapshotPoolNotFound = sdkerrors.Register(ModuleName, 9, "snapshot pool not found")
)
