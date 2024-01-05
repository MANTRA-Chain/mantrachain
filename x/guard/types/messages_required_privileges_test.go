package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	utils "github.com/MANTRA-Finance/aumega/types"
	"github.com/MANTRA-Finance/aumega/x/guard/types"
)

func TestMsgUpdateRequiredPrivileges(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgUpdateRequiredPrivileges)
		expectedErr string // empty means no error
	}{
		{
			"happy case",
			func(msg *types.MsgUpdateRequiredPrivileges) {},
			"",
		},
		{
			"invalid creator",
			func(msg *types.MsgUpdateRequiredPrivileges) {
				msg.Creator = "invalidaddr"
			},
			"invalid creator address (decoding bech32 failed: invalid separator index -1): invalid address",
		},
		{
			"invalid index",
			func(msg *types.MsgUpdateRequiredPrivileges) {
				msg.Index = []byte{}
			},
			"index should not be empty: invalid request",
		},
		{
			"wrong size privileges",
			func(msg *types.MsgUpdateRequiredPrivileges) {
				msg.Privileges = make([]byte, 31)
			},
			"invalid privileges length (31): invalid request",
		},
		{
			"invalid kind",
			func(msg *types.MsgUpdateRequiredPrivileges) {
				msg.Kind = "invalid"
			},
			"kind is invalid: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgUpdateRequiredPrivileges(utils.TestAddress(0).String(), []byte{0x01}, make([]byte, 32), "coin")
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgUpdateRequiredPrivileges, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
