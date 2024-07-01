package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:                    DefaultParams(),
		GlobalPlanId:              0,
		PlanRecords:               []PlanRecord{},
		StakingRecords:            []StakingRecord{},
		QueuedStakingRecords:      []QueuedStakingRecord{},
		TotalStakingsRecords:      []TotalStakingsRecord{},
		HistoricalRewardsRecords:  []HistoricalRewardsRecord{},
		OutstandingRewardsRecords: []OutstandingRewardsRecord{},
		UnharvestedRewardsRecords: []UnharvestedRewardsRecord{},
		CurrentEpochRecords:       []CurrentEpochRecord{},
		RewardPoolCoins:           sdk.Coins{},
		LastEpochTime:             nil,
		CurrentEpochDays:          DefaultCurrentEpochDays,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
