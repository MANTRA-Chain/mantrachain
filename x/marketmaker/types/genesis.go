package types

import (
	"fmt"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:         DefaultParams(),
		MarketMakers:   []MarketMaker{},
		Incentives:     []Incentive{},
		DepositRecords: []DepositRecord{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	for _, record := range gs.MarketMakers {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.Incentives {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.DepositRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	if err := ValidateDepositRecords(gs.MarketMakers, gs.DepositRecords); err != nil {
		return err
	}

	return gs.Params.Validate()
}

func ValidateDepositRecords(mms []MarketMaker, DepositRecords []DepositRecord) error {
	// not eligible market maker must have deposit record
	for _, mm := range mms {
		if !mm.Eligible {
			found := false
			for _, record := range DepositRecords {
				if record.PairId == mm.PairId && record.Address == mm.Address {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("deposit invariant failed, not eligible market maker must have deposit record")
			}
		}
	}

	// deposit record's market maker must not be eligible
	for _, record := range DepositRecords {
		found := false
		for _, mm := range mms {
			if !mm.Eligible && record.PairId == mm.PairId && record.Address == mm.Address {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("deposit invariant failed, deposit record's market maker must not be eligible")
		}
	}
	return nil
}
