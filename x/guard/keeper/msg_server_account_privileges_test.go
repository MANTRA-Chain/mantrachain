package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestUpdateAccountPrivileges() {
	testCases := []struct {
		name   string
		req    *types.MsgUpdateAccountPrivileges
		expErr bool
		errMsg string
	}{
		{
			name: "invalid account address",
			req: &types.MsgUpdateAccountPrivileges{
				Creator:    s.params.AdminAccount,
				Account:    "",
				Privileges: testPrivileges,
			},
			expErr: true,
			errMsg: "invalid account address",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateAccountPrivileges(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
