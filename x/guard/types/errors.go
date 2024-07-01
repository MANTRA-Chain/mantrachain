package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/guard module sentinel errors
var (
	ErrInvalidSigner                                       = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidTokenCollectionCreator                       = errors.Register(ModuleName, 1101, "token collection creator is invalid")
	ErrInvalidTokenCollectionId                            = errors.Register(ModuleName, 1102, "token collection id is invalid")
	ErrMissingSoulBondNft                                  = errors.Register(ModuleName, 1103, "missing soul bond nft")
	ErrAccountPrivilegesNotFound                           = errors.Register(ModuleName, 1104, "account privileges not found")
	ErrInvalidAccountPrivilegesTokenCollectionCreatorParam = errors.Register(ModuleName, 1105, "invalid account privileges token collection creator param")
	ErrInvalidAccountPrivilegesTokenCollectionIdParam      = errors.Register(ModuleName, 1106, "invalid account privileges token collection id param")
	ErrInsufficientPrivileges                              = errors.Register(ModuleName, 1107, "insufficient privileges")
	ErrInvalidPrivileges                                   = errors.Register(ModuleName, 1108, "invalid privileges")
	ErrRequiredPrivilegesNotFound                          = errors.Register(ModuleName, 1109, "required privileges not found")
	ErrCoinRequiredPrivilegesNotFound                      = errors.Register(ModuleName, 1110, "coin required privileges not found")
	ErrCoinAdminNotFound                                   = errors.Register(ModuleName, 1111, "coin admin not found")
	ErrInvalidDenom                                        = errors.Register(ModuleName, 1112, "invalid denom")
	ErrCoinRequiredPrivilegesNotSet                        = errors.Register(ModuleName, 1113, "coin required privileges not set")
	ErrAccountRequiredPrivilegesNotSet                     = errors.Register(ModuleName, 1114, "account required privileges not set")
	ErrAccountRequiredPrivilegesNotFound                   = errors.Register(ModuleName, 1115, "account required privileges not found")
	ErrCoinsRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 1116, "coins required privileges not found")
	ErrAuthzRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 1117, "authz required privileges not found")
	ErrAuthzRequiredPrivilegesNotSet                       = errors.Register(ModuleName, 1118, "authz required privileges not set")
)
