package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/vault module sentinel errors
var (
	ErrKeyFormatNotSupported          = sdkerrors.Register(ModuleName, 1111, "key format not supported")
	ErrInvalidMarketplaceId           = sdkerrors.Register(ModuleName, 1112, "marketplace id provided is invalid")
	ErrInvalidCollectionId            = sdkerrors.Register(ModuleName, 1113, "collection id provided is invalid")
	ErrInvalidNftId                   = sdkerrors.Register(ModuleName, 1114, "nft id provided is invalid")
	ErrValidatorDoesNotExist          = sdkerrors.Register(ModuleName, 1115, "validator does not exists")
	ErrInvalidStakingValidatorAddress = sdkerrors.Register(ModuleName, 1116, "staking validator address is invalid")
	ErrNftStakeDoesNotExist           = sdkerrors.Register(ModuleName, 1117, "nft stake does not exists")
	ErrLastEpochBlockNotFound         = sdkerrors.Register(ModuleName, 1118, "last epoch block not found")
	ErrUnavailable                    = sdkerrors.Register(ModuleName, 1119, "unavailable")
	ErrUnauthorized                   = sdkerrors.Register(ModuleName, 1120, "unauthorized")
	ErrInvalidBlockStart              = sdkerrors.Register(ModuleName, 1121, "invalid block start")
	ErrInvalidChain                   = sdkerrors.Register(ModuleName, 1122, "invalid chain")
	ErrInvalidValidator               = sdkerrors.Register(ModuleName, 1123, "invalid validator")
	ErrAdminAccountParamMismatch      = sdkerrors.Register(ModuleName, 1124, "admin account param mismatch")
	ErrInvalidAdminAccount            = sdkerrors.Register(ModuleName, 1125, "admin account param is invalid")
	ErrBridgeDoesNotExist             = sdkerrors.Register(ModuleName, 1126, "bridge does not exists")
	ErrChainValidatorBridgeNotFound   = sdkerrors.Register(ModuleName, 1127, "chain validator bridge not found")
)
