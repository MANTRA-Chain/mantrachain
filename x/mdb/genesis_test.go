package mdb_test

import (
	"testing"

	keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
	"github.com/LimeChain/mantrachain/testutil/nullify"
	"github.com/LimeChain/mantrachain/x/mdb"
	"github.com/LimeChain/mantrachain/x/mdb/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MdbKeeper(t)
	mdb.InitGenesis(ctx, *k, genesisState)
	got := mdb.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	

	// this line is used by starport scaffolding # genesis/test/assert
}
