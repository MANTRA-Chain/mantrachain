package types

import "cosmossdk.io/errors"

// DONTCOVER

// x/token module sentinel errors
var (
	ErrKeyFormatNotSupported              = errors.Register(ModuleName, 1111, "key format not supported")
	ErrNftCollectionAlreadyExists         = errors.Register(ModuleName, 1112, "nft collection already exists")
	ErrInvalidNftCollectionCategory       = errors.Register(ModuleName, 1113, "invalid nft collection category")
	ErrUnauthorized                       = errors.Register(ModuleName, 1114, "unauthorized")
	ErrNftCollectionDoesNotExist          = errors.Register(ModuleName, 1115, "nft collection does not exists")
	ErrInvalidNftCollectionId             = errors.Register(ModuleName, 1116, "nft collection id provided is invalid")
	ErrInvalidNftCollectionName           = errors.Register(ModuleName, 1117, "nft collection name provided is invalid")
	ErrInvalidNftCollectionUrl            = errors.Register(ModuleName, 1118, "nft collection url provided is invalid")
	ErrInvalidNftCollectionDescription    = errors.Register(ModuleName, 1119, "nft collection description provided is invalid")
	ErrInvalidNftCollectionSymbol         = errors.Register(ModuleName, 1120, "invalid nft collection symbol")
	ErrInvalidNftCollectionOption         = errors.Register(ModuleName, 1121, "nft collection option provided is invalid")
	ErrInvalidNftCollectionImage          = errors.Register(ModuleName, 1122, "nft collection image provided is invalid")
	ErrInvalidNftCollectionLink           = errors.Register(ModuleName, 1123, "nft collection link provided is invalid")
	ErrInvalidNftCollectionOptionsCount   = errors.Register(ModuleName, 1124, "nft collection options count provided is invalid")
	ErrInvalidNftCollectionImagesCount    = errors.Register(ModuleName, 1125, "nft collection images count provided is invalid")
	ErrInvalidNftCollectionLinksCount     = errors.Register(ModuleName, 1126, "nft collection links count provided is invalid")
	ErrInvalidNftCollectionOpened         = errors.Register(ModuleName, 1127, "nft collection opened provided is invalid")
	ErrInvalidNftCollectionRestricted     = errors.Register(ModuleName, 1128, "nft collection restricted provided is invalid")
	ErrInvalidNftCollectionSoulBondedNfts = errors.Register(ModuleName, 1129, "nft collection soul bonded nfts provided is invalid")

	ErrInvalidNftId              = errors.Register(ModuleName, 1130, "nft id provided is invalid")
	ErrInvalidNftTitle           = errors.Register(ModuleName, 1131, "nft title provided is invalid")
	ErrInvalidNftUrl             = errors.Register(ModuleName, 1132, "nft url provided is invalid")
	ErrInvalidNftDescription     = errors.Register(ModuleName, 1133, "nft description provided is invalid")
	ErrInvalidNftImage           = errors.Register(ModuleName, 1134, "nft image provided is invalid")
	ErrInvalidNftLink            = errors.Register(ModuleName, 1135, "nft link provided is invalid")
	ErrInvalidNftAttribute       = errors.Register(ModuleName, 1136, "nft attribute provided is invalid")
	ErrInvalidNftImagesCount     = errors.Register(ModuleName, 1137, "nft images count provided is invalid")
	ErrInvalidNftLinksCount      = errors.Register(ModuleName, 1138, "nft links count provided is invalid")
	ErrInvalidNftAttributesCount = errors.Register(ModuleName, 1139, "nft attributes count provided is invalid")
	ErrInvalidNftsCount          = errors.Register(ModuleName, 1140, "nfts count provided is invalid")
	ErrInvalidNft                = errors.Register(ModuleName, 1141, "nfts provided is invalid")

	ErrApproveNftDisabled                       = errors.Register(ModuleName, 1142, "approve nft disabled")
	ErrApproveNftsDisabled                      = errors.Register(ModuleName, 1143, "approve nfts disabled")
	ErrTransferNftDisabled                      = errors.Register(ModuleName, 1144, "transfer nft disabled")
	ErrTransferNftsDisabled                     = errors.Register(ModuleName, 1145, "transfer nfts disabled")
	ErrApproveSoulBondedNftNotSupported         = errors.Register(ModuleName, 1146, "approve soul bonded nft not supported")
	ErrApproveSoulBondedNftsNotSupported        = errors.Register(ModuleName, 1147, "approve soul bonded nfts not supported")
	ErrTransferSoulBondedNftNotSupported        = errors.Register(ModuleName, 1148, "transfer soul bonded nft not supported")
	ErrTransferSoulBondedNftsNotSupported       = errors.Register(ModuleName, 1149, "transfer soul bonded nfts not supported")
	ErrNftModuleTransferNftDisabled             = errors.Register(ModuleName, 1150, "nft module transfer nft disabled")
	ErrSoulBondedNftCollectionOperationDisabled = errors.Register(ModuleName, 1151, "soul bonded nft collection operation disabled")
)
