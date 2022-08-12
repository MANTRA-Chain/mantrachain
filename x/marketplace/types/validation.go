package types

import (
	"regexp"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateMarketplaceEarningType(earningType MarketplaceEarningType) error {
	switch earningType {
	case Initially, Repetitive:
		return nil
	default:
		return errors.Wrapf(ErrInvalidMarketplaceEarningType, "invalid marketplace earning type: %s", earningType)
	}
}

func ValidateMarketplaceId(validMarketplaceId string, marketplaceId string, isValidMarketplaceId func(s string) bool) error {
	if strings.TrimSpace(marketplaceId) == "" {
		return errors.Wrapf(ErrInvalidMarketplaceId, "invalid marketplace id %s", marketplaceId)
	}

	if validMarketplaceId == "" && isValidMarketplaceId == nil {
		return errors.Wrap(ErrInvalidMarketplaceId, "missing marketplace id regex param")
	}

	if isValidMarketplaceId == nil {
		isValidMarketplaceId = regexp.MustCompile(validMarketplaceId).MatchString
	}

	if !isValidMarketplaceId(marketplaceId) {
		return errors.Wrap(ErrInvalidMarketplaceId, marketplaceId)
	}

	return nil
}
