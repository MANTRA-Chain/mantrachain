package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestMintNfts() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgMintNfts{
				Creator:           "",
				Receiver:          "",
				CollectionCreator: "",
				CollectionId:      "1",
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail minting when collection does not exist",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          "",
				CollectionCreator: "",
				CollectionId:      "2",
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should successfully mint nfts for the same address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "1", Title: "nft1", Description: "nft1"},
						{Id: "2", Title: "nft2", Description: "nft2"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail mint nfts for the same address when strict flag is set to true and collection creator is empty",
			req: &types.MsgMintNfts{
				Creator:      s.addrs[0].String(),
				CollectionId: "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
					},
				},
				Strict: true,
			},
			cnt:    2,
			expErr: true,
			errMsg: "invalid collection creator",
		},
		{
			name: "should successfully mint nfts for another address",
			req: &types.MsgMintNfts{
				Creator:      s.addrs[0].String(),
				Receiver:     s.addrs[1].String(),
				CollectionId: "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
					},
				},
				Strict: false,
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint existing nfts",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "1", Title: "nft1", Description: "nft1"},
						{Id: "2", Title: "nft2", Description: "nft2"},
					},
				},
			},
			expErr: true,
			errMsg: "nfts count provided is invalid",
		},
		{
			name: "should successfully mint only the nfts that do not exist",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
						{Id: "4", Title: "nft4", Description: "nft4"},
					},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nfts with no permission",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		},
		{
			name: "should fail when mint nfts with did set to true",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
				Did: true,
			},
			expErr: true,
			errMsg: "cannot use did",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "1",
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftsRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nfts for the same address",
			req: &types.MsgMintNfts{
				Creator:           s.testAdminAccount,
				CollectionId:      "1",
				CollectionCreator: s.testAdminAccount,
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "1", Title: "nft1", Description: "nft1"},
						{Id: "2", Title: "nft2", Description: "nft2"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nfts for another address",
			req: &types.MsgMintNfts{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
						{Id: "4", Title: "nft4", Description: "nft4"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nfts with no permission",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "should fail when mint nfts with did set to true",
			req: &types.MsgMintNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
				Did: true,
			},
			cnt:    1,
			expErr: true,
			errMsg: "cannot use did",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "1",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftsDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nfts in the default collection for the same address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "1", Title: "nft1", Description: "nft1"},
						{Id: "2", Title: "nft2", Description: "nft2"},
					},
				},
				Strict: false,
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully mint nfts in the default collection for another address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
						{Id: "4", Title: "nft4", Description: "nft4"},
					},
				},
				Strict: false,
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully mint nfts in the default collection when it is set explicitly",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "default",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
				Strict: true,
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nfts in the default collection with strict flag set to true",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "6", Title: "nft6", Description: "nft6"},
					},
				},
				Strict: true,
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.CollectionId, "default")
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftsSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nfts for the same address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionId:      "2",
				CollectionCreator: s.addrs[0].String(),
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "1", Title: "nft1", Description: "nft1"},
						{Id: "2", Title: "nft2", Description: "nft2"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nfts for another address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "2",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
						{Id: "4", Title: "nft4", Description: "nft4"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nfts with no permission",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "2",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "should fail when mint with did set to true",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "2",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
				Did: true,
			},
			expErr: true,
			errMsg: "cannot use did",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "2",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftsOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nfts for the same address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionId:      "3",
				CollectionCreator: s.addrs[0].String(),
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "1", Title: "nft1", Description: "nft1"},
						{Id: "2", Title: "nft2", Description: "nft2"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nfts for another address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "3",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
						{Id: "4", Title: "nft4", Description: "nft4"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully mint nfts with another address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "3",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		}, {
			name: "should fail when mint with did set to true",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "3",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "6", Title: "nft6", Description: "nft6"},
					},
				},
				Did: true,
			},
			expErr: true,
			errMsg: "cannot use did",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:     "3",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftsRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nfts for the same address",
			req: &types.MsgMintNfts{
				Creator:           s.testAdminAccount,
				CollectionId:      "2",
				CollectionCreator: s.testAdminAccount,
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "1", Title: "nft1", Description: "nft1"},
						{Id: "2", Title: "nft2", Description: "nft2"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nfts for another address",
			req: &types.MsgMintNfts{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "2",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
						{Id: "4", Title: "nft4", Description: "nft4"},
					},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nfts with no permission",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "2",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "successfully mint nfts with did set to true",
			req: &types.MsgMintNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "2",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
					},
				},
				Did: true,
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "2",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNfts() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgBurnNfts{
				Creator:           "",
				CollectionCreator: "",
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1", "2"},
				},
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail burning when collection does not exist",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: "",
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1", "2"},
				},
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should fail burning nfts when not an owner",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1", "2"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner: nfts count provided is invalid",
		},
		{
			name: "should successfully burn only the nfts that do exist",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3", "5"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nfts",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1", "2"},
				},
			},
			cnt:    2,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "4",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "4",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.BurnNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftsRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nfts for the same address",
			req: &types.MsgBurnNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "3",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nfts for another address",
			req: &types.MsgBurnNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "3",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nfts with no permission",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "3",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"4"},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "should successfully burn only existing nfts",
			req: &types.MsgBurnNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "3",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2", "5"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "3",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "3",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "3",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.BurnNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftsDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nfts in the default collection for the same address",
			req: &types.MsgBurnNfts{
				Creator: s.addrs[0].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"7"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nfts for another address",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"9"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		}, {
			name: "should successfully burn nfts in the default collection when it is set explicitly",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "default",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"8"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully burn nfts in the default collection from another address",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"10"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "7", Title: "nft7", Description: "nft7"},
				{Id: "8", Title: "nft8", Description: "nft8"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "9", Title: "nft9", Description: "nft9"},
				{Id: "10", Title: "nft10", Description: "nft10"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.BurnNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftsSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nfts for the same address",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nfts with another address",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nfts with no permission",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"4"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		}, {
			name: "should successfully burn only existing nfts",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2", "5"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "5",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "5",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "5",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.BurnNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftsOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nfts for the same address",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "6",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nfts with another address",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "6",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nfts with no permission",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "6",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"4"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		}, {
			name: "should successfully burn only existing nfts",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "6",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2", "5"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:     "6",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "6",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "6",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.BurnNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftsRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nfts for the same address",
			req: &types.MsgBurnNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nfts for another address",
			req: &types.MsgBurnNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nfts with no permission",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"4"},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "should successfully burn only existing nfts",
			req: &types.MsgBurnNfts{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "4",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2", "5"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "4",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "4",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "4",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.BurnNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNfts() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgTransferNfts{
				Creator:           "",
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: "",
				CollectionId:      "7",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail transferring nft when collection does not exist",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "8",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should fail transferring nft when no transfer permission/not an owner",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "7",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or no transfer permission",
		},
		{
			name: "should successfully transfer nft",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "7",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "7",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "7",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.TransferNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftsRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully transfer nfts from the same address",
			req: &types.MsgTransferNfts{
				Creator:           s.testAdminAccount,
				Owner:             s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail transferring nfts from another address",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		},
		{
			name: "should successfully transfer nfts for another address",
			req: &types.MsgTransferNfts{
				Creator:           s.testAdminAccount,
				Owner:             s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully transfer only existing nfts",
			req: &types.MsgTransferNfts{
				Creator:           s.testAdminAccount,
				Owner:             s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "5",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2", "5"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "5",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "5",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "5",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.TransferNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftsDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully transfer nfts in the default collection for the same address",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"11"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		}, {
			name: "should fail transferring nfts for another address",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"12"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or no transfer permission",
		}, {
			name: "should successfully transfer only the existing nfts",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"12", "13"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[1].String(),
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "11", Title: "nft11", Description: "nft11"},
				{Id: "12", Title: "nft12", Description: "nft12"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.TransferNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftsSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should fail transferring nfts from the same address",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "8",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail transferring nfts from another address",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "8",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "8",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "8",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.TransferNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftsOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should fail transferring nft when no transfer permission/not an owner",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "9",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or no transfer permission",
		},
		{
			name: "should successfully transfer nft",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "9",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:     "9",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "9",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.TransferNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftsRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should fail transferring nfts from the same address",
			req: &types.MsgTransferNfts{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "6",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail transferring nfts from another address",
			req: &types.MsgTransferNfts{
				Creator:           s.testAdminAccount,
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "6",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "6",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "6",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.TransferNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNfts() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgApproveNfts{
				Creator:           "",
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "10",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail approving nft when collection does not exist",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "10",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should fail approving nft when not an owner",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "10",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		},
		{
			name: "should successfully approve nft",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "10",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "10",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "10",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.ApproveNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftsRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully approve nfts from the same address",
			req: &types.MsgApproveNfts{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "7",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail approving nfts from another address",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "7",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		},
		{
			name: "should successfully approve nfts for another address",
			req: &types.MsgApproveNfts{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "7",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"3"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully approve only existing nfts",
			req: &types.MsgApproveNfts{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "7",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2", "5"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "7",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "7",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "7",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.ApproveNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftsDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully approve nfts in the default collection for the same address",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"13"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		}, {
			name: "should fail approving nfts for another address",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"14"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		}, {
			name: "should successfully approve only the existing nfts",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"14", "15"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[1].String(),
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "13", Title: "nft13", Description: "nft13"},
				{Id: "14", Title: "nft14", Description: "nft14"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.ApproveNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftsSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should fail approving nfts from the same address",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "11",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail approving nfts from another address",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "11",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "11",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "11",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.ApproveNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftsOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should fail approving nft when not an owner",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "12",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		},
		{
			name: "should successfully approve nft",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "12",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:     "12",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "12",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.ApproveNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftsRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNfts
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should fail approving nfts from the same address",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "8",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail approving nfts from another address",
			req: &types.MsgApproveNfts{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "8",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"2"},
				},
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "8",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "8",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.ApproveNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
				s.Require().EqualValues(req.NftsCount, tc.cnt)
			}
		})
	}
}
