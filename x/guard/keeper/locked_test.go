package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

func (s *KeeperTestSuite) TestSetLocked() {
	index := []byte{0x01}

	s.guardKeeper.SetLocked(s.ctx, index, s.lkKind)

	has := s.guardKeeper.HasLocked(s.ctx, index, s.lkKind)
	s.Require().True(has)

	actual, has := s.guardKeeper.GetLocked(s.ctx, index, s.lkKind)
	s.Require().True(has)
	s.Require().EqualValues(types.Placeholder, actual)
}

func (s *KeeperTestSuite) TestGetLocked() {
	index := []byte{0x02}

	has := s.guardKeeper.HasLocked(s.ctx, index, s.lkKind)
	s.Require().False(has)

	actual, has := s.guardKeeper.GetLocked(s.ctx, index, s.lkKind)
	s.Require().False(has)
	s.Require().EqualValues([]byte{}, actual)

	s.guardKeeper.SetLocked(s.ctx, index, s.lkKind)

	actual, has = s.guardKeeper.GetLocked(s.ctx, index, s.lkKind)
	s.Require().True(has)
	s.Require().EqualValues(types.Placeholder, actual)
}

func (s *KeeperTestSuite) TestRemoveLocked() {
	index := []byte{0x03}

	has := s.guardKeeper.HasLocked(s.ctx, index, s.lkKind)
	s.Require().False(has)

	s.guardKeeper.SetLocked(s.ctx, index, s.lkKind)

	has = s.guardKeeper.HasLocked(s.ctx, index, s.lkKind)
	s.Require().True(has)

	s.guardKeeper.RemoveLocked(s.ctx, index, s.lkKind)

	has = s.guardKeeper.HasLocked(s.ctx, index, s.lkKind)
	s.Require().False(has)
}
