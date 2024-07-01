package farming

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/farming/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/farming/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	ctx, writeCache := ctx.CacheContext()

	k.SetCurrentEpochDays(ctx, genState.CurrentEpochDays)
	if addr := k.GetAccountKeeper().GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	k.SetGlobalPlanId(ctx, genState.GlobalPlanId)

	for _, record := range genState.PlanRecords {
		plan, _ := types.UnpackPlan(&record.Plan) // Already validated
		k.SetPlan(ctx, plan)
	}

	totalStakings := map[string]math.Int{} // (staking coin denom) => (amount)

	for _, record := range genState.StakingRecords {
		farmerAcc, err := sdk.AccAddressFromBech32(record.Farmer)
		if err != nil {
			panic(err)
		}
		k.SetStaking(ctx, record.StakingCoinDenom, farmerAcc, record.Staking)

		amt, ok := totalStakings[record.StakingCoinDenom]
		if !ok {
			amt = math.ZeroInt()
		}
		amt = amt.Add(record.Staking.Amount)
		totalStakings[record.StakingCoinDenom] = amt
	}

	for _, record := range genState.TotalStakingsRecords {
		if !record.Amount.Equal(totalStakings[record.StakingCoinDenom]) {
			panic(fmt.Sprintf("TotalStaking for %s differs from the actual value; have %s, want %s",
				record.StakingCoinDenom, totalStakings[record.StakingCoinDenom], record.Amount))
		}
		stakingReserveCoins := k.GetBankKeeper().GetAllBalances(ctx, types.StakingReserveAcc(record.StakingCoinDenom))
		if !record.StakingReserveCoins.Equal(stakingReserveCoins) {
			panic(fmt.Sprintf("StakingReserveCoins differs from the actual value; have %s, want %s",
				stakingReserveCoins, record.StakingReserveCoins))
		}
	}

	if len(totalStakings) != len(genState.TotalStakingsRecords) {
		panic(fmt.Sprintf("the number of TotalStaking differs from the actual value; have %d, want %d",
			len(totalStakings), len(genState.TotalStakingsRecords)))
	}

	for _, record := range genState.TotalStakingsRecords {
		k.SetTotalStakings(ctx, record.StakingCoinDenom, types.TotalStakings{Amount: record.Amount})
	}

	for _, record := range genState.QueuedStakingRecords {
		farmerAcc, err := sdk.AccAddressFromBech32(record.Farmer)
		if err != nil {
			panic(err)
		}
		k.SetQueuedStaking(ctx, record.EndTime, record.StakingCoinDenom, farmerAcc, record.QueuedStaking)
	}

	for _, record := range genState.HistoricalRewardsRecords {
		k.SetHistoricalRewards(ctx, record.StakingCoinDenom, record.Epoch, record.HistoricalRewards)
	}

	for _, record := range genState.OutstandingRewardsRecords {
		k.SetOutstandingRewards(ctx, record.StakingCoinDenom, record.OutstandingRewards)
	}

	for _, record := range genState.UnharvestedRewardsRecords {
		farmerAcc, err := sdk.AccAddressFromBech32(record.Farmer)
		if err != nil {
			panic(err)
		}
		k.SetUnharvestedRewards(ctx, farmerAcc, record.StakingCoinDenom, record.UnharvestedRewards)
	}

	for _, record := range genState.CurrentEpochRecords {
		k.SetCurrentEpoch(ctx, record.StakingCoinDenom, record.CurrentEpoch)
	}

	if genState.LastEpochTime != nil {
		k.SetLastEpochTime(ctx, *genState.LastEpochTime)
	}

	err := k.ValidateRemainingRewardsAmount(ctx)
	if err != nil {
		panic(err)
	}
	rewardsPoolCoins := k.GetBankKeeper().GetAllBalances(ctx, types.RewardsReserveAcc)
	if !genState.RewardPoolCoins.Equal(rewardsPoolCoins) {
		panic(fmt.Sprintf("RewardPoolCoins differs from the actual value; have %s, want %s",
			rewardsPoolCoins, genState.RewardPoolCoins))
	}

	err = k.ValidateStakingReservedAmount(ctx)
	if err != nil {
		panic(err)
	}

	if err := k.ValidateOutstandingRewardsAmount(ctx); err != nil {
		panic(err)
	}

	if err := k.ValidateUnharvestedRewardsAmount(ctx); err != nil {
		panic(err)
	}

	writeCache()
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PlanRecords = []types.PlanRecord{}
	for _, plan := range k.GetPlans(ctx) {
		any, err := types.PackPlan(plan)
		if err != nil {
			panic(err)
		}
		genesis.PlanRecords = append(genesis.PlanRecords, types.PlanRecord{
			Plan:             *any,
			FarmingPoolCoins: k.GetBankKeeper().GetAllBalances(ctx, plan.GetFarmingPoolAddress()),
		})
	}

	genesis.StakingRecords = []types.StakingRecord{}
	k.IterateStakings(ctx, func(stakingCoinDenom string, farmerAcc sdk.AccAddress, staking types.Staking) (stop bool) {
		genesis.StakingRecords = append(genesis.StakingRecords, types.StakingRecord{
			StakingCoinDenom: stakingCoinDenom,
			Farmer:           farmerAcc.String(),
			Staking:          staking,
		})
		return false
	})

	genesis.QueuedStakingRecords = []types.QueuedStakingRecord{}
	k.IterateQueuedStakings(ctx, func(endTime time.Time, stakingCoinDenom string, farmerAcc sdk.AccAddress, queuedStaking types.QueuedStaking) (stop bool) {
		genesis.QueuedStakingRecords = append(genesis.QueuedStakingRecords, types.QueuedStakingRecord{
			EndTime:          endTime,
			StakingCoinDenom: stakingCoinDenom,
			Farmer:           farmerAcc.String(),
			QueuedStaking:    queuedStaking,
		})
		return false
	})

	genesis.TotalStakingsRecords = []types.TotalStakingsRecord{}
	k.IterateTotalStakings(ctx, func(stakingCoinDenom string, ts types.TotalStakings) (stop bool) {
		genesis.TotalStakingsRecords = append(genesis.TotalStakingsRecords, types.TotalStakingsRecord{
			StakingCoinDenom:    stakingCoinDenom,
			Amount:              ts.Amount,
			StakingReserveCoins: k.GetBankKeeper().GetAllBalances(ctx, types.StakingReserveAcc(stakingCoinDenom)),
		})
		return false
	})

	genesis.HistoricalRewardsRecords = []types.HistoricalRewardsRecord{}
	k.IterateHistoricalRewards(ctx, func(stakingCoinDenom string, epoch uint64, rewards types.HistoricalRewards) (stop bool) {
		genesis.HistoricalRewardsRecords = append(genesis.HistoricalRewardsRecords, types.HistoricalRewardsRecord{
			StakingCoinDenom:  stakingCoinDenom,
			Epoch:             epoch,
			HistoricalRewards: rewards,
		})
		return false
	})

	genesis.OutstandingRewardsRecords = []types.OutstandingRewardsRecord{}
	k.IterateOutstandingRewards(ctx, func(stakingCoinDenom string, rewards types.OutstandingRewards) (stop bool) {
		genesis.OutstandingRewardsRecords = append(genesis.OutstandingRewardsRecords, types.OutstandingRewardsRecord{
			StakingCoinDenom:   stakingCoinDenom,
			OutstandingRewards: rewards,
		})
		return false
	})

	genesis.UnharvestedRewardsRecords = []types.UnharvestedRewardsRecord{}
	k.IterateAllUnharvestedRewards(ctx, func(farmerAcc sdk.AccAddress, stakingCoinDenom string, rewards types.UnharvestedRewards) (stop bool) {
		genesis.UnharvestedRewardsRecords = append(genesis.UnharvestedRewardsRecords, types.UnharvestedRewardsRecord{
			Farmer:             farmerAcc.String(),
			StakingCoinDenom:   stakingCoinDenom,
			UnharvestedRewards: rewards,
		})
		return false
	})

	genesis.CurrentEpochRecords = []types.CurrentEpochRecord{}
	k.IterateCurrentEpochs(ctx, func(stakingCoinDenom string, currentEpoch uint64) (stop bool) {
		genesis.CurrentEpochRecords = append(genesis.CurrentEpochRecords, types.CurrentEpochRecord{
			StakingCoinDenom: stakingCoinDenom,
			CurrentEpoch:     currentEpoch,
		})
		return false
	})

	var epochTime *time.Time
	tempEpochTime, found := k.GetLastEpochTime(ctx)
	if found {
		epochTime = &tempEpochTime
	}

	genesis.LastEpochTime = epochTime
	genesis.GlobalPlanId = k.GetGlobalPlanId(ctx)
	genesis.RewardPoolCoins = k.GetBankKeeper().GetAllBalances(ctx, types.RewardsReserveAcc)
	genesis.CurrentEpochDays = k.GetCurrentEpochDays(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
