package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// DONTCOVER

// x/token module sentinel errors
var (
	ErrNftCollectionAlreadyExists         = errors.Register(ModuleName, 2, "nft collection already exists")
	ErrInvalidNftCollectionCategory       = errors.Register(ModuleName, 3, "invalid nft collection category")
	ErrUnauthorized                       = errors.Register(ModuleName, 4, "unauthorized")
	ErrNftCollectionDoesNotExist          = errors.Register(ModuleName, 5, "nft collection does not exists")
	ErrInvalidNftCollectionId             = errors.Register(ModuleName, 6, "nft collection id provided is invalid")
	ErrInvalidNftCollectionName           = errors.Register(ModuleName, 7, "nft collection name provided is invalid")
	ErrInvalidNftCollectionUrl            = errors.Register(ModuleName, 8, "nft collection url provided is invalid")
	ErrInvalidNftCollectionDescription    = errors.Register(ModuleName, 9, "nft collection description provided is invalid")
	ErrInvalidNftCollectionSymbol         = errors.Register(ModuleName, 10, "invalid nft collection symbol")
	ErrInvalidNftCollectionOption         = errors.Register(ModuleName, 11, "nft collection option provided is invalid")
	ErrInvalidNftCollectionImage          = errors.Register(ModuleName, 12, "nft collection image provided is invalid")
	ErrInvalidNftCollectionLink           = errors.Register(ModuleName, 13, "nft collection link provided is invalid")
	ErrInvalidNftCollectionOptionsCount   = errors.Register(ModuleName, 14, "nft collection options count provided is invalid")
	ErrInvalidNftCollectionImagesCount    = errors.Register(ModuleName, 15, "nft collection images count provided is invalid")
	ErrInvalidNftCollectionLinksCount     = errors.Register(ModuleName, 16, "nft collection links count provided is invalid")
	ErrInvalidNftCollectionOpened         = errors.Register(ModuleName, 17, "nft collection opened provided is invalid")
	ErrInvalidNftCollectionRestricted     = errors.Register(ModuleName, 18, "nft collection restricted provided is invalid")
	ErrInvalidNftCollectionSoulBondedNfts = errors.Register(ModuleName, 19, "nft collection soul bonded nfts provided is invalid")

	ErrInvalidNftId              = errors.Register(ModuleName, 20, "nft id provided is invalid")
	ErrInvalidNftTitle           = errors.Register(ModuleName, 21, "nft title provided is invalid")
	ErrInvalidNftUrl             = errors.Register(ModuleName, 22, "nft url provided is invalid")
	ErrInvalidNftDescription     = errors.Register(ModuleName, 23, "nft description provided is invalid")
	ErrInvalidNftImage           = errors.Register(ModuleName, 24, "nft image provided is invalid")
	ErrInvalidNftLink            = errors.Register(ModuleName, 25, "nft link provided is invalid")
	ErrInvalidNftAttribute       = errors.Register(ModuleName, 26, "nft attribute provided is invalid")
	ErrInvalidNftImagesCount     = errors.Register(ModuleName, 27, "nft images count provided is invalid")
	ErrInvalidNftLinksCount      = errors.Register(ModuleName, 28, "nft links count provided is invalid")
	ErrInvalidNftAttributesCount = errors.Register(ModuleName, 29, "nft attributes count provided is invalid")
	ErrInvalidNftsCount          = errors.Register(ModuleName, 30, "nfts count provided is invalid")
	ErrInvalidNft                = errors.Register(ModuleName, 31, "nfts provided is invalid")

	ErrApproveNftDisabled                       = errors.Register(ModuleName, 32, "approve nft disabled")
	ErrApproveNftsDisabled                      = errors.Register(ModuleName, 33, "approve nfts disabled")
	ErrTransferNftDisabled                      = errors.Register(ModuleName, 34, "transfer nft disabled")
	ErrTransferNftsDisabled                     = errors.Register(ModuleName, 35, "transfer nfts disabled")
	ErrApproveSoulBondedNftNotSupported         = errors.Register(ModuleName, 36, "approve soul bonded nft not supported")
	ErrApproveSoulBondedNftsNotSupported        = errors.Register(ModuleName, 37, "approve soul bonded nfts not supported")
	ErrTransferSoulBondedNftNotSupported        = errors.Register(ModuleName, 38, "transfer soul bonded nft not supported")
	ErrTransferSoulBondedNftsNotSupported       = errors.Register(ModuleName, 39, "transfer soul bonded nfts not supported")
	ErrSoulBondedNftCollectionOperationDisabled = errors.Register(ModuleName, 40, "soul bonded nft collection operation disabled")

	ErrInvalidNftImageIndex = errors.Register(ModuleName, 41, "nft image index provided is invalid")
	ErrInvalidDid           = errors.Register(ModuleName, 42, "did provided is invalid")
	ErrInvalidAccount       = errors.Register(ModuleName, 43, "account provided is invalid")

	ErrNftCollectionNotRestricted = errors.Register(ModuleName, 44, "nft collection is not restricted")
	ErrKeyFormatNotSupported      = errors.Register(ModuleName, 45, "key format not supported")
)
