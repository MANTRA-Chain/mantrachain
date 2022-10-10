package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

    keepertest "github.com/LimeChain/mantrachain/testutil/keeper"
    "github.com/LimeChain/mantrachain/x/bridge/keeper"
    "github.com/LimeChain/mantrachain/x/bridge/types"
)

func TestCw20ContractMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.BridgeKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	expected := &types.MsgCreateCw20Contract{Creator: creator}
    _, err := srv.CreateCw20Contract(wctx, expected)
    require.NoError(t, err)
    rst, found := k.GetCw20Contract(ctx)
    require.True(t, found)
    require.Equal(t, expected.Creator, rst.Creator)
}

func TestCw20ContractMsgServerUpdate(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateCw20Contract
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateCw20Contract{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateCw20Contract{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.BridgeKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateCw20Contract{Creator: creator}
			_, err := srv.CreateCw20Contract(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateCw20Contract(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetCw20Contract(ctx)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestCw20ContractMsgServerDelete(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteCw20Contract
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteCw20Contract{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteCw20Contract{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.BridgeKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateCw20Contract(wctx, &types.MsgCreateCw20Contract{Creator: creator})
			require.NoError(t, err)
			_, err = srv.DeleteCw20Contract(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetCw20Contract(ctx)
				require.False(t, found)
			}
		})
	}
}
