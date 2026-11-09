package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NormalizeAddress parses a bech32 account address and returns its canonical
// string, so that mixed-case input can't create an entry the ante handler
// never matches.
func NormalizeAddress(addr string) (string, error) {
	acc, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return "", err
	}
	return acc.String(), nil
}

// IsValid reports whether the role can be assigned to a sanctioner.
func (r SanctionerRole) IsValid() bool {
	return r == SanctionerRole_SANCTIONER_ROLE_FULL || r == SanctionerRole_SANCTIONER_ROLE_ADD_ONLY
}

// CanUnsanction reports whether the role may remove sanctioned addresses and
// revoke add-only sanctioners.
func (r SanctionerRole) CanUnsanction() bool {
	return r == SanctionerRole_SANCTIONER_ROLE_FULL
}

// ValidateSanctioners checks a full sanctioner list: valid addresses, valid
// roles and no duplicates.
func ValidateSanctioners(sanctioners []Sanctioner) error {
	seen := make(map[string]struct{}, len(sanctioners))
	for _, s := range sanctioners {
		addr, err := NormalizeAddress(s.Address)
		if err != nil {
			return fmt.Errorf("invalid sanctioner address %q: %w", s.Address, err)
		}
		if !s.Role.IsValid() {
			return fmt.Errorf("sanctioner %s has invalid role %s", s.Address, s.Role)
		}
		if _, dup := seen[addr]; dup {
			return fmt.Errorf("duplicate sanctioner %s", s.Address)
		}
		seen[addr] = struct{}{}
	}
	return nil
}
