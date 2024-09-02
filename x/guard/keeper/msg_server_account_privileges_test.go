package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
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
		}, {
			name: "set account privileges",
			req: &types.MsgUpdateAccountPrivileges{
				Creator:    s.params.AdminAccount,
				Account:    s.testAccount,
				Privileges: testPrivileges,
			},
			expErr: false,
		}, {
			name: "not an admin",
			req: &types.MsgUpdateAccountPrivileges{
				Creator:    s.addrs[0].String(),
				Account:    s.testAccount,
				Privileges: testPrivileges,
			},
			expErr: true,
			errMsg: "unauthorized",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateAccountPrivileges(s.ctx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
