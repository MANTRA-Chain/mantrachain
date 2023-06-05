package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestUpdateRequiredPrivileges() {
	testCases := []struct {
		name   string
		req    *types.MsgUpdateRequiredPrivileges
		expErr bool
		errMsg string
	}{
		{
			name: "set required privileges",
			req: &types.MsgUpdateRequiredPrivileges{
				Creator:    s.params.AdminAccount,
				Index:      testIndex,
				Privileges: testPrivileges,
				Kind:       s.rpKind.String(),
			},
			expErr: false,
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateRequiredPrivileges(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
