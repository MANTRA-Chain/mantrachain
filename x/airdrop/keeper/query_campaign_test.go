package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCampaignQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.AirdropKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCampaign(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetCampaignRequest
		response *types.QueryGetCampaignResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetCampaignRequest{
				Id: msgs[0].Id,
			},
			response: &types.QueryGetCampaignResponse{Campaign: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetCampaignRequest{
				Id: msgs[1].Id,
			},
			response: &types.QueryGetCampaignResponse{Campaign: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetCampaignRequest{
				Id: 100000,
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
			response, err := keeper.Campaign(wctx, tc.request)
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

func TestCampaignQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.AirdropKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCampaign(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCampaignRequest {
		return &types.QueryAllCampaignRequest{
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
			resp, err := keeper.CampaignAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Campaign), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Campaign),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CampaignAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Campaign), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Campaign),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.CampaignAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Campaign),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.CampaignAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
