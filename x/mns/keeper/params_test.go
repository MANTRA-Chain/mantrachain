package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/LimeChain/mantrachain/testutil/keeper"
	"github.com/LimeChain/mantrachain/x/mns/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MnsKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
