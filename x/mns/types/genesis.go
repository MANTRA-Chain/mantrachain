package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DomainList:     []Domain{},
		DomainNameList: []DomainName{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in domain
	domainIndexMap := make(map[string]struct{})

	for _, elem := range gs.DomainList {
		index := string(GetDomainIndex(elem.Domain))
		if _, ok := domainIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for domain")
		}
		domainIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in domainName
	domainNameIndexMap := make(map[string]struct{})

	for _, elem := range gs.DomainNameList {
		index := string(GetDomainNameIndex(elem.Domain, elem.DomainName))
		if _, ok := domainNameIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for domainName")
		}
		domainNameIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
