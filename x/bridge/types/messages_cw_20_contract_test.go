package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/LimeChain/mantrachain/testutil/sample"
)

func TestMsgCreateCw20Contract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateCw20Contract
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateCw20Contract{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateCw20Contract{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateCw20Contract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateCw20Contract
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateCw20Contract{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateCw20Contract{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteCw20Contract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteCw20Contract
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteCw20Contract{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteCw20Contract{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
