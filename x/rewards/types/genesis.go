package types

import (
"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		SnapshotList: []Snapshot{},
// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in snapshot
snapshotIdMap := make(map[uint64]bool)
snapshotCount := gs.GetSnapshotCount()
for _, elem := range gs.SnapshotList {
	if _, ok := snapshotIdMap[elem.Id]; ok {
		return fmt.Errorf("duplicated id for snapshot")
	}
	if elem.Id >= snapshotCount {
		return fmt.Errorf("snapshot id should be lower or equal than the last id")
	}
	snapshotIdMap[elem.Id] = true
}
// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
