package tax_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/testutil/nullify"
	tax "github.com/MANTRA-Chain/mantrachain/x/tax/module"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx, _ := keepertest.TaxKeeper(t)
	err := tax.InitGenesis(ctx, k, genesisState)
	require.NoError(t, err)
	got, err := tax.ExportGenesis(ctx, k)
	require.NoError(t, err)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
