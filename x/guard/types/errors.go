package types

import "cosmossdk.io/errors"

// DONTCOVER

// x/guard module sentinel errors
var (
	ErrInvalidTokenCollectionCreator                       = errors.Register(ModuleName, 1111, "token collection creator is invalid")
	ErrInvalidTokenCollectionId                            = errors.Register(ModuleName, 1112, "token collection id is invalid")
	ErrIncorrectNftOwner                                   = errors.Register(ModuleName, 1113, "incorrect nft owner")
	ErrAccountPrivilegesNotFound                           = errors.Register(ModuleName, 1114, "account privileges not found")
	ErrInvalidAccountPrivilegesTokenCollectionCreatorParam = errors.Register(ModuleName, 1115, "invalid account privileges token collection creator param")
	ErrInvalidAccountPrivilegesTokenCollectionIdParam      = errors.Register(ModuleName, 1116, "invalid account privileges token collection id param")
	ErrInsufficientPrivileges                              = errors.Register(ModuleName, 1117, "insufficient privileges")
)
