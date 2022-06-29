package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/mdb module sentinel errors
var (
	ErrKeyFormatNotSupported               = sdkerrors.Register(ModuleName, 1111, "key format not supported")
	ErrNftCollectionAlreadyExists          = sdkerrors.Register(ModuleName, 1112, "nft collection already exists")
	ErrInvalidNftCollectionCategory        = sdkerrors.Register(ModuleName, 1113, "invalid nft collection category")
	ErrUnauthorized                        = sdkerrors.Register(ModuleName, 1114, "unauthorized")
	ErrNftCollectionDoesNotExist           = sdkerrors.Register(ModuleName, 1115, "nft collection does not exists")
	ErrInvalidNftCollectionId              = sdkerrors.Register(ModuleName, 1116, "nft collection id provided is invalid")
	ErrInvalidNftCollectionName            = sdkerrors.Register(ModuleName, 1117, "nft collection name provided is invalid")
	ErrInvalidNftCollectionUrl             = sdkerrors.Register(ModuleName, 1118, "nft collection url provided is invalid")
	ErrInvalidNftCollectionDescription     = sdkerrors.Register(ModuleName, 1119, "nft collection description provided is invalid")
	ErrInvalidNftCollectionDisplayTheme    = sdkerrors.Register(ModuleName, 1120, "invalid nft collection display theme")
	ErrInvalidNftCollectionSymbol          = sdkerrors.Register(ModuleName, 1121, "invalid nft collection symbol")
	ErrInvalidNftCollectionCreatorEarnings = sdkerrors.Register(ModuleName, 1122, "nft collection creator earnings provided is invalid")
	ErrInvalidNftCollectionImage           = sdkerrors.Register(ModuleName, 1123, "nft collection image provided is invalid")
	ErrInvalidNftCollectionLink            = sdkerrors.Register(ModuleName, 1124, "nft collection link provided is invalid")
	ErrInvalidNftCollectionOpened          = sdkerrors.Register(ModuleName, 1125, "nft collection opened provided is invalid")

	ErrInvalidNftId          = sdkerrors.Register(ModuleName, 1126, "nft id provided is invalid")
	ErrInvalidNftTitle       = sdkerrors.Register(ModuleName, 1127, "nft title provided is invalid")
	ErrInvalidNftUrl         = sdkerrors.Register(ModuleName, 1128, "nft url provided is invalid")
	ErrInvalidNftDescription = sdkerrors.Register(ModuleName, 1129, "nft description provided is invalid")
	ErrInvalidNftImage       = sdkerrors.Register(ModuleName, 1130, "nft image provided is invalid")
	ErrInvalidNftLink        = sdkerrors.Register(ModuleName, 1131, "nft link provided is invalid")
	ErrInvalidNftAttribute   = sdkerrors.Register(ModuleName, 1132, "nft attribute provided is invalid")
	ErrInvalidNftsLength     = sdkerrors.Register(ModuleName, 1133, "nfts length provided is invalid")
)
