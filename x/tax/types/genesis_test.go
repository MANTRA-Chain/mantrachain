package types_test

import (
	"testing"

	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// init() is required to set the bech32 prefix for the test
// most likely, we can remove this post-launch by simply using the default prefixes and updating the addresses in the test.
// it is here so that we can test against the precise values that will be used in genesis.
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
		{
			desc: "mca tax of 0.5 is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("0.5", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "mca tax of 0 is valid",
			genState: &types.GenesisState{
				Params: types.NewParams("0", types.DefaultMcaAddress),
			},
			valid: true,
		},
		{
			desc: "mca tax of 1 is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams("1", types.DefaultMcaAddress),
			},
			valid: false,
		},
		{
			desc: "mca address with wrong prefix is invalid",
			genState: &types.GenesisState{
				Params: types.NewParams(types.DefaultMcaTax, "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzutu9"),
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
