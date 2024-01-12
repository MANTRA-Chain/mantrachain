package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/AumegaChain/aumega/testutil/keeper"
	"github.com/AumegaChain/aumega/testutil/nullify"
	"github.com/AumegaChain/aumega/x/rewards/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestRewardsQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetRewardsRequest
		response *types.QueryGetRewardsResponse
		err      error
	}{
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Rewards(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
