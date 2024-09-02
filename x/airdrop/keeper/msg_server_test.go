package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(tb testing.TB) (types.MsgServer, context.Context) {
	tb.Helper()
	k, ctx := keepertest.AirdropKeeper(tb)
	return keeper.NewMsgServerImpl(*k), ctx
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
