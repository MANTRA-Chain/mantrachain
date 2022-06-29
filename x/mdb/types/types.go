package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type NftCollectionCategory string
type NftCollectionDisplayTheme string

const (
	ArtNftCollectionCat          NftCollectionCategory = "art"
	CollectiblesNftCollectionCat                       = "collectibles"
	MusicNftCollectionCat                              = "music"
	PhotographyNftCollectionCat                        = "photography"
	SportsNftCollectionCat                             = "sports"
	TradingCardsNftCollectionCat                       = "tradingCards"
	UtilityNftCollectionCat                            = "utility"
)

const (
	PaddedNftCollectionTheme    NftCollectionDisplayTheme = "art"
	ContainedNftCollectionTheme                           = "contained"
	CoveredNftCollectionTheme                             = "covered"
)

func ValidateNftCollectionCategory(cat NftCollectionCategory) error {
	switch cat {
	case ArtNftCollectionCat, CollectiblesNftCollectionCat, MusicNftCollectionCat,
		PhotographyNftCollectionCat, SportsNftCollectionCat, TradingCardsNftCollectionCat, UtilityNftCollectionCat:
		return nil
	default:
		return errors.Wrapf(ErrInvalidNftCollectionCategory, "invalid nft collection category: %s", cat)
	}
}

func ValidateNftCollectionDisplayTheme(theme NftCollectionDisplayTheme) error {
	switch theme {
	case PaddedNftCollectionTheme, ContainedNftCollectionTheme, CoveredNftCollectionTheme:
		return nil
	default:
		return errors.Wrapf(ErrInvalidNftCollectionDisplayTheme, "invalid nft collection display_theme: %s", theme)
	}
}
