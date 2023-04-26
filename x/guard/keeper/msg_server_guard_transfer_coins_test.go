package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestUpdateGuardTransferCoins() {
	testCases := []struct {
		name   string
		req    *types.MsgUpdateGuardTransferCoins
		expErr bool
		errMsg string
	}{
		{
			name: "set guard transfer",
			req: &types.MsgUpdateGuardTransferCoins{
				Creator: testAccount,
				Enabled: true,
			},
			expErr: false,
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateGuardTransferCoins(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
