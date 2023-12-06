package keeper_test

import (
	"strconv"
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNProvider(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Provider {
	items := make([]types.Provider, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
        
		keeper.SetProvider(ctx, items[i])
	}
	return items
}

func TestProviderGet(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	items := createNProvider(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetProvider(ctx,
		    item.Index,
            
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestProviderRemove(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	items := createNProvider(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveProvider(ctx,
		    item.Index,
            
		)
		_, found := keeper.GetProvider(ctx,
		    item.Index,
            
		)
		require.False(t, found)
	}
}

func TestProviderGetAll(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	items := createNProvider(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllProvider(ctx)),
	)
}
