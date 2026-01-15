package keeper_test

import (
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	userAddr1 = "mantra19nj8hate4y2w0va698uygf2jj4y6xq42k8a7ed"
	userAddr2 = "mantra1587ulw8vg2dnnyw8xuntuvlrejy4rguz9hz530"
)

func (suite *KeeperTestSuite) TestAddBlacklistAccounts() {
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	testCases := []struct {
		name      string
		preRun    func()
		msg       *types.MsgAddBlacklistAccounts
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid transaction",
			msg: &types.MsgAddBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1, userAddr2},
			},
			expErr: false,
		},
		{
			name: "invalid Authority address",
			msg: &types.MsgAddBlacklistAccounts{
				Authority:         userAddr1,
				BlacklistAccounts: []string{userAddr2},
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "single invalid blacklist address",
			msg: &types.MsgAddBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{"Invalid Address"},
			},
			expErr:    true,
			expErrMsg: "invalid account",
		},
		{
			name: "contains invalid blacklist address",
			msg: &types.MsgAddBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1, "Invalid Address"},
			},
			expErr:    true,
			expErrMsg: "invalid account",
		},
		{
			name: "contains invalid blacklist address",
			msg: &types.MsgAddBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1, "Invalid Address"},
			},
			expErr:    true,
			expErrMsg: "invalid account",
		},
		{
			name: "added blacklist address",
			msg: &types.MsgAddBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1},
			},
			expErr:    true,
			expErrMsg: "has already been blacklisted",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.preRun != nil {
				tc.preRun()
			}
			resp, err := suite.msgServer.AddBlacklistAccounts(suite.ctx, tc.msg)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRemoveBlacklistAccounts() {
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	resp, err := suite.msgServer.AddBlacklistAccounts(suite.ctx, &types.MsgAddBlacklistAccounts{
		Authority:         govAddr,
		BlacklistAccounts: []string{userAddr1, userAddr2},
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	testCases := []struct {
		name      string
		preRun    func()
		msg       *types.MsgRemoveBlacklistAccounts
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid transaction",
			msg: &types.MsgRemoveBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1, userAddr2},
			},
			expErr: false,
		},
		{
			name: "invalid Authority address",
			msg: &types.MsgRemoveBlacklistAccounts{
				Authority:         userAddr1,
				BlacklistAccounts: []string{userAddr2},
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "single invalid blacklist address",
			msg: &types.MsgRemoveBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{"Invalid Address"},
			},
			expErr:    true,
			expErrMsg: "invalid account",
		},
		{
			name: "contains invalid blacklist address",
			msg: &types.MsgRemoveBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1, "Invalid Address"},
			},
			expErr:    true,
			expErrMsg: "invalid account",
		},
		{
			name: "contains invalid blacklist address",
			msg: &types.MsgRemoveBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1, "Invalid Address"},
			},
			expErr:    true,
			expErrMsg: "invalid account",
		},
		{
			name: "added blacklist address",
			msg: &types.MsgRemoveBlacklistAccounts{
				Authority:         govAddr,
				BlacklistAccounts: []string{userAddr1},
			},
			expErr:    true,
			expErrMsg: "is not blacklisted",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.preRun != nil {
				tc.preRun()
			}
			resp, err := suite.msgServer.RemoveBlacklistAccounts(suite.ctx, tc.msg)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
			}
		})
	}
}
