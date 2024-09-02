package keeper_test

import (
	"context"
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/txfees/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(tb testing.TB) (types.MsgServer, context.Context) {
	tb.Helper()
	k, ctx := TxfeesKeeper(tb)
	return keeper.NewMsgServerImpl(*k), ctx
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
