package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/MANTRA-Finance/aumega/testutil/keeper"
	"github.com/MANTRA-Finance/aumega/x/rewards/keeper"
	"github.com/MANTRA-Finance/aumega/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.RewardsKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
