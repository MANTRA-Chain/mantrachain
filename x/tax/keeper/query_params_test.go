package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/tax/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
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
