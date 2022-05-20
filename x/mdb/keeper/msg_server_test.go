package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/LimeChain/mantrachain/x/mdb/types"
    "github.com/LimeChain/mantrachain/x/mdb/keeper"
    keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MdbKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
