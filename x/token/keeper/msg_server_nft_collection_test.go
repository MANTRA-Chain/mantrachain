package keeper_test

import (
	"github.com/MANTRA-Finance/aumega/x/token/types"

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
			name: "should successfully create collection for another creator",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[1].String(),
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
		{
			name: "should fail on creating a restricted collection when not an admin",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[0].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id:             "2",
					RestrictedNfts: true,
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		},
		{
			name: "should fail on creating a collection with empty id and true for an opened flag",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[0].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id:     "",
					Opened: true,
				},
			},
			expErr: true,
			errMsg: "nft collection id provided is invalid",
		},
		{
			name: "should fail on creating a collection with default for an id and true for an opened flag",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[0].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id:     "default",
					Opened: true,
				},
			},
			expErr: true,
			errMsg: "nft collection id provided is invalid",
		},
		{
			name: "should successfully create an opened collection",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[0].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id:     "2",
					Opened: true,
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully create an opened collection with soul-bonded nfts",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[0].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id:             "3",
					Opened:         true,
					SoulBondedNfts: true,
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail on creating a collection with non-existing category",
			req: &types.MsgCreateNftCollection{
				Creator: s.addrs[1].String(),
				Collection: &types.MsgCreateNftCollectionMetadata{
					Id:       "2",
					Category: "invalid",
				},
			},
			expErr: true,
			errMsg: "invalid nft collection category",
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
				if tc.req.Collection.Id == "" {
					s.Require().EqualValues(req.CollectionId, "default")
				} else {
					s.Require().EqualValues(req.CollectionId, tc.req.Collection.Id)
				}
				s.Require().EqualValues(req.CollectionCreator, tc.req.Creator)
			}
		})
	}
}
