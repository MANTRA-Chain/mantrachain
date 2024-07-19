package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	var plans []PlanI
	for _, record := range gs.PlanRecords {
		if err := record.Validate(); err != nil {
			return err
		}
		plan, _ := UnpackPlan(&record.Plan)
		if plan.GetId() > gs.GlobalPlanId {
			return fmt.Errorf("plan id is greater than the global last plan id")
		}
		plans = append(plans, plan)
	}

	if err := ValidateTotalEpochRatio(plans); err != nil {
		return err
	}

	for _, record := range gs.StakingRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.QueuedStakingRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.TotalStakingsRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.HistoricalRewardsRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.OutstandingRewardsRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.UnharvestedRewardsRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	for _, record := range gs.CurrentEpochRecords {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	if err := gs.RewardPoolCoins.Validate(); err != nil {
		return err
	}

	if gs.CurrentEpochDays == 0 {
		return fmt.Errorf("current epoch days must be positive")
	}

	return nil
}

// Validate validates PlanRecord.
func (record PlanRecord) Validate() error {
	plan, err := UnpackPlan(&record.Plan)
	if err != nil {
		return err
	}
	if err := plan.Validate(); err != nil {
		return err
	}
	if err := record.FarmingPoolCoins.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates StakingRecord.
func (record StakingRecord) Validate() error {
	if _, err := sdk.AccAddressFromBech32(record.Farmer); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(record.StakingCoinDenom); err != nil {
		return err
	}
	if !record.Staking.Amount.IsPositive() {
		return fmt.Errorf("staking amount must be positive: %s", record.Staking.Amount)
	}
	return nil
}

// Validate validates StakingRecord.
func (record TotalStakingsRecord) Validate() error {
	if err := sdk.ValidateDenom(record.StakingCoinDenom); err != nil {
		return err
	}
	if !record.Amount.IsPositive() {
		return fmt.Errorf("total staking amount must be positive: %s", record.Amount)
	}
	if err := record.StakingReserveCoins.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates QueuedStakingRecord.
func (record QueuedStakingRecord) Validate() error {
	if _, err := sdk.AccAddressFromBech32(record.Farmer); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(record.StakingCoinDenom); err != nil {
		return err
	}
	if !record.QueuedStaking.Amount.IsPositive() {
		return fmt.Errorf("queued staking amount must be positive: %s", record.QueuedStaking.Amount)
	}
	return nil
}

// Validate validates HistoricalRewardsRecord.
func (record HistoricalRewardsRecord) Validate() error {
	if err := sdk.ValidateDenom(record.StakingCoinDenom); err != nil {
		return err
	}
	if err := record.HistoricalRewards.CumulativeUnitRewards.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates OutstandingRewardsRecord.
func (record OutstandingRewardsRecord) Validate() error {
	if err := sdk.ValidateDenom(record.StakingCoinDenom); err != nil {
		return err
	}
	if err := record.OutstandingRewards.Rewards.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates UnharvestedRewardsRecord.
func (record UnharvestedRewardsRecord) Validate() error {
	if _, err := sdk.AccAddressFromBech32(record.Farmer); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(record.StakingCoinDenom); err != nil {
		return err
	}
	if err := record.UnharvestedRewards.Rewards.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates CurrentEpochRecord.
func (record CurrentEpochRecord) Validate() error {
	if err := sdk.ValidateDenom(record.StakingCoinDenom); err != nil {
		return err
	}
	return nil
}
