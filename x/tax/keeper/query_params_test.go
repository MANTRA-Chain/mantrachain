package keeper_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/v2/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/v2/x/tax/keeper"
	"github.com/MANTRA-Chain/mantrachain/v2/x/tax/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	k, ctx, _ := keepertest.TaxKeeper(t)

	qs := keeper.NewQueryServerImpl(k)
	params := types.DefaultParams()
	require.NoError(t, k.Params.Set(ctx, params))

	response, err := qs.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
