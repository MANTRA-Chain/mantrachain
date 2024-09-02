package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestSetAccountPrivileges() {
	account := s.addrs[0]
	privileges := []byte{0x02}

	s.guardKeeper.SetAccountPrivileges(s.ctx, account, privileges)

	has := s.guardKeeper.HasAccountPrivileges(s.ctx, account)
	s.Require().True(has)

	actual, has := s.guardKeeper.GetAccountPrivileges(s.ctx, account, s.defaultPrivileges)
	s.Require().True(has)
	s.Require().EqualValues(privileges, actual)

	list := s.guardKeeper.GetAccountPrivilegesMany(s.ctx, []sdk.AccAddress{account}, s.defaultPrivileges)
	s.Require().EqualValues(privileges, list[0])
}

func (s *KeeperTestSuite) TestGetAccountPrivileges() {
	account := s.addrs[1]

	has := s.guardKeeper.HasAccountPrivileges(s.ctx, account)
	s.Require().False(has)

	actual, has := s.guardKeeper.GetAccountPrivileges(s.ctx, account, nil)
	s.Require().False(has)
	s.Require().EqualValues([]byte{}, actual)

	list := s.guardKeeper.GetAccountPrivilegesMany(s.ctx, []sdk.AccAddress{account}, nil)
	s.Require().EqualValues([]byte(nil), list[0])

	actual, has = s.guardKeeper.GetAccountPrivileges(s.ctx, account, s.defaultPrivileges)
	s.Require().True(has)
	s.Require().EqualValues(s.defaultPrivileges, actual)

	list = s.guardKeeper.GetAccountPrivilegesMany(s.ctx, []sdk.AccAddress{account}, s.defaultPrivileges)
	s.Require().EqualValues(s.defaultPrivileges, list[0])
}

func (s *KeeperTestSuite) TestRemoveAccountPrivileges() {
	account := s.addrs[2]
	privileges := []byte{0x02}

	has := s.guardKeeper.HasAccountPrivileges(s.ctx, account)
	s.Require().False(has)

	s.guardKeeper.SetAccountPrivileges(s.ctx, account, privileges)

	has = s.guardKeeper.HasAccountPrivileges(s.ctx, account)
	s.Require().True(has)

	s.guardKeeper.RemoveAccountPrivileges(s.ctx, account)

	has = s.guardKeeper.HasAccountPrivileges(s.ctx, account)
	s.Require().False(has)
}
