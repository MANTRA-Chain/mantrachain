package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/token module sentinel errors
var (
	ErrKeyFormatNotSupported            = sdkerrors.Register(ModuleName, 1111, "key format not supported")
	ErrNftCollectionAlreadyExists       = sdkerrors.Register(ModuleName, 1112, "nft collection already exists")
	ErrInvalidNftCollectionCategory     = sdkerrors.Register(ModuleName, 1113, "invalid nft collection category")
	ErrUnauthorized                     = sdkerrors.Register(ModuleName, 1114, "unauthorized")
	ErrNftCollectionDoesNotExist        = sdkerrors.Register(ModuleName, 1115, "nft collection does not exists")
	ErrInvalidNftCollectionId           = sdkerrors.Register(ModuleName, 1116, "nft collection id provided is invalid")
	ErrInvalidNftCollectionName         = sdkerrors.Register(ModuleName, 1117, "nft collection name provided is invalid")
	ErrInvalidNftCollectionUrl          = sdkerrors.Register(ModuleName, 1118, "nft collection url provided is invalid")
	ErrInvalidNftCollectionDescription  = sdkerrors.Register(ModuleName, 1119, "nft collection description provided is invalid")
	ErrInvalidNftCollectionSymbol       = sdkerrors.Register(ModuleName, 1120, "invalid nft collection symbol")
	ErrInvalidNftCollectionOption       = sdkerrors.Register(ModuleName, 1121, "nft collection option provided is invalid")
	ErrInvalidNftCollectionImage        = sdkerrors.Register(ModuleName, 1122, "nft collection image provided is invalid")
	ErrInvalidNftCollectionLink         = sdkerrors.Register(ModuleName, 1123, "nft collection link provided is invalid")
	ErrInvalidNftCollectionOptionsCount = sdkerrors.Register(ModuleName, 1124, "nft collection options count provided is invalid")
	ErrInvalidNftCollectionImagesCount  = sdkerrors.Register(ModuleName, 1125, "nft collection images count provided is invalid")
	ErrInvalidNftCollectionLinksCount   = sdkerrors.Register(ModuleName, 1126, "nft collection links count provided is invalid")
	ErrInvalidNftCollectionOpened       = sdkerrors.Register(ModuleName, 1127, "nft collection opened provided is invalid")
	ErrInvalidNftCollectionSoulBonded   = sdkerrors.Register(ModuleName, 1128, "nft collection soul bonded provided is invalid")

	ErrInvalidNftId              = sdkerrors.Register(ModuleName, 1129, "nft id provided is invalid")
	ErrInvalidNftTitle           = sdkerrors.Register(ModuleName, 1130, "nft title provided is invalid")
	ErrInvalidNftUrl             = sdkerrors.Register(ModuleName, 1131, "nft url provided is invalid")
	ErrInvalidNftDescription     = sdkerrors.Register(ModuleName, 1132, "nft description provided is invalid")
	ErrInvalidNftImage           = sdkerrors.Register(ModuleName, 1133, "nft image provided is invalid")
	ErrInvalidNftLink            = sdkerrors.Register(ModuleName, 1134, "nft link provided is invalid")
	ErrInvalidNftAttribute       = sdkerrors.Register(ModuleName, 1135, "nft attribute provided is invalid")
	ErrInvalidNftImagesCount     = sdkerrors.Register(ModuleName, 1136, "nft images count provided is invalid")
	ErrInvalidNftLinksCount      = sdkerrors.Register(ModuleName, 1137, "nft links count provided is invalid")
	ErrInvalidNftAttributesCount = sdkerrors.Register(ModuleName, 1138, "nft attributes count provided is invalid")
	ErrInvalidNftsCount          = sdkerrors.Register(ModuleName, 1139, "nfts count provided is invalid")
	ErrInvalidNft                = sdkerrors.Register(ModuleName, 1140, "nfts provided is invalid")

	ErrApproveNftDisabled   = sdkerrors.Register(ModuleName, 1141, "approve nft disabled")
	ErrApproveNftsDisabled  = sdkerrors.Register(ModuleName, 1142, "approve nfts disabled")
	ErrTransferNftDisabled  = sdkerrors.Register(ModuleName, 1143, "transfer nft disabled")
	ErrTransferNftsDisabled = sdkerrors.Register(ModuleName, 1144, "transfer nfts disabled")
)
