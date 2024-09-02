package lpfarm

import (
	"time"

	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/lpfarm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	if genState.LastBlockTime != nil {
		k.SetLastBlockTime(ctx, *genState.LastBlockTime)
	}
	if genState.LastPlanId > 0 {
		k.SetLastPlanId(ctx, genState.LastPlanId)
	}
	k.SetNumPrivatePlans(ctx, genState.NumPrivatePlans)
	for _, plan := range genState.Plans {
		k.SetPlan(ctx, plan)
	}
	for _, farm := range genState.Farms {
		k.SetFarm(ctx, farm.Denom, farm.Farm)
	}
	for _, position := range genState.Positions {
		k.SetPosition(ctx, position)
	}
	for _, hist := range genState.HistoricalRewards {
		k.SetHistoricalRewards(ctx, hist.Denom, hist.Period, hist.HistoricalRewards)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	var lastBlockTimePtr *time.Time
	genesis.LastBlockTime = lastBlockTimePtr
	lastBlockTime, found := k.GetLastBlockTime(ctx)
	if found {
		genesis.LastBlockTime = &lastBlockTime
	}

	lastPlanId, _ := k.GetLastPlanId(ctx)
	genesis.LastPlanId = lastPlanId

	plans := []types.Plan{}
	k.IterateAllPlans(ctx, func(plan types.Plan) (stop bool) {
		plans = append(plans, plan)
		return false
	})
	genesis.NumPrivatePlans = k.GetNumPrivatePlans(ctx)
	genesis.Plans = plans

	farms := []types.FarmRecord{}
	k.IterateAllFarms(ctx, func(denom string, farm types.Farm) (stop bool) {
		farms = append(farms, types.FarmRecord{
			Denom: denom,
			Farm:  farm,
		})
		return false
	})
	genesis.Farms = farms

	positions := []types.Position{}
	k.IterateAllPositions(ctx, func(position types.Position) (stop bool) {
		positions = append(positions, position)
		return false
	})
	genesis.Positions = positions

	hists := []types.HistoricalRewardsRecord{}
	k.IterateAllHistoricalRewards(
		ctx, func(denom string, period uint64, hist types.HistoricalRewards) (stop bool) {
			hists = append(hists, types.HistoricalRewardsRecord{
				Denom:             denom,
				Period:            period,
				HistoricalRewards: hist,
			})
			return false
		})
	genesis.HistoricalRewards = hists

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
