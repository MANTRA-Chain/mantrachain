package types

import (
	"fmt"

	_ "github.com/MANTRA-Chain/mantrachain/v8/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:            DefaultParams(),
		BlacklistAccounts: []string{},
		Sanctioners:       []Sanctioner{},
		SanctionRecords:   []SanctionedAddress{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	blacklisted := make(map[string]struct{}, len(gs.BlacklistAccounts))
	for _, blacklistAccount := range gs.BlacklistAccounts {
		_, err := sdk.AccAddressFromBech32(blacklistAccount)
		if err != nil {
			return fmt.Errorf("invalid account %s", blacklistAccount)
		}
		blacklisted[blacklistAccount] = struct{}{}
	}

	if err := ValidateSanctioners(gs.Sanctioners); err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(gs.SanctionRecords))
	for _, r := range gs.SanctionRecords {
		if _, ok := blacklisted[r.Address]; !ok {
			return fmt.Errorf("sanction record for %s, which is not in blacklist_accounts", r.Address)
		}
		if _, dup := seen[r.Address]; dup {
			return fmt.Errorf("duplicate sanction record for %s", r.Address)
		}
		seen[r.Address] = struct{}{}
		if r.Record == nil {
			return fmt.Errorf("sanction record for %s is empty", r.Address)
		}
		if _, err := sdk.AccAddressFromBech32(r.Record.Actor); err != nil {
			return fmt.Errorf("invalid actor %q in sanction record for %s", r.Record.Actor, r.Address)
		}
	}

	return gs.Params.Validate()
}
