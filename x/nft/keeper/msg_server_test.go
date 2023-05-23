package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ExpClass = types.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}

	ExpNFT = types.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
)

func (s *TestSuite) TestSend() {
	err := s.nftKeeper.SaveClass(s.ctx, ExpClass)
	s.Require().NoError(err)

	actual, has := s.nftKeeper.GetClass(s.ctx, testClassID)
	s.Require().True(has)
	s.Require().EqualValues(ExpClass, actual)

	err = s.nftKeeper.Mint(s.ctx, ExpNFT, s.addrs[0])
	s.Require().NoError(err)

	expGenesis := &types.GenesisState{
		Classes: []*types.Class{&ExpClass},
		Entries: []*types.Entry{{
			Owner: s.addrs[0].String(),
			Nfts:  []*types.NFT{&ExpNFT},
		}},
	}
	genesis := s.nftKeeper.ExportGenesis(s.ctx)
	s.Require().Equal(expGenesis, genesis)

	testCases := []struct {
		name   string
		req    *types.MsgSend
		expErr bool
		errMsg string
	}{
		{
			name: "transfer nft disabled",
			req: &types.MsgSend{
				ClassId:  testClassID,
				Id:       "",
				Sender:   s.addrs[0].String(),
				Receiver: s.addrs[1].String(),
			},
			expErr: true,
			errMsg: "transfer nft disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.nftKeeper.Send(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
