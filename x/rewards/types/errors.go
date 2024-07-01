package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/rewards module sentinel errors
var (
	ErrInvalidSigner        = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrProviderNotFound     = errors.Register(ModuleName, 1101, "provider not found")
	ErrProviderPairNotFound = errors.Register(ModuleName, 1102, "provider pair not found")
	ErrProviderPoolNotFound = errors.Register(ModuleName, 1103, "provider pool not found")
	ErrPairNotFound         = errors.Register(ModuleName, 1104, "pair not found")
	ErrInvalidPairId        = errors.Register(ModuleName, 1105, "invalid pair id")
	ErrSnapshotNotFound     = errors.Register(ModuleName, 1106, "snapshot not found")
	ErrBalanceMismatch      = errors.Register(ModuleName, 1107, "balance mismatch")
	ErrSnapshotPoolNotFound = errors.Register(ModuleName, 1108, "snapshot pool not found")
)
