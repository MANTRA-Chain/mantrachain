package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bridge module sentinel errors
var (
	ErrInvalidAdminAccount        = sdkerrors.Register(ModuleName, 1111, "admin account is invalid")
	ErrInvalidBridgeId            = sdkerrors.Register(ModuleName, 1112, "bridge id is invalid")
	ErrUnauthorized               = sdkerrors.Register(ModuleName, 1113, "unauthorized")
	ErrBridgeDoesNotExist         = sdkerrors.Register(ModuleName, 1114, "bridge does not exist")
	ErrBridgeAlreadyExists        = sdkerrors.Register(ModuleName, 1115, "bridge already exists")
	ErrInvalidBridgeMetadata      = sdkerrors.Register(ModuleName, 1116, "bridge metadata is invalid")
	ErrBridgeAccountMismatch      = sdkerrors.Register(ModuleName, 1117, "bridge account mismatch")
	ErrInvalidMint                = sdkerrors.Register(ModuleName, 1118, "mint is invalid")
	ErrTxAlreadyProcessed         = sdkerrors.Register(ModuleName, 1119, "tx already processed")
	ErrInvalidTxHash              = sdkerrors.Register(ModuleName, 1120, "tx hash is invalid")
	ErrAdminAccountParamMismatch  = sdkerrors.Register(ModuleName, 1121, "admin account param mismatch")
	ErrInvalidStoreId             = sdkerrors.Register(ModuleName, 1122, "store id is invalid")
	ErrInvalidContractAddress     = sdkerrors.Register(ModuleName, 1123, "contract address is invalid")
	ErrInvalidPath                = sdkerrors.Register(ModuleName, 1124, "path is invalid")
	ErrInvalidCw20Name            = sdkerrors.Register(ModuleName, 1125, "cw20 name is invalid")
	ErrInvalidCw20Symbol          = sdkerrors.Register(ModuleName, 1126, "cw20 symbol is invalid")
	ErrInvalidCw20InitialBalances = sdkerrors.Register(ModuleName, 1127, "cw20 initial balances is invalid")
	ErrInvalidCw20Mint            = sdkerrors.Register(ModuleName, 1128, "cw20 mint is invalid")
	ErrInvalidCw20ContractAddress = sdkerrors.Register(ModuleName, 1129, "cw20 contract address is invalid")
)
