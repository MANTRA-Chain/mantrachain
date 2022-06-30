package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type NftCollectionCategory string
type NftCollectionDisplayTheme string

const (
	GeneralNftCollectionCat      NftCollectionCategory = "general"
	ArtNftCollectionCat          NftCollectionCategory = "art"
	CollectiblesNftCollectionCat                       = "collectibles"
	MusicNftCollectionCat                              = "music"
	PhotographyNftCollectionCat                        = "photography"
	SportsNftCollectionCat                             = "sports"
	TradingCardsNftCollectionCat                       = "tradingCards"
	UtilityNftCollectionCat                            = "utility"
)

func ValidateNftCollectionCategory(cat NftCollectionCategory) error {
	switch cat {
	case GeneralNftCollectionCat, ArtNftCollectionCat, CollectiblesNftCollectionCat, MusicNftCollectionCat,
		PhotographyNftCollectionCat, SportsNftCollectionCat, TradingCardsNftCollectionCat, UtilityNftCollectionCat:
		return nil
	default:
		return errors.Wrapf(ErrInvalidNftCollectionCategory, "invalid nft collection category: %s", cat)
	}
}
