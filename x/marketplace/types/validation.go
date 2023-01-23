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

func ValidateMarketplaceId(validMarketplaceId string, marketplaceId string) error {
	if strings.TrimSpace(marketplaceId) == "" {
		return errors.Wrapf(ErrInvalidMarketplaceId, "invalid marketplace, id: %s", marketplaceId)
	}

	if validMarketplaceId == "" {
		return errors.Wrap(ErrInvalidMarketplaceId, "missing marketplace id regex param")
	}

	if !regexp.MustCompile(validMarketplaceId).MatchString(marketplaceId) {
		return errors.Wrap(ErrInvalidMarketplaceId, marketplaceId)
	}

	return nil
}
