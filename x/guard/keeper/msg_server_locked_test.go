package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestUpdateLocked() {
	testCases := []struct {
		name   string
		req    *types.MsgUpdateLocked
		expErr bool
		errMsg string
	}{
		{
			name: "invalid denom",
			req: &types.MsgUpdateLocked{
				Creator: testAccount,
				Index:   []byte{0x01},
				Locked:  true,
				Kind:    s.lkKind.String(),
			},
			expErr: true,
			errMsg: "invalid denom",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateLocked(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
