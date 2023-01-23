package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DidDocumentList: []DidDocument{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

func (gs GenesisState) ValidateDidDocument() error {
	// Check for duplicated index in didDocument
	didDocumentIndexMap := make(map[string]struct{})

	for _, elem := range gs.DidDocumentList {
		var key []byte
		indexBytes := []byte(StoreKey)
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, DidDocumentKey...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Id...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := didDocumentIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for didDocument")
		}
		didDocumentIndexMap[index] = struct{}{}
	}

	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateDidDocument()
	if err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
