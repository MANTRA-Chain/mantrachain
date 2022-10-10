package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/LimeChain/mantrachain/x/bridge/keeper"
	"github.com/LimeChain/mantrachain/x/bridge/types"
	keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
    "github.com/LimeChain/mantrachain/testutil/nullify"
)

func createTestCw20Contract(keeper *keeper.Keeper, ctx sdk.Context) types.Cw20Contract {
	item := types.Cw20Contract{}
	keeper.SetCw20Contract(ctx, item)
	return item
}

func TestCw20ContractGet(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	item := createTestCw20Contract(keeper, ctx)
	rst, found := keeper.GetCw20Contract(ctx)
    require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestCw20ContractRemove(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	createTestCw20Contract(keeper, ctx)
	keeper.RemoveCw20Contract(ctx)
    _, found := keeper.GetCw20Contract(ctx)
    require.False(t, found)
}
