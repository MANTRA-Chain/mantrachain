package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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
					AdminAccount:                            "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
					AccountPrivilegesTokenCollectionCreator: "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
					AccountPrivilegesTokenCollectionId:      "id",
					DefaultPrivileges:                       make([]byte, 32),
					BaseDenom:                               "uom",
				},
				AccountPrivilegesList: []AccountPrivileges{
					{
						Account:    sdk.MustAccAddressFromBech32("cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw").Bytes(),
						Privileges: []byte{0x01},
					},
				},
				GuardTransferCoins: []byte{0x01},
				RequiredPrivilegesList: []RequiredPrivileges{
					{
						Index:      []byte{0x01},
						Privileges: []byte{0x01},
						Kind:       "coin",
					},
				},
				WhitelistTransfersAccAddrs: []WhitelistTransfersAccAddrs{
					{
						Index:         []byte{0x01},
						Account:       sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t").Bytes(),
						IsWhitelisted: true,
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
					AdminAccount:                            "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
					AccountPrivilegesTokenCollectionCreator: "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
					AccountPrivilegesTokenCollectionId:      "id",
					DefaultPrivileges:                       make([]byte, 31),
					BaseDenom:                               "uom",
				},
			},
			true,
		},
		{
			"dup account privileges",
			GenesisState{
				Params: Params{
					AdminAccount:                            "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
					AccountPrivilegesTokenCollectionCreator: "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
					AccountPrivilegesTokenCollectionId:      "id",
					DefaultPrivileges:                       make([]byte, 32),
					BaseDenom:                               "uom",
				},
				AccountPrivilegesList: []AccountPrivileges{
					{
						Account:    sdk.MustAccAddressFromBech32("cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw").Bytes(),
						Privileges: []byte{0x01},
					},
					{
						Account:    sdk.MustAccAddressFromBech32("cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw").Bytes(),
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
