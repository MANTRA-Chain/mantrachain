package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:            DefaultParams(),
		BlacklistAccounts: []string{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	for _, blacklistAccount := range gs.BlacklistAccounts {
		_, err := sdk.AccAddressFromBech32(blacklistAccount)
		if err != nil {
			return fmt.Errorf("invalid account %s", blacklistAccount)
		}
	}
	return gs.Params.Validate()
}
