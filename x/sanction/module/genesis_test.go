package sanction_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/v8/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/testutil/nullify"
	sanction "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/module"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx, _ := keepertest.SanctionKeeper(t)
	err := sanction.InitGenesis(ctx, k, genesisState)
	require.NoError(t, err)
	got, err := sanction.ExportGenesis(ctx, k)
	require.NoError(t, err)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
