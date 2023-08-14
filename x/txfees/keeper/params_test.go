package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "mantrachain/testutil/keeper"
	"mantrachain/x/txfees/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TxfeesKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
