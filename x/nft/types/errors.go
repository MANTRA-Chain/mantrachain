package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrClassExists                  = sdkerrors.Register(ModuleName, 1111, "nft class already exists")
	ErrClassNotExists               = sdkerrors.Register(ModuleName, 1112, "nft class does not exist")
	ErrNFTExists                    = sdkerrors.Register(ModuleName, 1113, "nft already exists")
	ErrNFTNotExists                 = sdkerrors.Register(ModuleName, 1114, "nft does not exist")
	ErrEmptyClassID                 = sdkerrors.Register(ModuleName, 1115, "empty class id")
	ErrEmptyNFTID                   = sdkerrors.Register(ModuleName, 1116, "empty nft id")
	ErrInvalidID                    = sdkerrors.Register(ModuleName, 1117, "invalid id")
	ErrInvalidClassID               = sdkerrors.Register(ModuleName, 1118, "invalid class id")
	ErrNftModuleTransferNftDisabled = sdkerrors.Register(ModuleName, 1119, "nft module transfer nft disabled")
)
