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

func ValidateNftCollectionId(validNftCollectionId string, collectionId string, isValidNftCollectionId func(s string) bool) error {
	if strings.TrimSpace(collectionId) == "" {
		return errors.Wrapf(ErrInvalidNftCollectionId, "invalid nft collection id %s", collectionId)
	}

	if validNftCollectionId == "" && isValidNftCollectionId == nil {
		return errors.Wrap(ErrInvalidNftCollectionId, "missing nft collection id regex param")
	}

	if isValidNftCollectionId == nil {
		isValidNftCollectionId = regexp.MustCompile(validNftCollectionId).MatchString
	}

	if !isValidNftCollectionId(collectionId) {
		return errors.Wrap(ErrInvalidNftCollectionId, collectionId)
	}

	return nil
}

func ValidateNftId(validNftId string, nftId string, isValidNftId func(s string) bool) error {
	if strings.TrimSpace(nftId) == "" {
		return errors.Wrapf(ErrInvalidNftId, "invalid nft id %s", nftId)
	}

	if validNftId == "" && isValidNftId == nil {
		return errors.Wrap(ErrInvalidNftId, "missing nft id regex param")
	}

	if isValidNftId == nil {
		isValidNftId = regexp.MustCompile(validNftId).MatchString
	}

	if !isValidNftId(nftId) {
		return errors.Wrapf(ErrInvalidNftId, "invalid nft id %s", nftId)
	}

	return nil
}