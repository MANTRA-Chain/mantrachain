package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/rewards module sentinel errors
var (
	ErrSample               = sdkerrors.Register(ModuleName, 1801, "sample error")
	ErrProviderNotFound     = sdkerrors.Register(ModuleName, 1802, "provider not found")
	ErrProviderPairNotFound = sdkerrors.Register(ModuleName, 1803, "provider pair not found")
	ErrProviderPoolNotFound = sdkerrors.Register(ModuleName, 1804, "provider pool not found")
	ErrPairNotFound         = sdkerrors.Register(ModuleName, 1805, "pair not found")
	ErrInvalidPairId        = sdkerrors.Register(ModuleName, 1806, "invalid pair id")
	ErrSnapshotNotFound     = sdkerrors.Register(ModuleName, 1807, "snapshot not found")
	ErrBalanceMismatch      = sdkerrors.Register(ModuleName, 1808, "balance mismatch")
)
