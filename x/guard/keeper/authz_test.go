package keeper_test

import (
	"math/big"

	"github.com/MANTRA-Finance/aumega/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var authz = "module:coinfactory:CreateDenom"

func (s *KeeperTestSuite) TestValidateAuthz() {
	goCtx := sdk.WrapSDKContext(s.ctx)
	authzBytes := []byte(authz)

	err := s.guardKeeper.CheckHasAuthz(s.ctx, s.testAdminAccount, authz)
	s.Require().NoError(err)

	err = s.guardKeeper.CheckHasAuthz(s.ctx, s.addrs[0].String(), authz)
	s.Require().Contains(err.Error(), "authz required privileges not found")

	privileges := types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      authzBytes,
		Privileges: privileges.Bytes(),
		Kind:       "authz",
	})
	s.Require().NoError(err)
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckHasAuthz(s.ctx, s.addrs[0].String(), authz)
	s.Require().NoError(err)

	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      authzBytes,
		Privileges: privileges.Bytes(),
		Kind:       "authz",
	})
	s.Require().NoError(err)
	privileges = privileges.SwitchOff([]*big.Int{big.NewInt(64)}).SwitchOn([]*big.Int{big.NewInt(65)})
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckHasAuthz(s.ctx, s.addrs[0].String(), authz)
	s.Require().Contains(err.Error(), "insufficient privileges")
}
