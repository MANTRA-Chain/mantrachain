package types

import "cosmossdk.io/errors"

// DONTCOVER

// x/guard module sentinel errors
var (
	ErrInvalidTokenCollectionCreator      = errors.Register(ModuleName, 1111, "token collection creator is invalid")
	ErrInvalidTokenCollectionId           = errors.Register(ModuleName, 1112, "token collection id is invalid")
	ErrTokenNftNotFound                   = errors.Register(ModuleName, 1113, "token nft not found")
	ErrIncorrectNftOwner                  = errors.Register(ModuleName, 1114, "incorrect nft owner")
	ErrAccPermNotFound                    = errors.Register(ModuleName, 1115, "account permission not found")
	ErrAccPermIdIncorrectOrNotFound       = errors.Register(ModuleName, 1116, "incorrect account permission id or not found")
	ErrInvalidGuardTransfer               = errors.Register(ModuleName, 1117, "invalid transfer enabled")
	ErrInvalidTokenCollectionCreatorParam = errors.Register(ModuleName, 1118, "invalid token collection creator param")
	ErrInvalidTokenCollectionIdParam      = errors.Register(ModuleName, 1119, "invalid token collection id param")
	ErrGuardTransferNotFound              = errors.Register(ModuleName, 1120, "transfer enabled not found")
	ErrAdminAccountParamMismatch          = errors.Register(ModuleName, 1121, "admin account param mismatch")
	ErrInsufficientPriviliges             = errors.Register(ModuleName, 1122, "insufficient priviliges")
)
