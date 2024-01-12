package keeper_test

import (
	"testing"

	testkeeper "github.com/AumegaChain/aumega/testutil/keeper"
	"github.com/AumegaChain/aumega/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.RewardsKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
