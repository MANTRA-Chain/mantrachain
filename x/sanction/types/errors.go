package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/sanction module sentinel errors
var (
	ErrAccountBlacklisted = errorsmod.Register(ModuleName, 500, "Account blacklisted")
	ErrNotSanctioner      = errorsmod.Register(ModuleName, 501, "signer is not a sanctioner")
	ErrRoleNotPermitted   = errorsmod.Register(ModuleName, 502, "sanctioner role does not permit this action")
	ErrProtectedAddress   = errorsmod.Register(ModuleName, 503, "address cannot be sanctioned by a sanctioner")
	ErrGovernanceEntry    = errorsmod.Register(ModuleName, 504, "entry was added by governance and only governance can remove it")
	ErrTooManyAddresses   = errorsmod.Register(ModuleName, 505, "too many addresses in one message")
	ErrReasonTooLong      = errorsmod.Register(ModuleName, 506, "reason is too long")
	ErrInvalidSanctioner  = errorsmod.Register(ModuleName, 507, "invalid sanctioner")
)
