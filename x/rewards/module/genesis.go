package rewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Set all the snapshot
	for _, elem := range genState.SnapshotList {
		k.SetSnapshot(ctx, elem)
	}

	// Set all the snapshots counts
	for _, elem := range genState.SnapshotCountList {
		k.SetSnapshotCount(ctx, elem)
	}

	// Set all the snapshots counts
	for _, elem := range genState.SnapshotStartIdList {
		k.SetSnapshotStartId(ctx, elem)
	}

	// Set all the provider
	for _, elem := range genState.ProviderList {
		k.SetProvider(ctx, elem)
	}
	k.SetDistributionPairsIdsBytes(ctx, genState.DistributionPairsIds)
	k.SetPurgePairsIdsBytes(ctx, genState.PurgePairsIds)
	k.SetSnapshotsLastDistributedAt(ctx, genState.SnapshotsLastDistributedAt)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.SnapshotList = k.GetAllSnapshot(ctx, 0)
	genesis.SnapshotCountList = k.GetAllSnapshotCount(ctx)
	genesis.SnapshotStartIdList = k.GetAllSnapshotStartId(ctx)
	genesis.ProviderList = k.GetAllProvider(ctx)
	genesis.DistributionPairsIds = k.GetDistributionPairsIdsBytes(ctx)
	genesis.PurgePairsIds = k.GetPurgePairsIdsBytes(ctx)
	genesis.SnapshotsLastDistributedAt = k.GetSnapshotsLastDistributedAt(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
