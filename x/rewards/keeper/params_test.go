package keeper_test

import (
	"testing"

	testkeeper "github.com/MANTRA-Finance/aumega/testutil/keeper"
	"github.com/MANTRA-Finance/aumega/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.RewardsKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
