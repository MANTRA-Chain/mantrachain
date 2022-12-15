package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ChainValidatorBridgeList: []ChainValidatorBridge{},
		EpochList:                []Epoch{},
		LastEpochBlockList:       []LastEpochBlock{},
		NftStakeList:             []NftStake{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

func (gs GenesisState) ValidateChainValidatorBridge() error {
	// Check for duplicated index in chainValidatorBridge
	chainValidatorBridgeIndexMap := make(map[string]struct{})

	for _, elem := range gs.ChainValidatorBridgeList {
		var key []byte
		indexBytes := []byte(ChainValidatorBridgeKey(elem.Chain))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, GetChainValidatorBridgeIndex(elem.Validator)...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := chainValidatorBridgeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for chainValidatorBridge")
		}
		chainValidatorBridgeIndexMap[index] = struct{}{}
	}

	return nil
}

func (gs GenesisState) ValidateEpoch() error {
	// Check for duplicated index in epoch
	epochIndexMap := make(map[string]struct{})

	for _, elem := range gs.EpochList {
		var key []byte
		indexBytes := []byte(EpochStoreKey(elem.StakingChain))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Index...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := epochIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for epoch")
		}
		epochIndexMap[index] = struct{}{}
	}

	return nil
}

func (gs GenesisState) ValidateLastEpochBlock() error {
	// Check for duplicated index in lastEpochBlock
	lastEpochBlockIndexMap := make(map[string]struct{})

	for _, elem := range gs.LastEpochBlockList {
		var key []byte
		indexBytes := []byte(LastEpochBlockStoreKey(elem.StakingChain))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, GetLastEpochBlockIndex(elem.StakingValidator)...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := lastEpochBlockIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for lastEpochBlock")
		}
		lastEpochBlockIndexMap[index] = struct{}{}
	}

	return nil
}

func (gs GenesisState) ValidateNftStake() error {
	// Check for duplicated index in nftStake
	nftStakeIndexMap := make(map[string]struct{})

	for _, elem := range gs.NftStakeList {
		var key []byte
		indexBytes := []byte(NftStakeStoreKey(elem.MarketplaceIndex, elem.CollectionIndex))
		key = append(key, indexBytes...)
		key = append(key, []byte("/")...)
		key = append(key, elem.Index...)
		key = append(key, []byte("/")...)

		index := string(key)
		if _, ok := nftStakeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for nftStake")
		}
		nftStakeIndexMap[index] = struct{}{}
	}

	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateChainValidatorBridge()
	if err != nil {
		return err
	}

	err = gs.ValidateEpoch()
	if err != nil {
		return err
	}

	err = gs.ValidateLastEpochBlock()
	if err != nil {
		return err
	}

	err = gs.ValidateNftStake()
	if err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
