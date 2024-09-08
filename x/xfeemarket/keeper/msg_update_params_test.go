package keeper_test

import (
	"testing"

	keepertest "github.com/MANTRA-Chain/mantrachain/testutil/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateParams(t *testing.T) {
	k, ctx, _ := keepertest.XfeemarketKeeper(t, nil)
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
				Authority: "invalid",
				Params:    params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "send enabled param",
			input: &types.MsgUpdateParams{
				Authority: k.GetAuthority(),
				Params:    types.Params{},
			},
			expErr: false,
		},
		{
			name: "all good",
			input: &types.MsgUpdateParams{
				Authority: k.GetAuthority(),
				Params:    params,
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
			}
		})
	}
}
