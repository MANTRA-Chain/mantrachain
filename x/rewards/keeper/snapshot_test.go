package keeper_test

import (
	"testing"

	keepertest "github.com/MANTRA-Finance/aumega/testutil/keeper"
	"github.com/MANTRA-Finance/aumega/testutil/nullify"
	"github.com/MANTRA-Finance/aumega/x/rewards/keeper"
	"github.com/MANTRA-Finance/aumega/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNSnapshot(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Snapshot {
	items := make([]types.Snapshot, n)
	for i := range items {
		items[i].Id = keeper.AppendSnapshot(ctx, items[i])
	}
	return items
}

func TestSnapshotGet(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	items := createNSnapshot(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetSnapshot(ctx, item.PairId, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestSnapshotRemove(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	items := createNSnapshot(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveSnapshot(ctx, item.PairId, item.Id)
		_, found := keeper.GetSnapshot(ctx, item.PairId, item.Id)
		require.False(t, found)
	}
}

func TestSnapshotGetAll(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	items := createNSnapshot(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllSnapshot(ctx, 0)),
	)
}
