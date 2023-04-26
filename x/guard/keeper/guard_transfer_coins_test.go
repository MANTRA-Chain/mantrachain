package keeper_test

func (s *KeeperTestSuite) TestSetuardTransferCoins() {
	s.guardKeeper.SetGuardTransferCoins(s.ctx)

	has := s.guardKeeper.HasGuardTransferCoins(s.ctx)
	s.Require().True(has)
}

func (s *KeeperTestSuite) TestRemoveGuardTransferCoins() {
	has := s.guardKeeper.HasGuardTransferCoins(s.ctx)
	s.Require().False(has)

	s.guardKeeper.SetGuardTransferCoins(s.ctx)

	has = s.guardKeeper.HasGuardTransferCoins(s.ctx)
	s.Require().True(has)

	s.guardKeeper.RemoveGuardTransferCoins(s.ctx)

	has = s.guardKeeper.HasGuardTransferCoins(s.ctx)
	s.Require().False(has)
}
