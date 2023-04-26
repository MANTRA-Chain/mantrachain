package keeper_test

func (s *KeeperTestSuite) TestSetRequiredPrivileges() {
	index := []byte{0x01}
	privileges := []byte{0x02}

	s.guardKeeper.SetRequiredPrivileges(s.ctx, index, s.kind, privileges)

	has := s.guardKeeper.HasRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().True(has)

	actual, has := s.guardKeeper.GetRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().True(has)
	s.Require().EqualValues(privileges, actual)

	list := s.guardKeeper.GetRequiredPrivilegesMany(s.ctx, [][]byte{index}, s.kind)
	s.Require().EqualValues(privileges, list[0])
}

func (s *KeeperTestSuite) TestGetRequiredPrivileges() {
	index := []byte{0x02}
	privileges := []byte{0x02}

	has := s.guardKeeper.HasRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().False(has)

	actual, has := s.guardKeeper.GetRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().False(has)
	s.Require().EqualValues([]byte{}, actual)

	list := s.guardKeeper.GetRequiredPrivilegesMany(s.ctx, [][]byte{index}, s.kind)
	s.Require().EqualValues([]byte(nil), list[0])

	s.guardKeeper.SetRequiredPrivileges(s.ctx, index, s.kind, privileges)

	actual, has = s.guardKeeper.GetRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().True(has)
	s.Require().EqualValues(privileges, actual)

	list = s.guardKeeper.GetRequiredPrivilegesMany(s.ctx, [][]byte{index}, s.kind)
	s.Require().EqualValues(privileges, list[0])
}

func (s *KeeperTestSuite) TestRemoveRequiredPrivileges() {
	index := []byte{0x03}
	privileges := []byte{0x02}

	has := s.guardKeeper.HasRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().False(has)

	s.guardKeeper.SetRequiredPrivileges(s.ctx, index, s.kind, privileges)

	has = s.guardKeeper.HasRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().True(has)

	s.guardKeeper.RemoveRequiredPrivileges(s.ctx, index, s.kind)

	has = s.guardKeeper.HasRequiredPrivileges(s.ctx, index, s.kind)
	s.Require().False(has)
}
