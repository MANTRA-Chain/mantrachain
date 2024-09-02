package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNCampaign(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Campaign {
	items := make([]types.Campaign, n)
	for i := range items {
		items[i].Index = []byte(strconv.Itoa(i))
		items[i].Id = uint64(i)

		keeper.SetCampaign(ctx, items[i])
	}
	return items
}

func TestCampaignGet(t *testing.T) {
	keeper, ctx := keepertest.AirdropKeeper(t)
	items := createNCampaign(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetCampaign(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestCampaignRemove(t *testing.T) {
	keeper, ctx := keepertest.AirdropKeeper(t)
	items := createNCampaign(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveCampaign(ctx,
			item.Index,
		)
		_, found := keeper.GetCampaign(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestCampaignGetAll(t *testing.T) {
	keeper, ctx := keepertest.AirdropKeeper(t)
	items := createNCampaign(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllCampaign(ctx)),
	)
}
