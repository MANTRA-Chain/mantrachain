package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestClaimMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.RewardsKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"
	for i := 0; i < 5; i++ {
		expected := &types.MsgClaim{Creator: creator} // Index: strconv.Itoa(i),

		_, err := srv.Claim(wctx, expected)
		require.Error(t, err, "provider not found: provider not found")
	}
}
