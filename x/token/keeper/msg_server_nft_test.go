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
			name: "should fail when mint nfts with did set to true",
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
			name: "should fail when mint nfts with did set to true",
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
			name: "should fail burning nfts when collection does not exist",
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
			errMsg: "not existing nfts or not an owner",
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
				Creator:           s.addrs[1].String(),
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
			name: "should fail transferring nfts when collection does not exist",
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
			name: "should fail transferring nfts when no transfer permission/not an owner",
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
			name: "should successfully transfer nfts",
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
			name: "should fail transferring nfts when no transfer permission/not an owner",
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
			name: "should successfully transfer nfts",
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
				Approved: true,
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
				Approved: true,
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should fail approving nfts when not an owner",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "10",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
				Approved: true,
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		},
		{
			name: "should successfully approve nfts",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "10",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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
			name: "should fail approving nfts when not an owner",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "12",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
				Approved: true,
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner",
		},
		{
			name: "should successfully approve nfts",
			req: &types.MsgApproveNfts{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "12",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1"},
				},
				Approved: true,
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
				Approved: true,
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
				Approved: true,
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

func (s *KeeperTestSuite) TestApproveTransferNfts() {
	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "13",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
				{Id: "5", Title: "nft5", Description: "nft5"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "6", Title: "nft6", Description: "nft6"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"1"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	req1, err := s.msgServer.ApproveNfts(goCtx, &types.MsgApproveNfts{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"1"},
		},
		Approved: true,
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req1.NftsCount, 1)

	req2, err := s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"1"},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req2.NftsCount, 1)

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"1"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	req3, err := s.msgServer.ApproveNfts(goCtx, &types.MsgApproveNfts{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"2"},
		},
		Approved: true,
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req3.NftsCount, 1)

	req4, err := s.msgServer.ApproveNfts(goCtx, &types.MsgApproveNfts{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"2"},
		},
		Approved: false,
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req4.NftsCount, 1)

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"2"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	req5, err := s.msgServer.ApproveNfts(goCtx, &types.MsgApproveNfts{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"2"},
		},
		Approved: true,
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req5.NftsCount, 1)

	req6, err := s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"2", "6"},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req6.NftsCount, 1)

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[1].String(),
		Receiver: s.addrs[0].String(),
		Approved: true,
	})
	s.Require().NoError(err)

	req7, err := s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"3", "4", "6"},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req7.NftsCount, 2)

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[1].String(),
		Receiver: s.addrs[0].String(),
		Approved: false,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"5"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[2].String(),
		Receiver: s.addrs[0].String(),
		Approved: true,
	})
	s.Require().NoError(err)

	req8, err := s.msgServer.ApproveNfts(goCtx, &types.MsgApproveNfts{
		Creator:           s.addrs[2].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6", "7"},
		},
		Approved: true,
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req8.NftsCount, 1)

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[1].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	req9, err := s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req9.NftsCount, 1)

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	req10, err := s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[1].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req10.NftsCount, 1)

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[2].String(),
		Receiver: s.addrs[0].String(),
		Approved: false,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nfts or no transfer permission")

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[1].String(),
		Receiver: s.addrs[0].String(),
		Approved: true,
	})
	s.Require().NoError(err)

	req11, err := s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[2].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req11.NftsCount, 1)

	req12, err := s.msgServer.TransferNfts(goCtx, &types.MsgTransferNfts{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "13",
		Nfts: &types.MsgNftsIds{
			NftsIds: []string{"6"},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(req12.NftsCount, 1)
}

func (s *KeeperTestSuite) TestMintNft() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNft
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgMintNft{
				Creator:           "",
				Receiver:          "",
				CollectionCreator: "",
				CollectionId:      "14",
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail minting when collection does not exist",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				Receiver:          "",
				CollectionCreator: "",
				CollectionId:      "15",
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should successfully mint nft for the same address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "14",
				Nft: &types.MsgNftMetadata{
					Id: "1", Title: "nft1", Description: "nft1",
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail mint nft for the same address when strict flag is set to true and collection creator is empty",
			req: &types.MsgMintNft{
				Creator:      s.addrs[0].String(),
				CollectionId: "14",
				Nft: &types.MsgNftMetadata{
					Id: "2", Title: "nft2", Description: "nft2",
				},
				Strict: true,
			},
			expErr: true,
			errMsg: "invalid collection creator",
		},
		{
			name: "should successfully mint nft for another address",
			req: &types.MsgMintNft{
				Creator:      s.addrs[0].String(),
				Receiver:     s.addrs[1].String(),
				CollectionId: "14",
				Nft: &types.MsgNftMetadata{
					Id: "2", Title: "nft2", Description: "nft2",
				},
				Strict: false,
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint existing nft",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "14",
				Nft: &types.MsgNftMetadata{
					Id: "2", Title: "nft2", Description: "nft2",
				},
			},
			expErr: true,
			errMsg: "existing or invalid nft",
		},
		{
			name: "should fail when mint nfts with no permission",
			req: &types.MsgMintNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "14",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		},
		{
			name: "should fail when mint nft with did set to true",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "14",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
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
			Id: "14",
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNft(goCtx, tc.req)
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
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nft for the same address",
			req: &types.MsgMintNft{
				Creator:           s.testAdminAccount,
				CollectionId:      "9",
				CollectionCreator: s.testAdminAccount,
				Nft: &types.MsgNftMetadata{
					Id: "1", Title: "nft1", Description: "nft1",
				},
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nft for another address",
			req: &types.MsgMintNft{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "9",
				Nft: &types.MsgNftMetadata{
					Id: "2", Title: "nft2", Description: "nft2",
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nft with no permission",
			req: &types.MsgMintNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "9",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "should fail when mint nft with did set to true",
			req: &types.MsgMintNft{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "9",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
				},
				Did: true,
			},
			expErr: true,
			errMsg: "cannot use did",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "9",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNft(goCtx, tc.req)
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
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nft in the default collection for the same address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[1].String(),
				Nft: &types.MsgNftMetadata{
					Id: "15", Title: "nft15", Description: "nft15",
				},
				Strict: false,
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully mint nft in the default collection for another address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[1].String(),
				Nft: &types.MsgNftMetadata{
					Id: "16", Title: "nft16", Description: "nft16",
				},
				Strict: false,
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully mint nft in the default collection when it is set explicitly",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "default",
				Nft: &types.MsgNftMetadata{
					Id: "17", Title: "nft17", Description: "nft17",
				},
				Strict: true,
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nft in the default collection with strict flag set to true",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "",
				Nft: &types.MsgNftMetadata{
					Id: "18", Title: "nft18", Description: "nft18",
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
			req, err := s.msgServer.MintNft(goCtx, tc.req)
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
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nft for the same address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionId:      "15",
				CollectionCreator: s.addrs[0].String(),
				Nft: &types.MsgNftMetadata{
					Id: "1", Title: "nft1", Description: "nft1",
				},
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nft for another address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "15",
				Nft: &types.MsgNftMetadata{
					Id: "2", Title: "nft2", Description: "nft2",
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nft with no permission",
			req: &types.MsgMintNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "15",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "should fail when mint nft with did set to true",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "15",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
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
			Id:             "15",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNft(goCtx, tc.req)
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
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nft for the same address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				CollectionId:      "16",
				CollectionCreator: s.addrs[0].String(),
				Nft: &types.MsgNftMetadata{
					Id: "1", Title: "nft1", Description: "nft1",
				},
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nft for another address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "16",
				Nft: &types.MsgNftMetadata{
					Id: "2", Title: "nft2", Description: "nft2",
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully mint nft with another address",
			req: &types.MsgMintNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "16",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
				},
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should fail when mint nft with did set to true",
			req: &types.MsgMintNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "16",
				Nft: &types.MsgNftMetadata{
					Id: "4", Title: "nft4", Description: "nft4",
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
			Id:     "16",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNft(goCtx, tc.req)
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
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgMintNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully mint nft for the same address",
			req: &types.MsgMintNft{
				Creator:           s.testAdminAccount,
				CollectionId:      "10",
				CollectionCreator: s.testAdminAccount,
				Nft: &types.MsgNftMetadata{
					Id: "1", Title: "nft1", Description: "nft1",
				},
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully mint nft for another address",
			req: &types.MsgMintNft{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "10",
				Nft: &types.MsgNftMetadata{
					Id: "2", Title: "nft2", Description: "nft2",
				},
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail when mint nft with no permission",
			req: &types.MsgMintNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "10",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		}, {
			name: "successfully mint nft with did set to true",
			req: &types.MsgMintNft{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "10",
				Nft: &types.MsgNftMetadata{
					Id: "3", Title: "nft3", Description: "nft3",
				},
				Did: true,
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "10",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNft(goCtx, tc.req)
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
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNft() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNft
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgBurnNft{
				Creator:           "",
				CollectionCreator: "",
				CollectionId:      "17",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail burning nft when collection does not exist",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: "",
				CollectionId:      "18",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should fail burning nft when not an owner",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "17",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "not existing nft or not an owner",
		},
		{
			name: "should successfully burn nft",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "17",
				NftId:             "1",
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "17",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "17",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.BurnNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nft for the same address",
			req: &types.MsgBurnNft{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "11",
				NftId:             "1",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nft for another address",
			req: &types.MsgBurnNft{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "11",
				NftId:             "3",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nft with no permission",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "11",
				NftId:             "2",
			},
			expErr: true,
			errMsg: "unauthorized",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "11",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "11",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "3", Title: "nft3", Description: "nft3"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "11",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "2", Title: "nft2", Description: "nft2"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.BurnNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nft in the default collection for the same address",
			req: &types.MsgBurnNft{
				Creator: s.addrs[0].String(),
				NftId:   "19",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nft for another address",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				NftId:             "21",
			},
			expErr: true,
			errMsg: "not existing nft or not an owner",
		}, {
			name: "should successfully burn nft in the default collection when it is set explicitly",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "default",
				NftId:             "20",
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should successfully burn nft in the default collection from another address",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				NftId:             "21",
			},
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
				{Id: "19", Title: "nft19", Description: "nft19"},
				{Id: "20", Title: "nft20", Description: "nft20"},
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
				{Id: "21", Title: "nft21", Description: "nft21"},
				{Id: "22", Title: "nft22", Description: "nft22"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.BurnNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNft
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nft for the same address",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "18",
				NftId:             "1",
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nft with another address",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "18",
				NftId:             "3",
			},
			cnt:    1,
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nft with no permission",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "18",
				NftId:             "4",
			},
			expErr: true,
			errMsg: "not existing nft or not an owner",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "18",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "18",
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
		CollectionId:      "18",
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
			_, err := s.msgServer.BurnNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nft for the same address",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "19",
				NftId:             "1",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nft with another address",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "19",
				NftId:             "3",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nft with no permission",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "19",
				NftId:             "4",
			},
			expErr: true,
			errMsg: "not existing nft or not an owner",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:     "19",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "19",
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
		CollectionId:      "19",
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
			_, err := s.msgServer.BurnNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNftRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully burn nft for the same address",
			req: &types.MsgBurnNft{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "12",
				NftId:             "1",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should successfully burn nft for another address",
			req: &types.MsgBurnNft{
				Creator:           s.testAdminAccount,
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "12",
				NftId:             "3",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail burning nft with no permission",
			req: &types.MsgBurnNft{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "12",
				NftId:             "4",
			},
			expErr: true,
			errMsg: "unauthorized",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "12",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "12",
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
		CollectionId:      "12",
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
			_, err := s.msgServer.BurnNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNft() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNft
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgTransferNft{
				Creator:           "",
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: "",
				CollectionId:      "20",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail transferring nft when collection does not exist",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "20",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should fail transferring nft when no transfer permission/not an owner",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "20",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "not existing nft or no transfer permission",
		},
		{
			name: "should successfully transfer nft",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "20",
				NftId:             "1",
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "20",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "20",
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
			_, err := s.msgServer.TransferNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully transfer nft from the same address",
			req: &types.MsgTransferNft{
				Creator:           s.testAdminAccount,
				Owner:             s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "13",
				NftId:             "1",
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail transferring nft from another address",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "13",
				NftId:             "3",
			},
			expErr: true,
			errMsg: "unauthorized",
		},
		{
			name: "should successfully transfer nft for another address",
			req: &types.MsgTransferNft{
				Creator:           s.testAdminAccount,
				Owner:             s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "13",
				NftId:             "3",
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "13",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "13",
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
		CollectionId:      "13",
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
			_, err := s.msgServer.TransferNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully transfer nft in the default collection for the same address",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				NftId:             "23",
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should fail transferring nft for another address",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				NftId:             "24",
			},
			expErr: true,
			errMsg: "not existing nft or no transfer permission",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[1].String(),
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "23", Title: "nft23", Description: "nft23"},
				{Id: "24", Title: "nft24", Description: "nft24"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.TransferNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNft
		expErr bool
		errMsg string
	}{
		{
			name: "should fail transferring nft from the same address",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "21",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail transferring nft from another address",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "21",
				NftId:             "2",
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "21",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "21",
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
			_, err := s.msgServer.TransferNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNft
		expErr bool
		errMsg string
	}{
		{
			name: "should fail transferring nft when no transfer permission/not an owner",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "22",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "not existing nft or no transfer permission",
		},
		{
			name: "should successfully transfer nft",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "22",
				NftId:             "1",
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:     "22",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "22",
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
			_, err := s.msgServer.TransferNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestTransferNftRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgTransferNft
		expErr bool
		errMsg string
	}{
		{
			name: "should fail transferring nft from the same address",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "14",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail transferring nft from another address",
			req: &types.MsgTransferNft{
				Creator:           s.testAdminAccount,
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "14",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "14",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "14",
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
			_, err := s.msgServer.TransferNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNft() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNft
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgApproveNft{
				Creator:           "",
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "23",
				NftId:             "1",
				Approved:          true,
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail approving nft when collection does not exist",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "23",
				NftId:             "1",
				Approved:          true,
			},
			expErr: true,
			errMsg: "nft collection does not exists",
		},
		{
			name: "should fail approving nft when not an owner",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "23",
				NftId:             "1",
				Approved:          true,
			},
			expErr: true,
			errMsg: "not existing nft or not an owner",
		},
		{
			name: "should successfully approve nft",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "23",
				NftId:             "1",
				Approved:          true,
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "23",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "23",
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
			_, err := s.msgServer.ApproveNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftRestrictedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully approve nft from the same address",
			req: &types.MsgApproveNft{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: "",
				CollectionId:      "15",
				NftId:             "1",
				Approved:          true,
			},
			expErr: false,
			errMsg: "",
		},
		{
			name: "should fail approving nft from another address",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "15",
				NftId:             "3",
				Approved:          true,
			},
			expErr: true,
			errMsg: "unauthorized",
		},
		{
			name: "should successfully approve nft for another address",
			req: &types.MsgApproveNft{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "15",
				NftId:             "3",
				Approved:          true,
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "15",
			RestrictedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "15",
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
		CollectionId:      "15",
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
			_, err := s.msgServer.ApproveNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftDefaultNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNft
		expErr bool
		errMsg string
	}{
		{
			name: "should successfully approve nft in the default collection for the same address",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				NftId:             "25",
				Approved:          true,
			},
			expErr: false,
			errMsg: "",
		}, {
			name: "should fail approving nft for another address",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				NftId:             "26",
				Approved:          true,
			},
			expErr: true,
			errMsg: "not existing nft or not an owner",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[1].String(),
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "25", Title: "nft25", Description: "nft25"},
				{Id: "26", Title: "nft26", Description: "nft26"},
			},
		},
	})
	s.Require().NoError(err)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.ApproveNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNft
		expErr bool
		errMsg string
	}{
		{
			name: "should fail approving nft from the same address",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "24",
				NftId:             "1",
				Approved:          true,
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail approving nft from another address",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "24",
				NftId:             "2",
				Approved:          true,
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "24",
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "24",
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
			_, err := s.msgServer.ApproveNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftOpenedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNft
		cnt    int
		expErr bool
		errMsg string
	}{
		{
			name: "should fail approving nft when not an owner",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "25",
				NftId:             "1",
				Approved:          true,
			},
			expErr: true,
			errMsg: "not existing nft or not an owner",
		},
		{
			name: "should successfully approve nft",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "25",
				NftId:             "1",
				Approved:          true,
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:     "25",
			Opened: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "25",
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
			_, err := s.msgServer.ApproveNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveNftRestrictedAndSoulBondedNftCollection() {
	testCases := []struct {
		name   string
		req    *types.MsgApproveNft
		expErr bool
		errMsg string
	}{
		{
			name: "should fail approving nft from the same address",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "16",
				NftId:             "1",
				Approved:          true,
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		}, {
			name: "should fail approving nft from another address",
			req: &types.MsgApproveNft{
				Creator:           s.testAdminAccount,
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.testAdminAccount,
				CollectionId:      "16",
				NftId:             "2",
				Approved:          true,
			},
			expErr: true,
			errMsg: "soul bonded nft collection operation disabled",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.testAdminAccount,
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "16",
			RestrictedNfts: true,
			SoulBondedNfts: true,
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.testAdminAccount,
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.testAdminAccount,
		CollectionId:      "16",
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
			_, err := s.msgServer.ApproveNft(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestApproveTransferNft() {
	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "26",
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "1", Title: "nft1", Description: "nft1"},
				{Id: "2", Title: "nft2", Description: "nft2"},
				{Id: "3", Title: "nft3", Description: "nft3"},
				{Id: "4", Title: "nft4", Description: "nft4"},
				{Id: "5", Title: "nft5", Description: "nft5"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{
				{Id: "6", Title: "nft6", Description: "nft6"},
			},
		},
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "1",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.ApproveNft(goCtx, &types.MsgApproveNft{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "1",
		Approved:          true,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "1",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "1",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.ApproveNft(goCtx, &types.MsgApproveNft{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "2",
		Approved:          true,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.ApproveNft(goCtx, &types.MsgApproveNft{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "2",
		Approved:          false,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "2",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.ApproveNft(goCtx, &types.MsgApproveNft{
		Creator:           s.addrs[1].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "2",
		Approved:          true,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "2",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[1].String(),
		Receiver: s.addrs[0].String(),
		Approved: true,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "3",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "4",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[1].String(),
		Receiver: s.addrs[0].String(),
		Approved: false,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "5",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[2].String(),
		Receiver: s.addrs[0].String(),
		Approved: true,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.ApproveNft(goCtx, &types.MsgApproveNft{
		Creator:           s.addrs[2].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
		Approved:          true,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.ApproveNft(goCtx, &types.MsgApproveNft{
		Creator:           s.addrs[2].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "7",
		Approved:          true,
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or not an owner")

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[1].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[0].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[1].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[2].String(),
		Receiver: s.addrs[0].String(),
		Approved: false,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "not existing nft or no transfer permission")

	_, err = s.msgServer.ApproveAllNfts(goCtx, &types.MsgApproveAllNfts{
		Creator:  s.addrs[1].String(),
		Receiver: s.addrs[0].String(),
		Approved: true,
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[2].String(),
		Owner:             s.addrs[2].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().NoError(err)

	_, err = s.msgServer.TransferNft(goCtx, &types.MsgTransferNft{
		Creator:           s.addrs[0].String(),
		Owner:             s.addrs[1].String(),
		Receiver:          s.addrs[2].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "26",
		NftId:             "6",
	})
	s.Require().NoError(err)
}
