package keeper_test

import (
	"strconv"
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/txfees/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestFeeTokenMsgServerCreate(t *testing.T) {
	k, ctx := TxfeesKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateFeeToken{
			Denom: strconv.Itoa(i),
		}
		_, err := srv.CreateFeeToken(ctx, expected)
		require.NoError(t, err)
		_, found := k.GetFeeToken(ctx,
			expected.Denom,
		)
		require.True(t, found)
	}
}

func TestFeeTokenMsgServerUpdate(t *testing.T) {
	tests := []struct {
		desc    string
		request *types.MsgUpdateFeeToken
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateFeeToken{
				Denom: strconv.Itoa(0),
			},
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateFeeToken{
				Denom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := TxfeesKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			expected := &types.MsgCreateFeeToken{
				Denom: strconv.Itoa(0),
			}
			_, err := srv.CreateFeeToken(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateFeeToken(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetFeeToken(ctx,
					expected.Denom,
				)
				require.True(t, found)
			}
		})
	}
}

func TestFeeTokenMsgServerDelete(t *testing.T) {
	tests := []struct {
		desc    string
		request *types.MsgDeleteFeeToken
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteFeeToken{
				Denom: strconv.Itoa(0),
			},
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteFeeToken{
				Denom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := TxfeesKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)

			_, err := srv.CreateFeeToken(ctx, &types.MsgCreateFeeToken{
				Denom: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteFeeToken(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetFeeToken(ctx,
					tc.request.Denom,
				)
				require.False(t, found)
			}
		})
	}
}
