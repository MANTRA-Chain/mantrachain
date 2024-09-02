package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestBridgedQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	wctx := ctx
	msgs := createNBridged(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetBridgedRequest
		response *types.QueryGetBridgedResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetBridgedRequest{
				EthTxHash: msgs[0].EthTxHash,
			},
			response: &types.QueryGetBridgedResponse{Bridged: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetBridgedRequest{
				EthTxHash: msgs[1].EthTxHash,
			},
			response: &types.QueryGetBridgedResponse{Bridged: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetBridgedRequest{
				EthTxHash: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Bridged(wctx, tc.request)
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

func TestBridgedQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.BridgeKeeper(t)
	wctx := ctx
	msgs := createNBridged(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllBridgedRequest {
		return &types.QueryAllBridgedRequest{
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
			resp, err := keeper.BridgedAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bridged), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bridged),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.BridgedAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bridged), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bridged),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.BridgedAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Bridged),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.BridgedAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
