package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/LimeChain/mantrachain/x/mns/types"
    "github.com/LimeChain/mantrachain/x/mns/keeper"
    keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MnsKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
