package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/AumegaChain/aumega/x/guard/types"
)

func TestGRPCQueryRequiredPrivileges(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestRequiredPrivileges() {
	var (
		req                    *types.QueryGetRequiredPrivilegesRequest
		requiredPrivilegesResp types.QueryGetRequiredPrivilegesResponse
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryGetRequiredPrivilegesResponse)
	}{
		{
			"success",
			func(index int, require *require.Assertions) {
				requiredPrivilegesResp = types.QueryGetRequiredPrivilegesResponse{
					Index:      []byte{0x01},
					Privileges: []byte{0x02},
					Kind:       s.rpKind.String(),
				}
				req = &types.QueryGetRequiredPrivilegesRequest{
					Index: []byte{0x01},
					Kind:  s.rpKind.String(),
				}
				s.TestSetRequiredPrivileges()
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryGetRequiredPrivilegesResponse) {
				require.Equal(requiredPrivilegesResp, *res)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.RequiredPrivileges(gocontext.Background(), req)
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
