package types

import (
	"regexp"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/errors"
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

func ValidateNftCollectionId(validNftCollectionId string, collectionId string) error {
	if strings.TrimSpace(collectionId) == "" {
		return errors.Wrapf(ErrInvalidNftCollectionId, "invalid nft collection id %s", collectionId)
	}

	if validNftCollectionId == "" {
		return errors.Wrap(ErrInvalidNftCollectionId, "missing nft collection id regex param")
	}

	if !regexp.MustCompile(validNftCollectionId).MatchString(collectionId) {
		return errors.Wrap(ErrInvalidNftCollectionId, collectionId)
	}

	return nil
}

func ValidateNftId(validNftId string, nftId string) error {
	if strings.TrimSpace(nftId) == "" {
		return errors.Wrapf(ErrInvalidNftId, "invalid nft id %s", nftId)
	}

	if validNftId == "" {
		return errors.Wrap(ErrInvalidNftId, "missing nft id regex param")
	}

	if !regexp.MustCompile(validNftId).MatchString(nftId) {
		return errors.Wrapf(ErrInvalidNftId, "invalid nft id %s", nftId)
	}

	return nil
}
