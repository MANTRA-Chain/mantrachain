package types

// DONTCOVER

import (
	fmt "fmt"

	"cosmossdk.io/errors"
)

// x/coinfactory module sentinel errors
var (
	ErrInvalidSigner            = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrDenomExists              = errors.Register(ModuleName, 1101, "attempting to create a denom that already exists (has bank metadata)")
	ErrUnauthorized             = errors.Register(ModuleName, 1102, "unauthorized account")
	ErrInvalidDenom             = errors.Register(ModuleName, 1103, "invalid denom")
	ErrInvalidCreator           = errors.Register(ModuleName, 1104, "invalid creator")
	ErrInvalidAuthorityMetadata = errors.Register(ModuleName, 1105, "invalid authority metadata")
	ErrInvalidGenesis           = errors.Register(ModuleName, 1106, "invalid genesis")
	ErrSubdenomTooLong          = errors.Register(ModuleName, 1107, fmt.Sprintf("subdenom too long, max length is %d bytes", MaxSubdenomLength))
	ErrCreatorTooLong           = errors.Register(ModuleName, 1108, fmt.Sprintf("creator too long, max length is %d bytes", MaxCreatorLength))
	ErrDenomDoesNotExist        = errors.Register(ModuleName, 1109, "denom does not exist")
	ErrBurnFromModuleAccount    = errors.Register(ModuleName, 1110, "cannot burn from module account")
)
