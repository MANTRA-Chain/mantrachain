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
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, blacklistAccount := range gs.BlacklistAccounts {
		_, err := sdk.AccAddressFromBech32(blacklistAccount)
		if err != nil {
			return fmt.Errorf("invalid account %s", blacklistAccount)
		}
	}

	if err := ValidateSanctioners(gs.Sanctioners); err != nil {
		return err
	}

	return gs.Params.Validate()
}
