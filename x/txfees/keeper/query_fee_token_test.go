package keeper_test

import (
	"strconv"
	"testing"

	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestFeeTokenQuerySingle(t *testing.T) {
	keeper, ctx := TxfeesKeeper(t)
	msgs := createNFeeToken(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetFeeTokenRequest
		response *types.QueryGetFeeTokenResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetFeeTokenRequest{
				Denom: msgs[0].Denom,
			},
			response: &types.QueryGetFeeTokenResponse{FeeToken: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetFeeTokenRequest{
				Denom: msgs[1].Denom,
			},
			response: &types.QueryGetFeeTokenResponse{FeeToken: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetFeeTokenRequest{
				Denom: strconv.Itoa(100000),
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
			response, err := keeper.FeeToken(ctx, tc.request)
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

func TestFeeTokenQueryPaginated(t *testing.T) {
	keeper, ctx := TxfeesKeeper(t)
	msgs := createNFeeToken(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllFeeTokenRequest {
		return &types.QueryAllFeeTokenRequest{
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
			resp, err := keeper.FeeTokenAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.FeeToken), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.FeeToken),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.FeeTokenAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.FeeToken), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.FeeToken),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.FeeTokenAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.FeeToken),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.FeeTokenAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
