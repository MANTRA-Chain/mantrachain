package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/MANTRA-Finance/mantrachain/x/nft/types"
)

func TestGRPCQuery(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestBalance() {
	var req *types.QueryBalanceRequest
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		balance  uint64
		postTest func(index int, require *require.Assertions, res *types.QueryBalanceResponse, expBalance uint64)
	}{
		{
			"fail empty ClassId",
			func(index int, require *require.Assertions) {
				req = &types.QueryBalanceRequest{}
			},
			types.ErrEmptyClassID.Error(),
			0,
			func(index int, require *require.Assertions, res *types.QueryBalanceResponse, expBalance uint64) {},
		},
		{
			"fail invalid Owner addr",
			func(index int, require *require.Assertions) {
				req = &types.QueryBalanceRequest{
					ClassId: testClassID,
					Owner:   "owner",
				}
			},
			"decoding bech32 failed",
			0,
			func(index int, require *require.Assertions, res *types.QueryBalanceResponse, expBalance uint64) {},
		},
		{
			"Success",
			func(index int, require *require.Assertions) {
				s.TestMint()
				req = &types.QueryBalanceRequest{
					ClassId: testClassID,
					Owner:   s.addrs[0].String(),
				}
			},
			"",
			2,
			func(index int, require *require.Assertions, res *types.QueryBalanceResponse, expBalance uint64) {
				require.Equal(res.Amount, expBalance, "the error occurred on:%d", index)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.Balance(gocontext.Background(), req)
			if tc.expError == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expError)
			}
			tc.postTest(index, require, result, tc.balance)
		})
	}
}

func (s *TestSuite) TestOwner() {
	var (
		req   *types.QueryOwnerRequest
		owner string
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryOwnerResponse)
	}{
		{
			"fail empty ClassId",
			func(index int, require *require.Assertions) {
				req = &types.QueryOwnerRequest{
					Id: testID,
				}
			},
			types.ErrEmptyClassID.Error(),
			func(index int, require *require.Assertions, res *types.QueryOwnerResponse) {},
		},
		{
			"fail empty nft id",
			func(index int, require *require.Assertions) {
				req = &types.QueryOwnerRequest{
					ClassId: testClassID,
				}
			},
			types.ErrEmptyNFTID.Error(),
			func(index int, require *require.Assertions, res *types.QueryOwnerResponse) {},
		},
		{
			"success but nft id not exist",
			func(index int, require *require.Assertions) {
				req = &types.QueryOwnerRequest{
					ClassId: testClassID,
					Id:      "kitty2",
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryOwnerResponse) {
				require.Equal(res.Owner, owner, "the error occurred on:%d", index)
			},
		},
		{
			"success but class id not exist",
			func(index int, require *require.Assertions) {
				req = &types.QueryOwnerRequest{
					ClassId: "kitty1",
					Id:      testID,
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryOwnerResponse) {
				require.Equal(res.Owner, owner, "the error occurred on:%d", index)
			},
		},
		{
			"Success",
			func(index int, require *require.Assertions) {
				s.TestMint()
				req = &types.QueryOwnerRequest{
					ClassId: testClassID,
					Id:      testID,
				}
				owner = s.addrs[0].String()
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryOwnerResponse) {
				require.Equal(res.Owner, owner, "the error occurred on:%d", index)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.Owner(gocontext.Background(), req)
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

func (s *TestSuite) TestSupply() {
	var req *types.QuerySupplyRequest
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		supply   uint64
		postTest func(index int, require *require.Assertions, res *types.QuerySupplyResponse, supply uint64)
	}{
		{
			"fail empty ClassId",
			func(index int, require *require.Assertions) {
				req = &types.QuerySupplyRequest{}
			},
			types.ErrEmptyClassID.Error(),
			0,
			func(index int, require *require.Assertions, res *types.QuerySupplyResponse, supply uint64) {},
		},
		{
			"success but class id not exist",
			func(index int, require *require.Assertions) {
				req = &types.QuerySupplyRequest{
					ClassId: "kitty1",
				}
			},
			"",
			0,
			func(index int, require *require.Assertions, res *types.QuerySupplyResponse, supply uint64) {
				require.Equal(res.Amount, supply, "the error occurred on:%d", index)
			},
		},
		{
			"success but supply equal zero",
			func(index int, require *require.Assertions) {
				req = &types.QuerySupplyRequest{
					ClassId: testClassID,
				}
				s.TestSaveClass()
			},
			"",
			0,
			func(index int, require *require.Assertions, res *types.QuerySupplyResponse, supply uint64) {
				require.Equal(res.Amount, supply, "the error occurred on:%d", index)
			},
		},
		{
			"Success",
			func(index int, require *require.Assertions) {
				n := types.NFT{
					ClassId: testClassID,
					Id:      testID,
					Uri:     testURI,
				}
				err := s.nftKeeper.Mint(s.ctx, n, s.addrs[0])
				require.NoError(err, "the error occurred on:%d", index)

				req = &types.QuerySupplyRequest{
					ClassId: testClassID,
				}
			},
			"",
			1,
			func(index int, require *require.Assertions, res *types.QuerySupplyResponse, supply uint64) {
				require.Equal(res.Amount, supply, "the error occurred on:%d", index)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.Supply(gocontext.Background(), req)
			if tc.expError == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expError)
			}
			tc.postTest(index, require, result, tc.supply)
		})
	}
}

func (s *TestSuite) TestNFTs() {
	var (
		req  *types.QueryNFTsRequest
		nfts []*types.NFT
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryNFTsResponse)
	}{
		{
			"fail empty Owner and ClassId",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTsRequest{}
			},
			"must provide at least one of classID or owner",
			func(index int, require *require.Assertions, res *types.QueryNFTsResponse) {},
		},
		{
			"success,empty ClassId and no nft",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTsRequest{
					Owner: s.addrs[1].String(),
				}
				s.TestSaveClass()
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryNFTsResponse) {
				require.Len(res.Nfts, 0, "the error occurred on:%d", index)
			},
		},
		{
			"success, empty Owner and class id not exist",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTsRequest{
					ClassId: "kitty1",
				}
				n := types.NFT{
					ClassId: testClassID,
					Id:      testID,
					Uri:     testURI,
				}
				err := s.nftKeeper.Mint(s.ctx, n, s.addrs[0])
				require.NoError(err, "the error occurred on:%d", index)
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryNFTsResponse) {
				require.Len(res.Nfts, 0, "the error occurred on:%d", index)
			},
		},
		{
			"Success,query by owner",
			func(index int, require *require.Assertions) {
				err := s.nftKeeper.SaveClass(s.ctx, types.Class{
					Id: "MyKitty",
				})
				require.NoError(err)

				nfts = []*types.NFT{}
				for i := 0; i < 5; i++ {
					n := types.NFT{
						ClassId: "MyKitty",
						Id:      fmt.Sprintf("MyCat%d", i),
					}
					err := s.nftKeeper.Mint(s.ctx, n, s.addrs[2])
					require.NoError(err)
					nfts = append(nfts, &n)
				}

				req = &types.QueryNFTsRequest{
					Owner: s.addrs[2].String(),
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryNFTsResponse) {
				require.EqualValues(res.Nfts, nfts, "the error occurred on:%d", index)
			},
		},
		{
			"Success,query by classID",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTsRequest{
					ClassId: "MyKitty",
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryNFTsResponse) {
				require.EqualValues(res.Nfts, nfts, "the error occurred on:%d", index)
			},
		},
		{
			"Success,query by classId and owner",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTsRequest{
					ClassId: testClassID,
					Owner:   s.addrs[0].String(),
				}
				nfts = []*types.NFT{
					{
						ClassId: testClassID,
						Id:      testID,
						Uri:     testURI,
					},
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryNFTsResponse) {
				require.Equal(res.Nfts, nfts, "the error occurred on:%d", index)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.NFTs(gocontext.Background(), req)
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

func (s *TestSuite) TestNFT() {
	var (
		req    *types.QueryNFTRequest
		expNFT types.NFT
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryNFTResponse)
	}{
		{
			"fail empty ClassId",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTRequest{}
			},
			types.ErrEmptyClassID.Error(),
			func(index int, require *require.Assertions, res *types.QueryNFTResponse) {},
		},
		{
			"fail empty nft id",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTRequest{
					ClassId: testClassID,
				}
			},
			types.ErrEmptyNFTID.Error(),
			func(index int, require *require.Assertions, res *types.QueryNFTResponse) {},
		},
		{
			"fail ClassId not exist",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTRequest{
					ClassId: "kitty1",
					Id:      testID,
				}
				s.TestMint()
			},
			"not found nft",
			func(index int, require *require.Assertions, res *types.QueryNFTResponse) {},
		},
		{
			"fail nft id not exist",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTRequest{
					ClassId: testClassID,
					Id:      "kitty2",
				}
			},
			"not found nft",
			func(index int, require *require.Assertions, res *types.QueryNFTResponse) {},
		},
		{
			"success",
			func(index int, require *require.Assertions) {
				req = &types.QueryNFTRequest{
					ClassId: testClassID,
					Id:      testID,
				}
				expNFT = types.NFT{
					ClassId: testClassID,
					Id:      testID,
					Uri:     testURI,
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryNFTResponse) {
				require.Equal(*res.Nft, expNFT, "the error occurred on:%d", index)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.NFT(gocontext.Background(), req)
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

func (s *TestSuite) TestClass() {
	var (
		req   *types.QueryClassRequest
		class types.Class
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryClassResponse)
	}{
		{
			"fail empty ClassId",
			func(index int, require *require.Assertions) {
				req = &types.QueryClassRequest{}
			},
			types.ErrEmptyClassID.Error(),
			func(index int, require *require.Assertions, res *types.QueryClassResponse) {},
		},
		{
			"fail ClassId not exist",
			func(index int, require *require.Assertions) {
				req = &types.QueryClassRequest{
					ClassId: "kitty1",
				}
				s.TestSaveClass()
			},
			"not found class",
			func(index int, require *require.Assertions, res *types.QueryClassResponse) {},
		},
		{
			"success",
			func(index int, require *require.Assertions) {
				class = types.Class{
					Id:          testClassID,
					Name:        testClassName,
					Symbol:      testClassSymbol,
					Description: testClassDescription,
					Uri:         testClassURI,
					UriHash:     testClassURIHash,
				}
				req = &types.QueryClassRequest{
					ClassId: testClassID,
				}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryClassResponse) {
				require.Equal(*res.Class, class, "the error occurred on:%d", index)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.Class(gocontext.Background(), req)
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

func (s *TestSuite) TestClasses() {
	var (
		req     *types.QueryClassesRequest
		classes []types.Class
	)
	testCases := []struct {
		msg      string
		malleate func(index int, require *require.Assertions)
		expError string
		postTest func(index int, require *require.Assertions, res *types.QueryClassesResponse)
	}{
		{
			"success Class not exist",
			func(index int, require *require.Assertions) {
				req = &types.QueryClassesRequest{}
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryClassesResponse) {
				require.Len(res.Classes, 0)
			},
		},
		{
			"success",
			func(index int, require *require.Assertions) {
				req = &types.QueryClassesRequest{}
				classes = []types.Class{
					{
						Id:          testClassID,
						Name:        testClassName,
						Symbol:      testClassSymbol,
						Description: testClassDescription,
						Uri:         testClassURI,
						UriHash:     testClassURIHash,
					},
				}
				s.TestSaveClass()
			},
			"",
			func(index int, require *require.Assertions, res *types.QueryClassesResponse) {
				require.Len(res.Classes, 1, "the error occurred on:%d", index)
				require.Equal(*res.Classes[0], classes[0], "the error occurred on:%d", index)
			},
		},
	}
	for index, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := s.Require()
			tc.malleate(index, require)
			result, err := s.queryClient.Classes(gocontext.Background(), req)
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
