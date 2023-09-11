package types

import (
	"fmt"
	// this line is used by starport scaffolding # genesis/types/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 2

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DidDocuments: []GenesisDidDocument{},
	}
}

func (gs GenesisState) ValidateDidDocument() error {
	// Check for duplicated index in didDocument
	didDocumentIndexMap := make(map[string]struct{})

	for _, elem := range gs.DidDocuments {
		var key []byte
		key = append(key, DidDocumentKey...)
		key = append(key, Placeholder...)
		key = append(key, []byte(elem.DidDocument.Id)...)
		key = append(key, Placeholder...)

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
	return gs.ValidateDidDocument()
}
