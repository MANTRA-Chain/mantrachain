package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/AumegaChain/aumega/x/coinfactory/types"
)

func (suite *KeeperTestSuite) TestAdminMsgs() {
	addr0bal := int64(0)
	addr1bal := int64(0)

	bankKeeper := suite.app.BankKeeper

	suite.CreateDefaultDenom()
	// Make sure that the admin is set correctly
	queryRes, err := suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: suite.defaultDenom,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(suite.addrs[0].String(), queryRes.AuthorityMetadata.Admin)

	// Test minting to own accounts
	_, err = suite.msgServer.Mint(sdk.WrapSDKContext(suite.ctx), types.NewMsgMint(suite.addrs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, 10)))
	addr0bal += 10
	suite.Require().NoError(err)
	suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom).Amount.Int64() == addr0bal, bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom))

	// Test force transferring
	_, err = suite.msgServer.ForceTransfer(sdk.WrapSDKContext(suite.ctx), types.NewMsgForceTransfer(suite.addrs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, 5), suite.addrs[0].String(), suite.addrs[1].String()))
	addr0bal -= 5
	addr1bal += 5
	suite.Require().NoError(err)
	suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom).IsEqual(sdk.NewInt64Coin(suite.defaultDenom, addr0bal)))
	suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[1], suite.defaultDenom).IsEqual(sdk.NewInt64Coin(suite.defaultDenom, addr1bal)))

	_, err = suite.msgServer.ForceTransfer(sdk.WrapSDKContext(suite.ctx), types.NewMsgForceTransfer(suite.addrs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, 5), suite.addrs[1].String(), suite.addrs[0].String()))
	addr0bal += 5
	addr1bal -= 5
	suite.Require().NoError(err)
	suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom).IsEqual(sdk.NewInt64Coin(suite.defaultDenom, addr0bal)))
	suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[1], suite.defaultDenom).IsEqual(sdk.NewInt64Coin(suite.defaultDenom, addr1bal)))

	// Test burning from own account
	_, err = suite.msgServer.Burn(sdk.WrapSDKContext(suite.ctx), types.NewMsgBurn(suite.addrs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, 5)))
	addr0bal -= 5
	suite.Require().NoError(err)
	suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[1], suite.defaultDenom).Amount.Int64() == addr1bal)

	// Test Change Admin
	_, err = suite.msgServer.ChangeAdmin(sdk.WrapSDKContext(suite.ctx), types.NewMsgChangeAdmin(suite.addrs[0].String(), suite.defaultDenom, suite.addrs[1].String()))
	suite.Require().NoError(err)
	queryRes, err = suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: suite.defaultDenom,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(suite.addrs[1].String(), queryRes.AuthorityMetadata.Admin)

	// Make sure old admin can no longer do actions
	_, err = suite.msgServer.Burn(sdk.WrapSDKContext(suite.ctx), types.NewMsgBurn(suite.addrs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, 5)))
	suite.Require().Error(err)

	// Make sure the new admin works
	_, err = suite.msgServer.Mint(sdk.WrapSDKContext(suite.ctx), types.NewMsgMint(suite.addrs[1].String(), sdk.NewInt64Coin(suite.defaultDenom, 5)))
	addr1bal += 5
	suite.Require().NoError(err)
	suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[1], suite.defaultDenom).Amount.Int64() == addr1bal)

	// Try setting admin to empty
	_, err = suite.msgServer.ChangeAdmin(sdk.WrapSDKContext(suite.ctx), types.NewMsgChangeAdmin(suite.addrs[1].String(), suite.defaultDenom, ""))
	suite.Require().NoError(err)
	queryRes, err = suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: suite.defaultDenom,
	})
	suite.Require().NoError(err)
	suite.Require().Equal("", queryRes.AuthorityMetadata.Admin)
}

// TestMintDenom ensures the following properties of the MintMessage:
// * Noone can mint tokens for a denom that doesn't exist
// * Only the admin of a denom can mint tokens for it
// * The admin of a denom can mint tokens for it
func (suite *KeeperTestSuite) TestMintDenom() {
	var addr0bal int64

	// Create a denom
	suite.CreateDefaultDenom()

	for _, tc := range []struct {
		desc      string
		amount    int64
		mintDenom string
		admin     string
		valid     bool
	}{
		{
			desc:      "denom does not exist",
			amount:    10,
			mintDenom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/testcoin",
			admin:     suite.addrs[0].String(),
			valid:     false,
		},
		{
			desc:      "mint is not by the admin",
			amount:    10,
			mintDenom: suite.defaultDenom,
			admin:     suite.addrs[1].String(),
			valid:     false,
		},
		{
			desc:      "success case",
			amount:    10,
			mintDenom: suite.defaultDenom,
			admin:     suite.addrs[0].String(),
			valid:     true,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// Test minting to admins own account
			bankKeeper := suite.app.BankKeeper
			_, err := suite.msgServer.Mint(sdk.WrapSDKContext(suite.ctx), types.NewMsgMint(tc.admin, sdk.NewInt64Coin(tc.mintDenom, 10)))
			if tc.valid {
				addr0bal += 10
				suite.Require().NoError(err)
				suite.Require().Equal(bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom).Amount.Int64(), addr0bal, bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom))
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBurnDenom() {
	var addr0bal int64

	// Create a denom.
	suite.CreateDefaultDenom()

	// mint 10 default token for testAcc[0]
	suite.msgServer.Mint(sdk.WrapSDKContext(suite.ctx), types.NewMsgMint(suite.addrs[0].String(), sdk.NewInt64Coin(suite.defaultDenom, 10)))
	addr0bal += 10

	for _, tc := range []struct {
		desc      string
		amount    int64
		burnDenom string
		admin     string
		valid     bool
	}{
		{
			desc:      "denom does not exist",
			amount:    10,
			burnDenom: "factory/cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw/testcoin",
			admin:     suite.addrs[0].String(),
			valid:     false,
		},
		{
			desc:      "burn is not by the admin",
			amount:    10,
			burnDenom: suite.defaultDenom,
			admin:     suite.addrs[1].String(),
			valid:     false,
		},
		{
			desc:      "burn amount is bigger than minted amount",
			amount:    1000,
			burnDenom: suite.defaultDenom,
			admin:     suite.addrs[1].String(),
			valid:     false,
		},
		{
			desc:      "success case",
			amount:    10,
			burnDenom: suite.defaultDenom,
			admin:     suite.addrs[0].String(),
			valid:     true,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// Test minting to admins own account
			bankKeeper := suite.app.BankKeeper
			_, err := suite.msgServer.Burn(sdk.WrapSDKContext(suite.ctx), types.NewMsgBurn(tc.admin, sdk.NewInt64Coin(tc.burnDenom, 10)))
			if tc.valid {
				addr0bal -= 10
				suite.Require().NoError(err)
				suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom).Amount.Int64() == addr0bal, bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom))
			} else {
				suite.Require().Error(err)
				suite.Require().True(bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom).Amount.Int64() == addr0bal, bankKeeper.GetBalance(suite.ctx, suite.addrs[0], suite.defaultDenom))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestChangeAdminDenom() {
	for _, tc := range []struct {
		desc                    string
		msgChangeAdmin          func(denom string) *types.MsgChangeAdmin
		expectedChangeAdminPass bool
		expectedAdminIndex      int
		msgMint                 func(denom string) *types.MsgMint
		expectedMintPass        bool
	}{
		{
			desc: "creator admin can't mint after setting to '' ",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(suite.addrs[0].String(), denom, "")
			},
			expectedChangeAdminPass: true,
			expectedAdminIndex:      -1,
			msgMint: func(denom string) *types.MsgMint {
				return types.NewMsgMint(suite.addrs[0].String(), sdk.NewInt64Coin(denom, 5))
			},
			expectedMintPass: false,
		},
		{
			desc: "non-admins can't change the existing admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(suite.addrs[1].String(), denom, suite.addrs[2].String())
			},
			expectedChangeAdminPass: false,
			expectedAdminIndex:      0,
		},
		{
			desc: "success change admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(suite.addrs[0].String(), denom, suite.addrs[1].String())
			},
			expectedAdminIndex:      1,
			expectedChangeAdminPass: true,
			msgMint: func(denom string) *types.MsgMint {
				return types.NewMsgMint(suite.addrs[1].String(), sdk.NewInt64Coin(denom, 5))
			},
			expectedMintPass: true,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			suite.SetupTest()

			// Create a denom and mint
			res, err := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[0].String(), "bitcoin"))
			suite.Require().NoError(err)

			testDenom := res.GetNewTokenDenom()

			_, err = suite.msgServer.Mint(sdk.WrapSDKContext(suite.ctx), types.NewMsgMint(suite.addrs[0].String(), sdk.NewInt64Coin(testDenom, 10)))
			suite.Require().NoError(err)

			_, err = suite.msgServer.ChangeAdmin(sdk.WrapSDKContext(suite.ctx), tc.msgChangeAdmin(testDenom))
			if tc.expectedChangeAdminPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

			queryRes, err := suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
				Denom: testDenom,
			})
			suite.Require().NoError(err)

			// expectedAdminIndex with negative value is assumed as admin with value of ""
			const emptyStringAdminIndexFlag = -1
			if tc.expectedAdminIndex == emptyStringAdminIndexFlag {
				suite.Require().Equal("", queryRes.AuthorityMetadata.Admin)
			} else {
				suite.Require().Equal(suite.addrs[tc.expectedAdminIndex].String(), queryRes.AuthorityMetadata.Admin)
			}

			// we test mint to test if admin authority is performed properly after admin change.
			if tc.msgMint != nil {
				_, err := suite.msgServer.Mint(sdk.WrapSDKContext(suite.ctx), tc.msgMint(testDenom))
				if tc.expectedMintPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSetDenomMetaData() {
	// setup test
	suite.SetupTest()
	suite.CreateDefaultDenom()

	for _, tc := range []struct {
		desc                string
		msgSetDenomMetadata types.MsgSetDenomMetadata
		expectedPass        bool
	}{
		{
			desc: "successful set denom metadata",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(suite.addrs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    suite.defaultDenom,
						Exponent: 0,
					},
					{
						Denom:    "stake",
						Exponent: 6,
					},
				},
				Base:    suite.defaultDenom,
				Display: "stake",
				Name:    "AUM",
				Symbol:  "AUM",
			}),
			expectedPass: true,
		},
		{
			desc: "non existent factory denom name",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(suite.addrs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    fmt.Sprintf("factory/%s/litecoin", suite.addrs[0].String()),
						Exponent: 0,
					},
					{
						Denom:    "stake",
						Exponent: 6,
					},
				},
				Base:    fmt.Sprintf("factory/%s/litecoin", suite.addrs[0].String()),
				Display: "stake",
				Name:    "AUM",
				Symbol:  "AUM",
			}),
			expectedPass: false,
		},
		{
			desc: "non-factory denom",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(suite.addrs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    "stake",
						Exponent: 0,
					},
					{
						Denom:    "stakeo",
						Exponent: 6,
					},
				},
				Base:    "stake",
				Display: "stakeo",
				Name:    "AUM",
				Symbol:  "AUM",
			}),
			expectedPass: false,
		},
		{
			desc: "wrong admin",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(suite.addrs[1].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    suite.defaultDenom,
						Exponent: 0,
					},
					{
						Denom:    "stake",
						Exponent: 6,
					},
				},
				Base:    suite.defaultDenom,
				Display: "stake",
				Name:    "AUM",
				Symbol:  "AUM",
			}),
			expectedPass: false,
		},
		{
			desc: "invalid metadata (missing display denom unit)",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(suite.addrs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    suite.defaultDenom,
						Exponent: 0,
					},
				},
				Base:    suite.defaultDenom,
				Display: "stake",
				Name:    "AUM",
				Symbol:  "AUM",
			}),
			expectedPass: false,
		},
	} {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			bankKeeper := suite.app.BankKeeper
			res, err := suite.msgServer.SetDenomMetadata(sdk.WrapSDKContext(suite.ctx), &tc.msgSetDenomMetadata)
			if tc.expectedPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)

				md, found := bankKeeper.GetDenomMetaData(suite.ctx, suite.defaultDenom)
				suite.Require().True(found)
				suite.Require().Equal(tc.msgSetDenomMetadata.Metadata.Name, md.Name)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
