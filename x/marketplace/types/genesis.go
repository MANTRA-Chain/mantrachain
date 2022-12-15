package types

import (
	"fmt"
	// this line is used by starport scaffolding # genesis/types/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		MarketplaceList:           []Marketplace{},
		MarketplaceCollectionList: []MarketplaceCollection{},
		MarketplaceNftList:        []MarketplaceNft{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

func (gs GenesisState) ValidateMarketplace() error {
	// Check for duplicated index in marketplace
	marketplaceIndexMap := make(map[string]struct{})

	for _, elem := range gs.MarketplaceList {
		var key []byte
		indexBytes := []byte(MarketplaceStoreKey(elem.Creator))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Index...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := marketplaceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for marketplace")
		}
		marketplaceIndexMap[index] = struct{}{}
	}

	return nil
}

func (gs GenesisState) ValidateMarketplaceCollection() error {
	// Check for duplicated index in marketplace
	marketplaceIndexMap := make(map[string]struct{})

	for _, elem := range gs.MarketplaceCollectionList {
		var key []byte
		indexBytes := []byte(MarketplaceCollectionStoreKey(elem.MarketplaceIndex))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Index...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := marketplaceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for marketplace")
		}
		marketplaceIndexMap[index] = struct{}{}
	}

	return nil
}

func (gs GenesisState) ValidateMarketplaceNft() error {
	// Check for duplicated index in marketplaceNft
	marketplaceNftIndexMap := make(map[string]struct{})

	for _, elem := range gs.MarketplaceNftList {
		var key []byte
		indexBytes := []byte(MarketplaceNftStoreKey(elem.MarketplaceIndex, elem.CollectionIndex))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Index...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := marketplaceNftIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for marketplaceNft")
		}
		marketplaceNftIndexMap[index] = struct{}{}
	}

	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateMarketplace()
	if err != nil {
		return err
	}

	err = gs.ValidateMarketplaceCollection()
	if err != nil {
		return err
	}

	err = gs.ValidateMarketplaceNft()
	if err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
