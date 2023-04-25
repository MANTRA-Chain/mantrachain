package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestSetAccountPrivileges() {
	except := types.AccountPrivileges{
		Account:    s.TestAccs[0],
		Privileges: []byte{0x02},
	}

	s.keeper.SetAccountPrivileges(s.ctx, except.Account, except.Privileges)

	has := s.keeper.HasAccountPrivileges(s.ctx, except.Account)
	s.Require().True(has)

	actual, has := s.keeper.GetAccountPrivileges(s.ctx, except.Account, []byte{0x01})
	s.Require().True(has)
	s.Require().EqualValues(except.Privileges, actual)

	list := s.keeper.GetAccountPrivilegesMany(s.ctx, []sdk.AccAddress{except.Account}, []byte{0x01})
	s.Require().EqualValues(except.Privileges, list[0])

}
