package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/LimeChain/mantrachain/x/bridge/types"
    "github.com/LimeChain/mantrachain/x/bridge/keeper"
    keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.BridgeKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
