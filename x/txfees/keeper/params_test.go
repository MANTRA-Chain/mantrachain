package keeper_test

import (
	"testing"

	testkeeper "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TxfeesKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
