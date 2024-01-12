package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// DONTCOVER

// x/guard module sentinel errors
var (
	ErrInvalidTokenCollectionCreator                       = errors.Register(ModuleName, 1401, "token collection creator is invalid")
	ErrInvalidTokenCollectionId                            = errors.Register(ModuleName, 1402, "token collection id is invalid")
	ErrMissingSoulBondNft                                  = errors.Register(ModuleName, 1403, "missing soul bond nft")
	ErrAccountPrivilegesNotFound                           = errors.Register(ModuleName, 1404, "account privileges not found")
	ErrInvalidAccountPrivilegesTokenCollectionCreatorParam = errors.Register(ModuleName, 1405, "invalid account privileges token collection creator param")
	ErrInvalidAccountPrivilegesTokenCollectionIdParam      = errors.Register(ModuleName, 1406, "invalid account privileges token collection id param")
	ErrInsufficientPrivileges                              = errors.Register(ModuleName, 1407, "insufficient privileges")
	ErrInvalidPrivileges                                   = errors.Register(ModuleName, 1408, "invalid privileges")
	ErrRequiredPrivilegesNotFound                          = errors.Register(ModuleName, 1409, "required privileges not found")
	ErrCoinRequiredPrivilegesNotFound                      = errors.Register(ModuleName, 1410, "coin required privileges not found")
	ErrCoinAdminNotFound                                   = errors.Register(ModuleName, 1411, "coin admin not found")
	ErrInvalidDenom                                        = errors.Register(ModuleName, 1412, "invalid denom")
	ErrCoinRequiredPrivilegesNotSet                        = errors.Register(ModuleName, 1413, "coin required privileges not set")
	ErrAccountRequiredPrivilegesNotSet                     = errors.Register(ModuleName, 1414, "account required privileges not set")
	ErrAccountRequiredPrivilegesNotFound                   = errors.Register(ModuleName, 1415, "account required privileges not found")
	ErrCoinsRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 1416, "coins required privileges not found")
	ErrAuthzRequiredPrivilegesNotFound                     = errors.Register(ModuleName, 1417, "authz required privileges not found")
	ErrAuthzRequiredPrivilegesNotSet                       = errors.Register(ModuleName, 1418, "authz required privileges not set")
)
