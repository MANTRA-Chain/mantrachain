package keeper_test

import (
	"math/big"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
)

func (s *KeeperTestSuite) TestCheckCanTransferCoins() {
	goCtx := sdk.WrapSDKContext(s.ctx)

	err := s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], nil)
	s.Require().NoError(err)

	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin("mantra", 1000000000000000000)))
	s.Require().Contains(err.Error(), "coin required privileges not found")

	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(types.DefaultBaseDenom, 1000000000000000000)))
	s.Require().NoError(err)

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[0], true).Times(1)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().NoError(err)

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().Contains(err.Error(), "coin required privileges not found")

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: s.defaultPrivileges,
		Kind:       "coin",
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().Contains(err.Error(), "coin required privileges not set")

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges := types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().Contains(err.Error(), "account privileges not set")

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	privileges = privileges.SwitchOff([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().Contains(err.Error(), "account privileges not set")

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	privileges = privileges.SwitchOff([]*big.Int{big.NewInt(64)}).SwitchOn([]*big.Int{big.NewInt(65)})
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().Contains(err.Error(), "insufficient privileges")

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().NoError(err)

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)}).SwitchOff([]*big.Int{big.NewInt(0)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	privileges = privileges.SwitchOn([]*big.Int{big.NewInt(0)})
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().NoError(err)

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64), big.NewInt(65)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	privileges = privileges.SwitchOff([]*big.Int{big.NewInt(65)})
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().NoError(err)

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	privileges = privileges.SwitchOff([]*big.Int{big.NewInt(0)})
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().Contains(err.Error(), "insufficient privileges")

	s.coinFactoryKeeper.EXPECT().GetAdmin(gomock.Any(), gomock.Any()).Return(s.addrs[1], true).Times(1)
	s.nftKeeper.EXPECT().GetOwner(gomock.Any(), gomock.Any(), gomock.Any()).Return(s.addrs[0]).Times(1)
	privileges = types.PrivilegesFromBytes(s.defaultPrivileges).SwitchOn([]*big.Int{big.NewInt(64)}).SwitchOff([]*big.Int{big.NewInt(0)})
	_, err = s.msgServer.UpdateRequiredPrivileges(goCtx, &types.MsgUpdateRequiredPrivileges{
		Creator:    s.testAdminAccount,
		Index:      s.lkIndex,
		Privileges: privileges.Bytes(),
		Kind:       "coin",
	})
	s.Require().NoError(err)
	_, err = s.msgServer.UpdateAccountPrivileges(goCtx, &types.MsgUpdateAccountPrivileges{
		Creator:    s.testAdminAccount,
		Account:    s.addrs[0].String(),
		Privileges: privileges.Bytes(),
	})
	s.Require().NoError(err)
	err = s.guardKeeper.CheckCanTransferCoins(s.ctx, s.addrs[0], sdk.NewCoins(sdk.NewInt64Coin(string(s.lkIndex), 1000000000000000000)))
	s.Require().NoError(err)
}
