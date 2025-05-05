package types_test

import (
	"testing"

	"github.com/MANTRA-Chain/mantrachain/v5/x/sanction/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				// this line is used by starport scaffolding # types/genesis/validField
				BlacklistAccounts: []string{"mantra1hz88lcv4xmfzsrsvmtynhc2medke0m20zq0tex", "mantra1nmefvq9aa7t5p85vp8weuuwrmrppjgm6r9ccpa"},
			},
			valid: true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				// this line is used by starport scaffolding # types/genesis/validField
				BlacklistAccounts: []string{"not-a-valid-bech32"},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
