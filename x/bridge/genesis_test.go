package bridge_test

import (
	"testing"

	keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
	"github.com/LimeChain/mantrachain/testutil/nullify"
	"github.com/LimeChain/mantrachain/x/bridge"
	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		Cw20Contract: &types.Cw20Contract{
			CodeId: 81,
			Ver:    "36",
			Path:   "33",
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BridgeKeeper(t)
	bridge.InitGenesis(ctx, *k, genesisState)
	got := bridge.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.Cw20Contract, got.Cw20Contract)
	// this line is used by starport scaffolding # genesis/test/assert
}
