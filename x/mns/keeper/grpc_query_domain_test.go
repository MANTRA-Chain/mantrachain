package keeper_test

import (
    "strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/testutil/nullify"
	keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestDomainQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.MnsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDomain(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDomainRequest
		response *types.QueryGetDomainResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetDomainRequest{
			    Index: msgs[0].Index,
                
			},
			response: &types.QueryGetDomainResponse{Domain: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetDomainRequest{
			    Index: msgs[1].Index,
                
			},
			response: &types.QueryGetDomainResponse{Domain: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetDomainRequest{
			    Index:strconv.Itoa(100000),
                
			},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Domain(wctx, tc.request)
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

func TestDomainQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.MnsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDomain(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDomainRequest {
		return &types.QueryAllDomainRequest{
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
			resp, err := keeper.DomainAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Domain), step)
			require.Subset(t,
            	nullify.Fill(msgs),
            	nullify.Fill(resp.Domain),
            )
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DomainAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Domain), step)
			require.Subset(t,
            	nullify.Fill(msgs),
            	nullify.Fill(resp.Domain),
            )
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.DomainAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Domain),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.DomainAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
