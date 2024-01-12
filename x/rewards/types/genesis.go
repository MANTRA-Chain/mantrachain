package types

import (
	"encoding/binary"
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		SnapshotList:               []Snapshot{},
		ProviderList:               []Provider{},
		SnapshotCountList:          []SnapshotCount{},
		SnapshotStartIdList:        []SnapshotStartId{},
		DistributionPairsIds:       []byte{},
		PurgePairsIds:              []byte{},
		SnapshotsLastDistributedAt: 0,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in snapshot count
	snapshotIndexMap := make(map[string]struct{})

	snapshot := gs.GetSnapshotList()
	for _, elem := range snapshot {
		if elem.PairId == 0 {
			return fmt.Errorf("invalid pair id for snapshot")
		}
		var key []byte
		indexBytes := []byte(SnapshotStoreKey(elem.PairId))
		idBytes := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(idBytes, elem.Id)
		key = append(key, indexBytes...)
		key = append(key, Placeholder...)
		key = append(key, idBytes...)
		key = append(key, Placeholder...)

		index := string(key)
		if _, ok := snapshotIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for snapshot")
		}
		snapshotIndexMap[index] = struct{}{}
	}

	snapshotCount := gs.GetSnapshotCountList()
	snapshotCountMap := make(map[uint64]bool)
	for _, elem := range snapshotCount {
		if elem.PairId == 0 {
			return fmt.Errorf("invalid pair id for snapshot")
		}
		if _, ok := snapshotCountMap[elem.PairId]; ok {
			return fmt.Errorf("duplicated pair id for snapshot count")
		}
		snapshotCountMap[elem.PairId] = true
	}
	// Check for duplicated ID in snapshot start id
	snapshotStartId := gs.GetSnapshotStartIdList()
	snapshotStartIdMap := make(map[uint64]bool)
	for _, elem := range snapshotStartId {
		if elem.PairId == 0 {
			return fmt.Errorf("invalid pair id for snapshot")
		}
		if _, ok := snapshotStartIdMap[elem.PairId]; ok {
			return fmt.Errorf("duplicated pair id for snapshot start id")
		}
		snapshotStartIdMap[elem.PairId] = true
	}
	// Check for duplicated index in provider
	providerIndexMap := make(map[string]struct{})

	for _, elem := range gs.ProviderList {
		index := string(ProviderKey(elem.Index))
		if _, ok := providerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for provider")
		}
		providerIndexMap[index] = struct{}{}
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
