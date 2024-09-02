package keeper_test

import (
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := TxfeesKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)

	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
