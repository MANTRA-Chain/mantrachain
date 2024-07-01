package types

import fmt "fmt"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:       DefaultParams(),
		DidDocuments: []GenesisDidDocument{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	err := gs.ValidateDidDocument()
	if err != nil {
		return err
	}

	return gs.Params.Validate()
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
