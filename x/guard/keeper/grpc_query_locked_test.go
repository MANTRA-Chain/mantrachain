package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestGRPCQueryLocked(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestLocked() {
	s.bankKeeper.EXPECT().GetDenomMetaData(gomock.Any(), gomock.Any()).Return(banktypes.Metadata{}, true)
	var (
		req    *types.QueryGetLockedRequest
		locked types.Locked
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryGetLockedResponse)
	}{
		{
			"success",
			func(index int, require *require.Assertions) {
				locked = types.Locked{
					Index: s.lkIndex,
					Kind:  s.rpKind.String(),
				}
				req = &types.QueryGetLockedRequest{
					Index: s.lkIndex,
					Kind:  s.rpKind.String(),
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryGetLockedResponse) {
				require.Equal(types.QueryGetLockedResponse{
					Index:  locked.Index,
					Locked: false,
					Kind:   locked.Kind,
				}, *res)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.Locked(gocontext.Background(), req)
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
