package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestGRPCQueryParams(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestParams() {
	var (
		req        *types.QueryParamsRequest
		paramsResp types.Params
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryParamsResponse)
	}{
		{
			"success",
			func(index int, require *require.Assertions) {
				paramsResp = s.params
				req = &types.QueryParamsRequest{}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryParamsResponse) {
				require.Equal(paramsResp, res.Params)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.Params(gocontext.Background(), req)
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
