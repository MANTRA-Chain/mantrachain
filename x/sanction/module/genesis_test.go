package sanction_test

import (
	"strings"
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
	}
	require.NoError(t, genesisState.Validate())

	k, ctx, _ := keepertest.SanctionKeeper(t)
	require.NoError(t, sanction.InitGenesis(ctx, k, genesisState))
	got, err := sanction.ExportGenesis(ctx, k)
	require.NoError(t, err)

	require.Equal(t, genesisState.Params, got.Params)
	require.ElementsMatch(t, genesisState.BlacklistAccounts, got.BlacklistAccounts)
	require.ElementsMatch(t, genesisState.Sanctioners, got.Sanctioners)
}

// TestInitGenesisNormalizesSanctioners checks a mixed-case sanctioner address
// in genesis is stored in canonical form, so it matches the tx signer.
func TestInitGenesisNormalizesSanctioners(t *testing.T) {
	full := sdk.AccAddress([]byte("full________________")).String()
	gs := *types.DefaultGenesis()
	gs.Sanctioners = []types.Sanctioner{{Address: strings.ToUpper(full), Role: types.SanctionerRole_SANCTIONER_ROLE_FULL}}
	require.NoError(t, gs.Validate())

	k, ctx, _ := keepertest.SanctionKeeper(t)
	require.NoError(t, sanction.InitGenesis(ctx, k, gs))
	has, err := k.Sanctioners.Has(ctx, full)
	require.NoError(t, err)
	require.True(t, has)
}

func TestGenesisEmptyRoundTrip(t *testing.T) {
	k, ctx, _ := keepertest.SanctionKeeper(t)
	require.NoError(t, sanction.InitGenesis(ctx, k, *types.DefaultGenesis()))
	got, err := sanction.ExportGenesis(ctx, k)
	require.NoError(t, err)
	require.Equal(t, types.DefaultParams(), got.Params)
	require.Empty(t, got.BlacklistAccounts)
	require.Empty(t, got.Sanctioners)
}
