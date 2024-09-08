package xfeemarket_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/testutil/nullify"
	xfeemarket "github.com/MANTRA-Chain/mantrachain/x/xfeemarket/module"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.XfeemarketKeeper(t, nil)
	err := xfeemarket.InitGenesis(ctx, k, genesisState)
	require.NoError(t, err)
	got, err := xfeemarket.ExportGenesis(ctx, k)
	require.NoError(t, err)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
