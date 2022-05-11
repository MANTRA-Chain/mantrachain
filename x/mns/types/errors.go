package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/mns module sentinel errors
var (
	ErrKeyFormatNotSupported = sdkerrors.Register(ModuleName, 1111, "key format not supported")
)
