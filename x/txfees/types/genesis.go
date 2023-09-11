package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 5

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		FeeTokenList: []FeeToken{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in feeToken
	feeTokenDenomMap := make(map[string]struct{})

	for _, elem := range gs.FeeTokenList {
		index := string(FeeTokenKey(elem.Denom))
		if _, ok := feeTokenDenomMap[index]; ok {
			return fmt.Errorf("duplicated index for feeToken")
		}
		feeTokenDenomMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
