package mns_test

import (
	"testing"

	keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
	"github.com/LimeChain/mantrachain/testutil/nullify"
	"github.com/LimeChain/mantrachain/x/mns"
	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MnsKeeper(t)
	mns.InitGenesis(ctx, *k, genesisState)
	got := mns.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	

	// this line is used by starport scaffolding # genesis/test/assert
}
