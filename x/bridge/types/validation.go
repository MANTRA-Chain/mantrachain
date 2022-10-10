package types

import (
	"regexp"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateBridgeId(validBridgeId string, bridgeId string) error {
	if strings.TrimSpace(bridgeId) == "" {
		return errors.Wrapf(ErrInvalidBridgeId, "invalid bridge id %s", bridgeId)
	}

	if validBridgeId == "" {
		return errors.Wrap(ErrInvalidBridgeId, "missing bridge id regex param")
	}

	if !regexp.MustCompile(validBridgeId).MatchString(bridgeId) {
		return errors.Wrap(ErrInvalidBridgeId, bridgeId)
	}

	return nil
}
