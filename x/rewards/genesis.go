package rewards

import (
	"github.com/MANTRA-Finance/mantrachain/x/rewards/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the snapshot
for _, elem := range genState.SnapshotList {
	k.SetSnapshot(ctx, elem)
}

// Set snapshot count
k.SetSnapshotCount(ctx, genState.SnapshotCount)
// Set all the provider
for _, elem := range genState.ProviderList {
	k.SetProvider(ctx, elem)
}
// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.SnapshotList = k.GetAllSnapshot(ctx)
genesis.SnapshotCount = k.GetSnapshotCount(ctx)
genesis.ProviderList = k.GetAllProvider(ctx)
// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
