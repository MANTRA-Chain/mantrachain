package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AccountPrivilegesList:  []*AccountPrivileges{},
		GuardTransferCoins:     nil,
		RequiredPrivilegesList: []*RequiredPrivileges{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in accountPrivileges
	accountPrivilegesIndexMap := make(map[string]struct{})

	for _, elem := range gs.AccountPrivilegesList {
		index := string(elem.Account)
		if _, ok := accountPrivilegesIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for accountPrivileges")
		}
		accountPrivilegesIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in requiredPrivileges
	requiredPrivilegesIndexMap := make(map[string]struct{})

	for _, elem := range gs.RequiredPrivilegesList {
		var key []byte
		indexBytes := []byte(RequiredPrivilegesStoreKey([]byte(elem.Kind)))
		key = append(key, indexBytes...)
		key = append(key, Placeholder...)
		key = append(key, elem.Index...)
		key = append(key, Placeholder...)

		index := string(key)
		if _, ok := requiredPrivilegesIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for requiredPrivileges")
		}
		requiredPrivilegesIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
