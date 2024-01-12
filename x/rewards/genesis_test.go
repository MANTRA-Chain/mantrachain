package rewards_test

import (
	"testing"

	keepertest "github.com/AumegaChain/aumega/testutil/keeper"
	"github.com/AumegaChain/aumega/testutil/nullify"
	"github.com/AumegaChain/aumega/x/rewards"
	"github.com/AumegaChain/aumega/x/rewards/types"
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
		SnapshotCountList: []types.SnapshotCount{
			{
				PairId: 0,
			},
			{
				PairId: 1,
			},
		},
		SnapshotStartIdList: []types.SnapshotStartId{
			{
				PairId: 0,
			},
			{
				PairId: 1,
			},
		},
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
	require.ElementsMatch(t, genesisState.ProviderList, got.ProviderList)
	// this line is used by starport scaffolding # genesis/test/assert
}
