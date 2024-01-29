package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// DONTCOVER

// x/token module sentinel errors
var (
	ErrKeyFormatNotSupported              = errors.Register(ModuleName, 1901, "key format not supported")
	ErrNftCollectionAlreadyExists         = errors.Register(ModuleName, 1902, "nft collection already exists")
	ErrInvalidNftCollectionCategory       = errors.Register(ModuleName, 1903, "invalid nft collection category")
	ErrUnauthorized                       = errors.Register(ModuleName, 1904, "unauthorized")
	ErrNftCollectionDoesNotExist          = errors.Register(ModuleName, 1905, "nft collection does not exists")
	ErrInvalidNftCollectionId             = errors.Register(ModuleName, 1906, "nft collection id provided is invalid")
	ErrInvalidNftCollectionName           = errors.Register(ModuleName, 1907, "nft collection name provided is invalid")
	ErrInvalidNftCollectionUrl            = errors.Register(ModuleName, 1908, "nft collection url provided is invalid")
	ErrInvalidNftCollectionDescription    = errors.Register(ModuleName, 1909, "nft collection description provided is invalid")
	ErrInvalidNftCollectionSymbol         = errors.Register(ModuleName, 1910, "invalid nft collection symbol")
	ErrInvalidNftCollectionOption         = errors.Register(ModuleName, 1911, "nft collection option provided is invalid")
	ErrInvalidNftCollectionImage          = errors.Register(ModuleName, 1912, "nft collection image provided is invalid")
	ErrInvalidNftCollectionLink           = errors.Register(ModuleName, 1913, "nft collection link provided is invalid")
	ErrInvalidNftCollectionOptionsCount   = errors.Register(ModuleName, 1914, "nft collection options count provided is invalid")
	ErrInvalidNftCollectionImagesCount    = errors.Register(ModuleName, 1915, "nft collection images count provided is invalid")
	ErrInvalidNftCollectionLinksCount     = errors.Register(ModuleName, 1916, "nft collection links count provided is invalid")
	ErrInvalidNftCollectionOpened         = errors.Register(ModuleName, 1917, "nft collection opened provided is invalid")
	ErrInvalidNftCollectionRestricted     = errors.Register(ModuleName, 1918, "nft collection restricted provided is invalid")
	ErrInvalidNftCollectionSoulBondedNfts = errors.Register(ModuleName, 1919, "nft collection soul bonded nfts provided is invalid")

	ErrInvalidNftId              = errors.Register(ModuleName, 1920, "nft id provided is invalid")
	ErrInvalidNftTitle           = errors.Register(ModuleName, 1921, "nft title provided is invalid")
	ErrInvalidNftUrl             = errors.Register(ModuleName, 1922, "nft url provided is invalid")
	ErrInvalidNftDescription     = errors.Register(ModuleName, 1923, "nft description provided is invalid")
	ErrInvalidNftImage           = errors.Register(ModuleName, 1924, "nft image provided is invalid")
	ErrInvalidNftLink            = errors.Register(ModuleName, 1925, "nft link provided is invalid")
	ErrInvalidNftAttribute       = errors.Register(ModuleName, 1926, "nft attribute provided is invalid")
	ErrInvalidNftImagesCount     = errors.Register(ModuleName, 1927, "nft images count provided is invalid")
	ErrInvalidNftLinksCount      = errors.Register(ModuleName, 1928, "nft links count provided is invalid")
	ErrInvalidNftAttributesCount = errors.Register(ModuleName, 1929, "nft attributes count provided is invalid")
	ErrInvalidNftsCount          = errors.Register(ModuleName, 1930, "nfts count provided is invalid")
	ErrInvalidNft                = errors.Register(ModuleName, 1931, "nfts provided is invalid")

	ErrApproveNftDisabled                       = errors.Register(ModuleName, 1932, "approve nft disabled")
	ErrApproveNftsDisabled                      = errors.Register(ModuleName, 1933, "approve nfts disabled")
	ErrTransferNftDisabled                      = errors.Register(ModuleName, 1934, "transfer nft disabled")
	ErrTransferNftsDisabled                     = errors.Register(ModuleName, 1935, "transfer nfts disabled")
	ErrApproveSoulBondedNftNotSupported         = errors.Register(ModuleName, 1936, "approve soul bonded nft not supported")
	ErrApproveSoulBondedNftsNotSupported        = errors.Register(ModuleName, 1937, "approve soul bonded nfts not supported")
	ErrTransferSoulBondedNftNotSupported        = errors.Register(ModuleName, 1938, "transfer soul bonded nft not supported")
	ErrTransferSoulBondedNftsNotSupported       = errors.Register(ModuleName, 1939, "transfer soul bonded nfts not supported")
	ErrSoulBondedNftCollectionOperationDisabled = errors.Register(ModuleName, 1940, "soul bonded nft collection operation disabled")

	ErrInvalidNftImageIndex = errors.Register(ModuleName, 1941, "nft image index provided is invalid")
	ErrInvalidDid           = errors.Register(ModuleName, 1942, "did provided is invalid")
	ErrInvalidAccount       = errors.Register(ModuleName, 1943, "account provided is invalid")

	ErrNftCollectionNotRestricted = errors.Register(ModuleName, 1944, "nft collection is not restricted")
)
