package types_test

import (
	"testing"

	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func init() {
	accountAddressPrefix := "mantra"
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}

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
			desc: "negative proportion is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("-0.5", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "invalid bech32 address",
			genState: &types.GenesisState{
				Params: types.NewParams(types.DefaultMcaTax, "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qkb"),
			},
			valid: false,
		},
		{
			desc: "valid custom parameters",
			genState: &types.GenesisState{
				Params: types.NewParams("0.1", "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls"),
			},
			valid: true,
		},
		{
			desc: "mca tax greater than 1 is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("1.5", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "empty mca address is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams(types.DefaultMcaTax, ""),
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
