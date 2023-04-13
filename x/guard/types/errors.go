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
	ErrCoinLocked                                          = errors.Register(ModuleName, 1121, "coin locked")
)
