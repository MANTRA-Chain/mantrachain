package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestGRPCQueryAccountPrivileges(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestAccountPrivileges() {
	var (
		req                   *types.QueryGetAccountPrivilegesRequest
		accountPrivilegesResp types.QueryGetAccountPrivilegesResponse
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryGetAccountPrivilegesResponse)
	}{
		{
			"success",
			func(index int, require *require.Assertions) {
				accountPrivilegesResp = types.QueryGetAccountPrivilegesResponse{
					Account:    s.addrs[0].String(),
					Privileges: []byte{0x02},
				}
				req = &types.QueryGetAccountPrivilegesRequest{
					Account: s.addrs[0].String(),
				}
				s.TestSetAccountPrivileges()
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryGetAccountPrivilegesResponse) {
				require.Equal(accountPrivilegesResp, *res)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.AccountPrivileges(gocontext.Background(), req)
			if tc.expError == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expError)
			}
			tc.postTest(index, require, result)
		})
	}
}
