package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/nft module sentinel errors
var (
	ErrClassExists    = sdkerrors.Register(ModuleName, 1111, "nft class already exist")
	ErrClassNotExists = sdkerrors.Register(ModuleName, 1112, "nft class does not exist")
	ErrNFTExists      = sdkerrors.Register(ModuleName, 1113, "nft already exist")
	ErrNFTNotExists   = sdkerrors.Register(ModuleName, 1114, "nft does not exist")
	ErrInvalidID      = sdkerrors.Register(ModuleName, 1115, "invalid id")
	ErrInvalidClassID = sdkerrors.Register(ModuleName, 1116, "invalid class id")
)
