package types_test

import (
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/rewards/types"
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
				Params: types.DefaultParams(),
				SnapshotList: []types.Snapshot{
					{
						Id:     0,
						PairId: 1,
					},
					{
						Id:     1,
						PairId: 2,
					},
				},
				SnapshotCountList: []types.SnapshotCount{
					{
						PairId: 1,
						Count:  0,
					},
				},
				ProviderList: []types.Provider{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated snapshot",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				SnapshotList: []types.Snapshot{
					{
						Id:     0,
						PairId: 1,
					},
					{
						Id:     0,
						PairId: 1,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid snapshot count",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				SnapshotList: []types.Snapshot{
					{
						Id:     1,
						PairId: 1,
					},
				},
				SnapshotCountList: []types.SnapshotCount{
					{
						PairId: 1,
						Count:  0,
					},
					{
						PairId: 1,
						Count:  0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated provider",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ProviderList: []types.Provider{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
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
