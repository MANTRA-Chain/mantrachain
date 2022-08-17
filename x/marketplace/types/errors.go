package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/marketplace module sentinel errors
var (
	ErrKeyFormatNotSupported                                  = sdkerrors.Register(ModuleName, 1111, "key format not supported")
	ErrMarketplaceAlreadyExists                               = sdkerrors.Register(ModuleName, 1112, "marketplace already exists")
	ErrUnauthorized                                           = sdkerrors.Register(ModuleName, 1113, "unauthorized")
	ErrMarketplaceDoesNotExist                                = sdkerrors.Register(ModuleName, 1114, "marketplace does not exists")
	ErrInvalidMarketplaceId                                   = sdkerrors.Register(ModuleName, 1115, "marketplace id provided is invalid")
	ErrInvalidMarketplaceName                                 = sdkerrors.Register(ModuleName, 1116, "marketplace name provided is invalid")
	ErrInvalidMarketplaceUrl                                  = sdkerrors.Register(ModuleName, 1117, "marketplace url provided is invalid")
	ErrInvalidMarketplaceDescription                          = sdkerrors.Register(ModuleName, 1118, "marketplace description provided is invalid")
	ErrInvalidCollectionId                                    = sdkerrors.Register(ModuleName, 1119, "collection id provided is invalid")
	ErrCollectionDoesNotExist                                 = sdkerrors.Register(ModuleName, 1120, "collection does not exists")
	ErrInvalidInitiallyNftMinPrice                            = sdkerrors.Register(ModuleName, 1121, "initially nft min price provided is invalid")
	ErrInvalidNftsEarningsOnSaleMaxCount                      = sdkerrors.Register(ModuleName, 1122, "nfts earnings on sale count provided is invalid")
	ErrInvalidNftsEarningsOnSalePercentage                    = sdkerrors.Register(ModuleName, 1123, "nfts earnings on sale percentage provided is invalid")
	ErrInvalidNftsEarningsOnSaleAddress                       = sdkerrors.Register(ModuleName, 1124, "nfts earnings on sale address provided is invalid")
	ErrInvalidNftsEarningsOnYieldRewardMaxCount               = sdkerrors.Register(ModuleName, 1125, "nfts earnings on yield reward count provided is invalid")
	ErrInvalidNftsEarningsOnYieldRewardPercentage             = sdkerrors.Register(ModuleName, 1126, "nfts earnings on yield reward percentage provided is invalid")
	ErrInvalidNftsEarningsOnYieldRewardAddress                = sdkerrors.Register(ModuleName, 1127, "nfts earnings on yield reward address provided is invalid")
	ErrCollectionSettingsAlreadyExists                        = sdkerrors.Register(ModuleName, 1128, "collection settings already exists")
	ErrCollectionSettingsDoesNotExist                         = sdkerrors.Register(ModuleName, 1129, "collection settings does not exists")
	ErrInvalidInitiallyNftsVaultLockPercentage                = sdkerrors.Register(ModuleName, 1130, "initially nfts vault lock percentage provided is invalid")
	ErrInvalidInitiallyEarnAndLockSummaryPercentage           = sdkerrors.Register(ModuleName, 1131, "initially earn and lock summary percentage provided is invalid")
	ErrNftDoesNotExist                                        = sdkerrors.Register(ModuleName, 1132, "nft does not exists")
	ErrNftNotForSale                                          = sdkerrors.Register(ModuleName, 1133, "nft is not for sale")
	ErrInvalidNftBuyer                                        = sdkerrors.Register(ModuleName, 1134, "nft buyer provided is invalid")
	ErrInvalidMarketplaceNftsEarningsOnSaleEarningType        = sdkerrors.Register(ModuleName, 1135, "marketplace nfts earnings on sale earning type provided is invalid")
	ErrInvalidMarketplaceNftsEarningsOnYieldRewardEarningType = sdkerrors.Register(ModuleName, 1136, "marketplace nfts earnings on yield reward earning type provided is invalid")
	ErrInvalidMarketplaceEarningType                          = sdkerrors.Register(ModuleName, 1137, "marketplace earning type provided is invalid")
	ErrUnavailable                                            = sdkerrors.Register(ModuleName, 1138, "unavailable")
)
