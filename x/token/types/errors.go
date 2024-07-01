package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/token module sentinel errors
var (
	ErrInvalidSigner                      = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrNftCollectionAlreadyExists         = errors.Register(ModuleName, 1101, "nft collection already exists")
	ErrInvalidNftCollectionCategory       = errors.Register(ModuleName, 1102, "invalid nft collection category")
	ErrUnauthorized                       = errors.Register(ModuleName, 1103, "unauthorized")
	ErrNftCollectionDoesNotExist          = errors.Register(ModuleName, 1104, "nft collection does not exists")
	ErrInvalidNftCollectionId             = errors.Register(ModuleName, 1105, "nft collection id provided is invalid")
	ErrInvalidNftCollectionName           = errors.Register(ModuleName, 1106, "nft collection name provided is invalid")
	ErrInvalidNftCollectionUrl            = errors.Register(ModuleName, 1107, "nft collection url provided is invalid")
	ErrInvalidNftCollectionDescription    = errors.Register(ModuleName, 1108, "nft collection description provided is invalid")
	ErrInvalidNftCollectionSymbol         = errors.Register(ModuleName, 1109, "invalid nft collection symbol")
	ErrInvalidNftCollectionOption         = errors.Register(ModuleName, 1110, "nft collection option provided is invalid")
	ErrInvalidNftCollectionImage          = errors.Register(ModuleName, 1111, "nft collection image provided is invalid")
	ErrInvalidNftCollectionLink           = errors.Register(ModuleName, 1112, "nft collection link provided is invalid")
	ErrInvalidNftCollectionOptionsCount   = errors.Register(ModuleName, 1113, "nft collection options count provided is invalid")
	ErrInvalidNftCollectionImagesCount    = errors.Register(ModuleName, 1114, "nft collection images count provided is invalid")
	ErrInvalidNftCollectionLinksCount     = errors.Register(ModuleName, 1115, "nft collection links count provided is invalid")
	ErrInvalidNftCollectionOpened         = errors.Register(ModuleName, 1116, "nft collection opened provided is invalid")
	ErrInvalidNftCollectionRestricted     = errors.Register(ModuleName, 1117, "nft collection restricted provided is invalid")
	ErrInvalidNftCollectionSoulBondedNfts = errors.Register(ModuleName, 1118, "nft collection soul bonded nfts provided is invalid")

	ErrInvalidNftId              = errors.Register(ModuleName, 1119, "nft id provided is invalid")
	ErrInvalidNftTitle           = errors.Register(ModuleName, 1120, "nft title provided is invalid")
	ErrInvalidNftUrl             = errors.Register(ModuleName, 1121, "nft url provided is invalid")
	ErrInvalidNftDescription     = errors.Register(ModuleName, 1122, "nft description provided is invalid")
	ErrInvalidNftImage           = errors.Register(ModuleName, 1123, "nft image provided is invalid")
	ErrInvalidNftLink            = errors.Register(ModuleName, 1124, "nft link provided is invalid")
	ErrInvalidNftAttribute       = errors.Register(ModuleName, 1125, "nft attribute provided is invalid")
	ErrInvalidNftImagesCount     = errors.Register(ModuleName, 1126, "nft images count provided is invalid")
	ErrInvalidNftLinksCount      = errors.Register(ModuleName, 1127, "nft links count provided is invalid")
	ErrInvalidNftAttributesCount = errors.Register(ModuleName, 1128, "nft attributes count provided is invalid")
	ErrInvalidNftsCount          = errors.Register(ModuleName, 1129, "nfts count provided is invalid")
	ErrInvalidNft                = errors.Register(ModuleName, 1130, "nfts provided is invalid")

	ErrApproveNftDisabled                       = errors.Register(ModuleName, 1131, "approve nft disabled")
	ErrApproveNftsDisabled                      = errors.Register(ModuleName, 1132, "approve nfts disabled")
	ErrTransferNftDisabled                      = errors.Register(ModuleName, 1133, "transfer nft disabled")
	ErrTransferNftsDisabled                     = errors.Register(ModuleName, 1134, "transfer nfts disabled")
	ErrApproveSoulBondedNftNotSupported         = errors.Register(ModuleName, 1135, "approve soul bonded nft not supported")
	ErrApproveSoulBondedNftsNotSupported        = errors.Register(ModuleName, 1136, "approve soul bonded nfts not supported")
	ErrTransferSoulBondedNftNotSupported        = errors.Register(ModuleName, 1137, "transfer soul bonded nft not supported")
	ErrTransferSoulBondedNftsNotSupported       = errors.Register(ModuleName, 1138, "transfer soul bonded nfts not supported")
	ErrSoulBondedNftCollectionOperationDisabled = errors.Register(ModuleName, 1139, "soul bonded nft collection operation disabled")

	ErrInvalidNftImageIndex = errors.Register(ModuleName, 1140, "nft image index provided is invalid")
	ErrInvalidDid           = errors.Register(ModuleName, 1141, "did provided is invalid")
	ErrInvalidAccount       = errors.Register(ModuleName, 1142, "account provided is invalid")

	ErrNftCollectionNotRestricted = errors.Register(ModuleName, 1143, "nft collection is not restricted")
	ErrKeyFormatNotSupported      = errors.Register(ModuleName, 1144, "key format not supported")
)
