package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/guard module sentinel errors
var (
	ErrInvalidTokenCollectionCreator         = sdkerrors.Register(ModuleName, 1111, "token collection creator is invalid")
	ErrInvalidTokenCollectionId              = sdkerrors.Register(ModuleName, 1112, "token collection id is invalid")
	ErrTokenNftNotFound                      = sdkerrors.Register(ModuleName, 1113, "token nft not found")
	ErrIncorrectNftOwner                     = sdkerrors.Register(ModuleName, 1114, "incorrect nft owner")
	ErrTokenNftAttributesIncorrectOrNotFound = sdkerrors.Register(ModuleName, 1115, "incorrect nft attributes or not found")
	ErrAccPermNotFound                       = sdkerrors.Register(ModuleName, 1116, "account permission not found")
	ErrAccPermCatIncorrectOrNotFound         = sdkerrors.Register(ModuleName, 1117, "incorrect account permission category or not found")
	ErrInvalidGuardTransfer                  = sdkerrors.Register(ModuleName, 1118, "invalid transfer enabled")
	ErrInvalidTokenCollectionCreatorParam    = sdkerrors.Register(ModuleName, 1119, "invalid token collection creator param")
	ErrInvalidTokenCollectionIdParam         = sdkerrors.Register(ModuleName, 1120, "invalid token collection id param")
	ErrGuardTransferNotFound                 = sdkerrors.Register(ModuleName, 1121, "transfer enabled not found")
	ErrAdminAccountParamMismatch             = sdkerrors.Register(ModuleName, 1122, "admin account param mismatch")
)
