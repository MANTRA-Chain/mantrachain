package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TxHashList:   []TxHash{},
		BridgeList:   []Bridge{},
		Cw20Contract: nil,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

func (gs GenesisState) ValidateTxHash() error {
	// Check for duplicated index in txHash
	txHashIndexMap := make(map[string]struct{})

	for _, elem := range gs.TxHashList {
		var key []byte
		indexBytes := []byte(TxHashStoreKey(elem.BridgeIndex))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Index...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := txHashIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for txHash")
		}
		txHashIndexMap[index] = struct{}{}
	}

	return nil
}

func (gs GenesisState) ValidateBridge() error {
	// Check for duplicated index in bridge
	bridgeIndexMap := make(map[string]struct{})

	for _, elem := range gs.BridgeList {
		var key []byte
		indexBytes := []byte(BridgeStoreKey(elem.Creator))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Index...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := bridgeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for bridge")
		}
		bridgeIndexMap[index] = struct{}{}
	}

	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateBridge()
	if err != nil {
		return err
	}

	err = gs.ValidateTxHash()
	if err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
