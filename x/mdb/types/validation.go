package types

import (
	"regexp"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateNftCollectionId(validNftCollectionId string, collectionId string, isValidNftCollectionId func(s string) bool) error {
	if strings.TrimSpace(collectionId) == "" {
		return sdkerrors.Wrapf(ErrInvalidNftCollectionId, "invalid nft collection id %s", collectionId)
	}

	if validNftCollectionId == "" && isValidNftCollectionId == nil {
		return sdkerrors.Wrap(ErrInvalidNftCollectionId, "missing nft collection id regex param")
	}

	if isValidNftCollectionId == nil {
		isValidNftCollectionId = regexp.MustCompile(validNftCollectionId).MatchString
	}

	if !isValidNftCollectionId(collectionId) {
		return sdkerrors.Wrap(ErrInvalidNftCollectionId, collectionId)
	}

	return nil
}

func ValidateNftId(validNftId string, nftId string, isValidNftId func(s string) bool) error {
	if strings.TrimSpace(nftId) == "" {
		return sdkerrors.Wrapf(ErrInvalidNftId, "invalid nft id %s", nftId)
	}

	if validNftId == "" && isValidNftId == nil {
		return sdkerrors.Wrap(ErrInvalidNftId, "missing nft id regex param")
	}

	if isValidNftId == nil {
		isValidNftId = regexp.MustCompile(validNftId).MatchString
	}

	if !isValidNftId(nftId) {
		return sdkerrors.Wrapf(ErrInvalidNftId, "invalid nft id %s", nftId)
	}

	return nil
}
