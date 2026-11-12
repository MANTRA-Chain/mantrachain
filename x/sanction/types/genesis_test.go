package types_test

import (
	"testing"

	appparams "github.com/MANTRA-Chain/mantrachain/v8/app/params"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	appparams.SetAddressPrefixes()
	const (
		addr1 = "mantra1hz88lcv4xmfzsrsvmtynhc2medke0m20zq0tex"
		addr2 = "mantra1nmefvq9aa7t5p85vp8weuuwrmrppjgm6r9ccpa"
	)
	full := types.Sanctioner{Address: addr1, Role: types.SanctionerRole_SANCTIONER_ROLE_FULL}

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
				Params:            types.DefaultParams(),
				BlacklistAccounts: []string{addr1, addr2},
			},
			valid: true,
		},
		{
			desc: "invalid genesis state with non-bech32 account",
			genState: &types.GenesisState{
				Params:            types.DefaultParams(),
				BlacklistAccounts: []string{"not-a-valid-bech32"},
			},
			valid: false,
		},
		{
			desc:     "zero params rejected",
			genState: &types.GenesisState{},
			valid:    false,
		},
		{
			desc: "sanctioners are valid",
			genState: &types.GenesisState{
				Params:            types.DefaultParams(),
				BlacklistAccounts: []string{addr2},
				Sanctioners:       []types.Sanctioner{full},
			},
			valid: true,
		},
		{
			desc: "duplicate sanctioner rejected",
			genState: &types.GenesisState{
				Params:      types.DefaultParams(),
				Sanctioners: []types.Sanctioner{full, full},
			},
			valid: false,
		},
		{
			desc: "unspecified sanctioner role rejected",
			genState: &types.GenesisState{
				Params:      types.DefaultParams(),
				Sanctioners: []types.Sanctioner{{Address: addr1}},
			},
			valid: false,
		},
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
