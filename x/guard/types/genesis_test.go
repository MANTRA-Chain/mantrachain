package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisStateValidate(t *testing.T) {
	testCases := []struct {
		name         string
		genesisState GenesisState
		expErr       bool
	}{
		{
			"valid genesisState",
			GenesisState{
				Params: Params{
					AdminAccount:                            "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
					AccountPrivilegesTokenCollectionCreator: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
					AccountPrivilegesTokenCollectionId:      "id",
					DefaultPrivileges:                       make([]byte, 32),
					BaseDenom:                               "uaum",
				},
				AccountPrivilegesList: []*AccountPrivileges{
					{
						Account:    sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t").Bytes(),
						Privileges: []byte{0x01},
					},
				},
				GuardTransferCoins: []byte{0x01},
				RequiredPrivilegesList: []*RequiredPrivileges{
					{
						Index:      []byte{0x01},
						Privileges: []byte{0x01},
						Kind:       "coin",
					},
				},
			},
			false,
		},
		{"empty genesisState", GenesisState{}, true},
		{
			"invalid params ",
			GenesisState{
				Params: Params{
					AdminAccount:                            "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
					AccountPrivilegesTokenCollectionCreator: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
					AccountPrivilegesTokenCollectionId:      "id",
					DefaultPrivileges:                       make([]byte, 31),
					BaseDenom:                               "uaum",
				},
			},
			true,
		},
		{
			"dup account privileges",
			GenesisState{
				Params: Params{
					AdminAccount:                            "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
					AccountPrivilegesTokenCollectionCreator: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
					AccountPrivilegesTokenCollectionId:      "id",
					DefaultPrivileges:                       make([]byte, 32),
					BaseDenom:                               "uaum",
				},
				AccountPrivilegesList: []*AccountPrivileges{
					{
						Account:    sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t").Bytes(),
						Privileges: []byte{0x01},
					},
					{
						Account:    sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t").Bytes(),
						Privileges: []byte{0x01},
					},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesisState.Validate()

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
