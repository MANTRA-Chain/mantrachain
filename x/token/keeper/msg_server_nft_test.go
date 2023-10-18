package keeper_test

import (
	"github.com/MANTRA-Finance/mantrachain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestMintNftsNotRestrictedCollection() {
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
			errMsg: "not found: 2: nft collection does not exists",
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
			name: "should successfully mint nfts for another address",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
					},
				},
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
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "1",
		},
	})
	if err != nil {
		return
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftsRestrictedCollection() {
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
		}, {
			name: "should successfully mint nfts for another address",
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
				CollectionId:      "1",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "3", Title: "nft3", Description: "nft3"},
					},
				},
			},
			expErr: true,
			errMsg: "unauthorized",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id:             "2",
			RestrictedNfts: true,
		},
	})
	if err != nil {
		return
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := s.msgServer.MintNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestMintNftsDefaultCollection() {
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
				CollectionId:      "",
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
			name: "should fail when mint nfts in the default collection with strict flag set to true",
			req: &types.MsgMintNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "",
				Nfts: &types.MsgNftsMetadata{
					Nfts: []*types.MsgNftMetadata{
						{Id: "5", Title: "nft5", Description: "nft5"},
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
				if tc.req.Receiver == "" {
					s.Require().EqualValues(req.Receiver, tc.req.Creator)
				} else {
					s.Require().EqualValues(req.Receiver, tc.req.Receiver)
				}
				s.Require().EqualValues(req.NftsCount, tc.cnt)
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestBurnNfts() {
	testCases := []struct {
		name   string
		req    *types.MsgBurnNfts
		expErr bool
		errMsg string
	}{
		{
			name: "empty address error",
			req: &types.MsgBurnNfts{
				Creator:           "",
				CollectionCreator: "",
				CollectionId:      "1",
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
				CollectionId:      "2",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1", "2"},
				},
			},
			expErr: true,
			errMsg: "not found: 2: nft collection does not exists",
		},
		{
			name: "should fail burning nfts when not an owner",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[0].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1", "2"},
				},
			},
			expErr: true,
			errMsg: "not existing nfts or not an owner: nfts count provided is invalid",
		},
		{
			name: "should successfully burn nfts",
			req: &types.MsgBurnNfts{
				Creator:           s.addrs[1].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				Nfts: &types.MsgNftsIds{
					NftsIds: []string{"1", "2"},
				},
			},
			expErr: false,
			errMsg: "",
		},
	}

	goCtx := sdk.WrapSDKContext(s.ctx)

	_, err := s.msgServer.CreateNftCollection(goCtx, &types.MsgCreateNftCollection{
		Creator: s.addrs[0].String(),
		Collection: &types.MsgCreateNftCollectionMetadata{
			Id: "1",
		},
	})
	if err != nil {
		return
	}

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "1",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{{Id: "1", Title: "nft1", Description: "nft1"}, {Id: "2", Title: "nft1", Description: "nft2"}},
		},
	})

	s.Require().NoError(err)

	//s.nftKeeper.Mint(goCtx,  &types.Nft{Id: "1"}, s.addrs[0])

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.msgServer.BurnNfts(goCtx, tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
func (s *KeeperTestSuite) TestTransferNfts() {
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
				CollectionId:      "1",
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
				CollectionId:      "2",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "not found: 2: nft collection does not exists",
		},
		{
			name: "should fail transferring nft when no transfer permission/not an owner",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[0].String(),
				Owner:             s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "not existing nft or no transfer permission: nfts provided is invalid",
		},
		{
			name: "should successfully transfer nft",
			req: &types.MsgTransferNft{
				Creator:           s.addrs[1].String(),
				Owner:             s.addrs[1].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
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
			Id: "1",
		},
	})
	if err != nil {
		return
	}

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "1",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{{Id: "1", Title: "nft1", Description: "nft1"}, {Id: "2", Title: "nft1", Description: "nft2"}},
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
				Receiver:          s.addrs[2].String(),
				CollectionCreator: "",
				CollectionId:      "1",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "empty address string is not allowed",
		},
		{
			name: "should fail approving nft when collection does not exist",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[1].String(),
				CollectionId:      "2",
				NftId:             "1",
			},
			expErr: true,
			errMsg: "not found: 2: nft collection does not exists",
		},
		{
			name: "should fail transferring nft when no transfer permission/not an owner",
			req: &types.MsgApproveNft{
				Creator:           s.addrs[0].String(),
				Receiver:          s.addrs[2].String(),
				CollectionCreator: s.addrs[0].String(),
				CollectionId:      "1",
				NftId:             "1",
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
				CollectionId:      "1",
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
			Id: "1",
		},
	})
	if err != nil {
		return
	}

	_, err = s.msgServer.MintNfts(goCtx, &types.MsgMintNfts{
		Creator:           s.addrs[0].String(),
		Receiver:          s.addrs[1].String(),
		CollectionCreator: s.addrs[0].String(),
		CollectionId:      "1",
		Nfts: &types.MsgNftsMetadata{
			Nfts: []*types.MsgNftMetadata{{Id: "1", Title: "nft1", Description: "nft1"}, {Id: "2", Title: "nft1", Description: "nft2"}},
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
