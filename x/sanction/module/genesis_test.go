package sanction_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/v8/testutil/keeper"
	sanction "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/module"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	full := sdk.AccAddress([]byte("full________________")).String()
	addOnly := sdk.AccAddress([]byte("addonly_____________")).String()
	legacy := sdk.AccAddress([]byte("legacy______________")).String()
	flagged := sdk.AccAddress([]byte("flagged_____________")).String()

	genesisState := types.GenesisState{
		Params:            types.NewParams(10, 64),
		BlacklistAccounts: []string{legacy, flagged},
		Sanctioners: []types.Sanctioner{
			{Address: full, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL},
			{Address: addOnly, Role: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
		},
		SanctionRecords: []types.SanctionedAddress{{
			Address: flagged,
			Record: &types.SanctionRecord{
				Actor:     addOnly,
				ActorRole: types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY,
				Reason:    "rule 7",
				Height:    42,
			},
		}},
	}
	require.NoError(t, genesisState.Validate())

	k, ctx, _ := keepertest.SanctionKeeper(t)
	require.NoError(t, sanction.InitGenesis(ctx, k, genesisState))
	got, err := sanction.ExportGenesis(ctx, k)
	require.NoError(t, err)

	require.Equal(t, genesisState.Params, got.Params)
	require.ElementsMatch(t, genesisState.BlacklistAccounts, got.BlacklistAccounts)
	require.ElementsMatch(t, genesisState.Sanctioners, got.Sanctioners)
	require.Equal(t, genesisState.SanctionRecords, got.SanctionRecords)
}

func TestGenesisEmptyRoundTrip(t *testing.T) {
	k, ctx, _ := keepertest.SanctionKeeper(t)
	require.NoError(t, sanction.InitGenesis(ctx, k, *types.DefaultGenesis()))
	got, err := sanction.ExportGenesis(ctx, k)
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams(), got.Params)
	require.Empty(t, got.BlacklistAccounts)
	require.Empty(t, got.Sanctioners)
	require.Empty(t, got.SanctionRecords)
}
