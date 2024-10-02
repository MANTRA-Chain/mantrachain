package keeper_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/tax/keeper"
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

func TestMsgUpdateParams(t *testing.T) {
	k, ctx, _ := keepertest.TaxKeeper(t)
	ms := keeper.NewMsgServerImpl(k)

	params := types.DefaultParams()
	require.NoError(t, k.Params.Set(ctx, params))

	// default params
	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgUpdateParams{
				Authority:  "invalid",
				McaTax:     "",
				McaAddress: "",
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "update mca tax",
			input: &types.MsgUpdateParams{
				Authority:  "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
				McaTax:     "0.200000000000000000",
				McaAddress: "",
			},
			expErr: false,
		},
		{
			name: "update mca address",
			input: &types.MsgUpdateParams{
				Authority:  "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
				McaTax:     "",
				McaAddress: "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
			},
			expErr: false,
		},
		{
			name: "old authority address no longer work",
			input: &types.MsgUpdateParams{
				Authority:  "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
				McaTax:     "",
				McaAddress: "",
			},
			expErr:    true,
			expErrMsg: "invalid sender; expected mcaAddress",
		},
		{
			name: "update both",
			input: &types.MsgUpdateParams{
				Authority:  "mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka",
				McaTax:     "0.200000000000000000",
				McaAddress: "mantra15m77x4pe6w9vtpuqm22qxu0ds7vn4ehzwx8pls",
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.UpdateParams(ctx, tc.input)

			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				params, err := k.Params.Get(ctx)
				require.NoError(t, err)
				if tc.input.McaTax != "" {
					require.Equal(t, tc.input.McaTax, params.McaTax.String())
				}
				if tc.input.McaAddress != "" {
					require.Equal(t, tc.input.McaAddress, params.McaAddress)
				}
			}
		})
	}
}
