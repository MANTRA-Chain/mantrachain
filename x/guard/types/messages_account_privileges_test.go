package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

func TestMsgUpdateAccountPrivileges(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgUpdateAccountPrivileges)
		expectedErr string // empty means no error
	}{
		{
			"happy case",
			func(msg *types.MsgUpdateAccountPrivileges) {},
			"",
		},
		{
			"invalid creator",
			func(msg *types.MsgUpdateAccountPrivileges) {
				msg.Creator = "invalidaddr"
			},
			"invalid creator address (decoding bech32 failed: invalid separator index -1): invalid address",
		},
		{
			"invalid account",
			func(msg *types.MsgUpdateAccountPrivileges) {
				msg.Account = "invalidaddr"
			},
			"invalid account address (decoding bech32 failed: invalid separator index -1): invalid address",
		},
		{
			"wrong size privileges",
			func(msg *types.MsgUpdateAccountPrivileges) {
				msg.Privileges = make([]byte, 31)
			},
			"invalid privileges length (31): invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgUpdateAccountPrivileges(utils.TestAddress(0).String(), utils.TestAddress(1).String(), make([]byte, 32))
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgUpdateAccountPrivileges, msg.Type())
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
