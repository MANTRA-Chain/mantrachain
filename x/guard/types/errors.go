package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// DONTCOVER

// x/guard module sentinel errors
var (
	ErrInvalidTokenCollectionCreator                       = errors.Register(ModuleName, 1111, "token collection creator is invalid")
	ErrInvalidTokenCollectionId                            = errors.Register(ModuleName, 1112, "token collection id is invalid")
	ErrMissingSoulBondNft                                  = errors.Register(ModuleName, 1113, "missing soul bond nft")
	ErrAccountPrivilegesNotFound                           = errors.Register(ModuleName, 1114, "account privileges not found")
	ErrInvalidAccountPrivilegesTokenCollectionCreatorParam = errors.Register(ModuleName, 1115, "invalid account privileges token collection creator param")
	ErrInvalidAccountPrivilegesTokenCollectionIdParam      = errors.Register(ModuleName, 1116, "invalid account privileges token collection id param")
	ErrInsufficientPrivileges                              = errors.Register(ModuleName, 1117, "insufficient privileges")
	ErrInvalidPrivileges                                   = errors.Register(ModuleName, 1118, "invalid privileges")
	ErrRequiredPrivilegesNotFound                          = errors.Register(ModuleName, 1119, "required privileges not found")
	ErrCoinRequiredPrivilegesNotFound                      = errors.Register(ModuleName, 1120, "coin required privileges not found")
	ErrCoinAdminNotFound                                   = errors.Register(ModuleName, 1121, "coin admin not found")
	ErrInvalidDenom                                        = errors.Register(ModuleName, 1122, "invalid denom")
	ErrCoinRequiredPrivilegesNotSet                        = errors.Register(ModuleName, 1123, "coin required privileges not set")
	ErrAccountRequiredPrivilegesNotSet                     = errors.Register(ModuleName, 1124, "account required privileges not set")
	ErrAccountRequiredPrivilegesNotFound                   = errors.Register(ModuleName, 1125, "account required privileges not found")
	ErrCoinsRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 1126, "coins required privileges not found")
	ErrAuthzRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 1127, "authz required privileges not found")
	ErrAuthzRequiredPrivilegesNotSet                       = errors.Register(ModuleName, 1128, "authz required privileges not set")
)
