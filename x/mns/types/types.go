package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type DomainType string

const (
	// OpenDomain is the domain type in which an account owner is the only entity that can perform actions on the account
	OpenDomain DomainType = "open"
	// ClosedDomain is the domain type in which the domain owner has control over accounts too
	ClosedDomain = "closed"
)

func ValidateDomainType(typ DomainType) error {
	switch typ {
	case OpenDomain, ClosedDomain:
		return nil
	default:
		return errors.Wrapf(ErrInvalidDomainType, "invalid domain type: %s", typ)
	}
}
