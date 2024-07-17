package liquidity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/liquidity/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	k.SetLastPairId(ctx, genState.LastPairId)
	k.SetLastPoolId(ctx, genState.LastPoolId)
	for _, pair := range genState.Pairs {
		k.SetPair(ctx, pair)
		k.SetPairIndex(ctx, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
		k.SetPairLookupIndex(ctx, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
		k.SetPairLookupIndex(ctx, pair.QuoteCoinDenom, pair.BaseCoinDenom, pair.Id)
	}
	for _, pool := range genState.Pools {
		k.SetPool(ctx, pool)
		k.SetPoolByReserveIndex(ctx, pool)
		k.SetPoolsByPairIndex(ctx, pool)
	}
	for _, req := range genState.DepositRequests {
		k.SetDepositRequest(ctx, req)
		k.SetDepositRequestIndex(ctx, req)
	}
	for _, req := range genState.WithdrawRequests {
		k.SetWithdrawRequest(ctx, req)
		k.SetWithdrawRequestIndex(ctx, req)
	}
	for _, order := range genState.Orders {
		k.SetOrder(ctx, order)
		k.SetOrderIndex(ctx, order)
	}
	for _, index := range genState.MarketMakingOrderIndexes {
		k.SetMMOrderIndex(ctx, index)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.Params = k.GetParams(ctx)
	genesis.LastPairId = k.GetLastPairId(ctx)
	genesis.LastPoolId = k.GetLastPoolId(ctx)
	genesis.Pairs = k.GetAllPairs(ctx)
	genesis.Pools = k.GetAllPools(ctx)
	genesis.DepositRequests = k.GetAllDepositRequests(ctx)
	genesis.WithdrawRequests = k.GetAllWithdrawRequests(ctx)
	genesis.Orders = k.GetAllOrders(ctx)
	genesis.MarketMakingOrderIndexes = k.GetAllMMOrderIndexes(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
