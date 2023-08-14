package txfees_test

import (
	"testing"

	keepertest "mantrachain/testutil/keeper"
	"mantrachain/testutil/nullify"
	"mantrachain/x/txfees"
	"mantrachain/x/txfees/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TxfeesKeeper(t)
	txfees.InitGenesis(ctx, *k, genesisState)
	got := txfees.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	

	// this line is used by starport scaffolding # genesis/test/assert
}
