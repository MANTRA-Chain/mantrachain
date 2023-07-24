package keeper_test

import (
	"mantrachain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestCreateCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgCreateNftCollection
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgCreateNftCollection{
				Creator: "",
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id: "1",
				},
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should successfully create collection",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[0].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id: "1",
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail on creating an existing collection",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[0].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id: "1",
				},
			},
			expErr: true,
			errMsg: "nft collection already exists",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.CreateNftCollection(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.CollectionId, tc.req.Collection.Id)
				s.Require().EqualValues(req.CollectionCreator, tc.req.Creator)
			}
		})
	}
}
