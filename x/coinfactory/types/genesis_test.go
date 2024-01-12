package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AumegaChain/aumega/x/coinfactory/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
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
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "different admin from creator",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "cosmos15ejrsrfts5jfd8vekdje4t3t56nvflry92uegz",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "empty admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "no admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
					},
				},
			},
			valid: true,
		},
		{
			desc: "invalid admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "moose",
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "multiple denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/litecoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "duplicate denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/zcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: false,
		},
	} {
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
