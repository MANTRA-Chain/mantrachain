package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/LimeChain/mantrachain/testutil/keeper"
	"github.com/LimeChain/mantrachain/x/mdb/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MdbKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
