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
		NftCollectionList:            []NftCollection{},
		NftList:                      []Nft{},
		SoulBondedNftsCollectionList: [][]byte{},
		RestrictedNftsCollectionList: [][]byte{},
		OpenedNftsCollectionList:     [][]byte{},
		NftCollectionOwnerList:       []*NftCollectionOwner{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

func (gs GenesisState) ValidateNftCollection() error {
	// Check for duplicated index in nftCollection
	nftCollectionIndexMap := make(map[string]struct{})

	for _, elem := range gs.NftCollectionList {
		var key []byte
		indexBytes := []byte(NftCollectionStoreKey(elem.Creator))
		key = append(key, indexBytes...)
		key = append(key, Placeholder...)
		key = append(key, elem.Index...)
		key = append(key, Placeholder...)

		index := string(key)
		if _, ok := nftCollectionIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for nftCollection")
		}
		nftCollectionIndexMap[index] = struct{}{}
	}

	return nil
}

func (gs GenesisState) ValidateNft() error {
	// Check for duplicated index in nft
	nftIndexMap := make(map[string]struct{})

	for _, elem := range gs.NftList {
		var key []byte
		indexBytes := []byte(NftStoreKey(elem.CollectionIndex))
		key = append(key, indexBytes...)
		key = append(key, Placeholder...)
		key = append(key, elem.Index...)
		key = append(key, Placeholder...)

		index := string(key)
		if _, ok := nftIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for nft")
		}
		nftIndexMap[index] = struct{}{}
	}

	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateNftCollection()
	if err != nil {
		return err
	}

	err = gs.ValidateNft()
	if err != nil {
		return err
	}

	// Check for duplicated index in soulBondedNftsCollection
	soulBondedNftsCollectionIndexMap := make(map[string]struct{})

	for _, elem := range gs.SoulBondedNftsCollectionList {
		index := string(elem)
		if _, ok := soulBondedNftsCollectionIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for soulBondedNftsCollection")
		}
		soulBondedNftsCollectionIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in restrictedNftsCollection
	restrictedNftsCollectionIndexMap := make(map[string]struct{})

	for _, elem := range gs.RestrictedNftsCollectionList {
		index := string(elem)
		if _, ok := restrictedNftsCollectionIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for restrictedNftsCollection")
		}
		restrictedNftsCollectionIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in openedNftsCollection
	openedNftsCollectionIndexMap := make(map[string]struct{})

	for _, elem := range gs.OpenedNftsCollectionList {
		index := string(elem)
		if _, ok := openedNftsCollectionIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for openedNftsCollection")
		}
		openedNftsCollectionIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in nftCollectionOwner
	nftCollectionOwnerIndexMap := make(map[string]struct{})

	for _, elem := range gs.NftCollectionOwnerList {
		index := string(elem.Index)
		if _, ok := nftCollectionOwnerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for nftCollectionOwner")
		}
		nftCollectionOwnerIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
