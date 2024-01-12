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
	ErrPairNotFound         = sdkerrors.Register(ModuleName, 1804, "pair not found")
	ErrInvalidPairId        = sdkerrors.Register(ModuleName, 1805, "invalid pair id")
)
