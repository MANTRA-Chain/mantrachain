package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bridge module sentinel errors
var (
	ErrNoInputs            = sdkerrors.Register(ModuleName, 2, "no inputs")
	ErrMultipleSenders     = sdkerrors.Register(ModuleName, 3, "multiple senders")
	ErrNoOutputs           = sdkerrors.Register(ModuleName, 4, "no outputs")
	ErrInputOutputMismatch = sdkerrors.Register(ModuleName, 5, "input output mismatch")
)
