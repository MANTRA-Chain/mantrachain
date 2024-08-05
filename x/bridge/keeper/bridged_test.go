package keeper_test

import (
	"strconv"
	"testing"

	sdkmath "cosmossdk.io/math"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	"github.com/MANTRA-Finance/mantrachain/testutil/sample"
	"github.com/MANTRA-Finance/mantrachain/x/bridge/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNBridged(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Bridged {
	items := make([]types.Bridged, n)
	for i := range items {
		items[i].Index = types.BridgedKey(strconv.Itoa(i))
		items[i].EthTxHash = strconv.Itoa(i)
		coin := sdk.NewCoin("uom", sdkmath.NewInt(100))
		items[i].Amount = &coin
		items[i].Receiver = sample.AccAddress()

		keeper.SetBridged(ctx, items[i])
	}
	return items
}

func TestBridgedGet(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	items := createNBridged(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetBridged(ctx,
			item.EthTxHash,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestBridgedGetAll(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	items := createNBridged(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllBridged(ctx)),
	)
}
