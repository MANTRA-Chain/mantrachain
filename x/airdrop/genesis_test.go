package airdrop_test

import (
	"testing"

	keepertest "github.com/MANTRA-Finance/mantrachain/testutil/keeper"
	"github.com/MANTRA-Finance/mantrachain/testutil/nullify"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		CampaignList: []types.Campaign{
			{
				Index: []byte("0"),
			},
			{
				Index: []byte("1"),
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AirdropKeeper(t)
	airdrop.InitGenesis(ctx, *k, genesisState)
	got := airdrop.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.CampaignList, got.CampaignList)
	// this line is used by starport scaffolding # genesis/test/assert
}
