package keeper_test

import (
	"strconv"
	"testing"

	"github.com/AumegaChain/aumega/testutil/nullify"
	"github.com/AumegaChain/aumega/x/txfees/keeper"
	"github.com/AumegaChain/aumega/x/txfees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNFeeToken(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.FeeToken {
	items := make([]types.FeeToken, n)
	for i := range items {
		items[i].Denom = strconv.Itoa(i)

		keeper.SetFeeToken(ctx, items[i])
	}
	return items
}

func TestFeeTokenGet(t *testing.T) {
	keeper, ctx := TxfeesKeeper(t)
	items := createNFeeToken(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetFeeToken(ctx,
			item.Denom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestFeeTokenRemove(t *testing.T) {
	keeper, ctx := TxfeesKeeper(t)
	items := createNFeeToken(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveFeeToken(ctx,
			item.Denom,
		)
		_, found := keeper.GetFeeToken(ctx,
			item.Denom,
		)
		require.False(t, found)
	}
}

func TestFeeTokenGetAll(t *testing.T) {
	keeper, ctx := TxfeesKeeper(t)
	items := createNFeeToken(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllFeeToken(ctx)),
	)
}
