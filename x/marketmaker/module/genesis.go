package marketmaker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	ctx, writeCache := ctx.CacheContext()

	// init to prevent nil slice, []types.IncentivePairs(nil)
	if len(genState.Params.IncentivePairs) == 0 {
		genState.Params.IncentivePairs = []types.IncentivePair{}
	}

	// validations
	if err := k.ValidateDepositReservedAmount(ctx); err != nil {
		panic(err)
	}

	if err := k.ValidateIncentiveReservedAmount(ctx, genState.Incentives); err != nil {
		panic(err)
	}

	for _, record := range genState.MarketMakers {
		if err := record.Validate(); err != nil {
			panic(err)
		}
		k.SetMarketMaker(ctx, record)
	}

	for _, record := range genState.Incentives {
		if err := record.Validate(); err != nil {
			panic(err)
		}
		k.SetIncentive(ctx, record)
	}

	for _, record := range genState.DepositRecords {
		if err := record.Validate(); err != nil {
			panic(err)
		}
		k.SetDeposit(ctx, record.GetAccAddress(), record.PairId, record.Amount)
	}

	writeCache()
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// init to prevent empty slice
	if len(genesis.Params.IncentivePairs) == 0 {
		genesis.Params.IncentivePairs = []types.IncentivePair{}
	}

	genesis.MarketMakers = k.GetAllMarketMakers(ctx)
	genesis.Incentives = k.GetAllIncentives(ctx)
	genesis.DepositRecords = k.GetAllDepositRecords(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
