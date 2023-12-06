package rewards_test

import (
	"testing"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	"github.com/MANTRA-Finance/mantrachain/x/rewards"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		SnapshotList: []types.Snapshot{
		{
			Id: 0,
		},
		{
			Id: 1,
		},
	},
	SnapshotCount: 2,
	ProviderList: []types.Provider{
		{
			Index: "0",
},
		{
			Index: "1",
},
	},
	// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.RewardsKeeper(t)
	rewards.InitGenesis(ctx, *k, genesisState)
	got := rewards.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.SnapshotList, got.SnapshotList)
require.Equal(t, genesisState.SnapshotCount, got.SnapshotCount)
require.ElementsMatch(t, genesisState.ProviderList, got.ProviderList)
// this line is used by starport scaffolding # genesis/test/assert
}
