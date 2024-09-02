package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestGRPCQueryGuardTransferCoins(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestGuardTransferCoins() {
	var (
		req                    *types.QueryGetGuardTransferCoinsRequest
		guardTransferCoinsResp bool
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryGetGuardTransferCoinsResponse)
	}{
		{
			"success",
			func(index int, require *require.Assertions) {
				guardTransferCoinsResp = true
				req = &types.QueryGetGuardTransferCoinsRequest{}
				s.TestSetGuardTransferCoins()
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryGetGuardTransferCoinsResponse) {
				require.Equal(guardTransferCoinsResp, res.GuardTransferCoins)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.GuardTransferCoins(gocontext.Background(), req)
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
