package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// DONTCOVER

// x/guard module sentinel errors
var (
	ErrInvalidTokenCollectionCreator                       = errors.Register(ModuleName, 2, "token collection creator is invalid")
	ErrInvalidTokenCollectionId                            = errors.Register(ModuleName, 3, "token collection id is invalid")
	ErrMissingSoulBondNft                                  = errors.Register(ModuleName, 4, "missing soul bond nft")
	ErrAccountPrivilegesNotFound                           = errors.Register(ModuleName, 5, "account privileges not found")
	ErrInvalidAccountPrivilegesTokenCollectionCreatorParam = errors.Register(ModuleName, 6, "invalid account privileges token collection creator param")
	ErrInvalidAccountPrivilegesTokenCollectionIdParam      = errors.Register(ModuleName, 7, "invalid account privileges token collection id param")
	ErrInsufficientPrivileges                              = errors.Register(ModuleName, 8, "insufficient privileges")
	ErrInvalidPrivileges                                   = errors.Register(ModuleName, 9, "invalid privileges")
	ErrRequiredPrivilegesNotFound                          = errors.Register(ModuleName, 10, "required privileges not found")
	ErrCoinRequiredPrivilegesNotFound                      = errors.Register(ModuleName, 11, "coin required privileges not found")
	ErrCoinAdminNotFound                                   = errors.Register(ModuleName, 12, "coin admin not found")
	ErrInvalidDenom                                        = errors.Register(ModuleName, 13, "invalid denom")
	ErrCoinRequiredPrivilegesNotSet                        = errors.Register(ModuleName, 14, "coin required privileges not set")
	ErrAccountRequiredPrivilegesNotSet                     = errors.Register(ModuleName, 15, "account required privileges not set")
	ErrAccountRequiredPrivilegesNotFound                   = errors.Register(ModuleName, 16, "account required privileges not found")
	ErrCoinsRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 17, "coins required privileges not found")
	ErrAuthzRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 18, "authz required privileges not found")
	ErrAuthzRequiredPrivilegesNotSet                       = errors.Register(ModuleName, 19, "authz required privileges not set")
)
