package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/mns module sentinel errors
var (
	ErrKeyFormatNotSupported = sdkerrors.Register(ModuleName, 1111, "key format not supported")
	ErrInvalidDomainType     = sdkerrors.Register(ModuleName, 1112, "invalid domain type")
	ErrInvalidDomainName     = sdkerrors.Register(ModuleName, 1113, "domain name provided is invalid")
	ErrDomainAlreadyExists   = sdkerrors.Register(ModuleName, 1114, "domain already exists")
	ErrDomainDoesNotExist    = sdkerrors.Register(ModuleName, 1115, "domain does not exist")
	ErrDomainExpired         = sdkerrors.Register(ModuleName, 1116, "domain has expired")
	ErrDomainClosed          = sdkerrors.Register(ModuleName, 1117, "domain closed")
)
