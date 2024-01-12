package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/AumegaChain/aumega/testutil/keeper"
	"github.com/AumegaChain/aumega/testutil/nullify"
	"github.com/AumegaChain/aumega/x/rewards/types"
)

func TestSnapshotQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNSnapshot(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetSnapshotRequest
		response *types.QueryGetSnapshotResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetSnapshotRequest{Id: msgs[0].Id},
			response: &types.QueryGetSnapshotResponse{Snapshot: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetSnapshotRequest{Id: msgs[1].Id},
			response: &types.QueryGetSnapshotResponse{Snapshot: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetSnapshotRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Snapshot(wctx, tc.request)
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

func TestSnapshotQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.RewardsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNSnapshot(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllSnapshotRequest {
		return &types.QueryAllSnapshotRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SnapshotAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Snapshot), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Snapshot),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SnapshotAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Snapshot), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Snapshot),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.SnapshotAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Snapshot),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.SnapshotAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
