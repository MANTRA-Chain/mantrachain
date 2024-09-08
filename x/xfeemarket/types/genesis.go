package types

import "errors"

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:           DefaultParams(),
		DenomMultipliers: []DenomMultiplier{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	for _, denomMuliplier := range gs.DenomMultipliers {
		if err := denomMuliplier.Validate(); err != nil {
			return err
		}
	}
	return gs.Params.Validate()
}

func (m DenomMultiplier) Validate() error {
	if m.Denom == "" {
		return errors.New("denom cannot be empty")
	}
	return nil
}
