package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

func (suite *KeeperTestSuite) TestMsgCreateDenom() {
	// Creating a denom should work
	res, err := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[0].String(), "bitcoin"))
	suite.Require().NoError(err)
	suite.Require().NotEmpty(res.GetNewTokenDenom())

	// Make sure that the admin is set correctly
	queryRes, err := suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: res.GetNewTokenDenom(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(suite.addrs[0].String(), queryRes.AuthorityMetadata.Admin)

	// Make sure that a second version of the same denom can't be recreated
	res, err = suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[0].String(), "bitcoin"))
	suite.Require().Error(err)

	// Creating a second denom should work
	res, err = suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[0].String(), "litecoin"))
	suite.Require().NoError(err)
	suite.Require().NotEmpty(res.GetNewTokenDenom())

	// Try querying all the denoms created by suite.addrs[0]
	queryRes2, err := suite.queryClient.DenomsFromCreator(suite.ctx.Context(), &types.QueryDenomsFromCreatorRequest{
		Creator: suite.addrs[0].String(),
	})
	suite.Require().NoError(err)
	suite.Require().Len(queryRes2.Denoms, 2)

	// Make sure that cannot create a denom with account that is not the admin
	_, err = suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[1].String(), "bitcoin"))
	suite.Require().Error(err)

	// Make sure that an address with a "/" in it can't create denoms
	_, err = suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom("mantrachain.eth/creator", "bitcoin"))
	suite.Require().Error(err)
}

func (suite *KeeperTestSuite) TestCreateDenom() {
	for _, tc := range []struct {
		desc             string
		denomCreationFee types.Params
		setup            func()
		subdenom         string
		valid            bool
	}{
		{
			desc:     "subdenom too long",
			subdenom: "assadsadsadasdasdsadsadsadsadsadsadsklkadaskkkdasdasedskhanhassyeunganassfnlksdflksafjlkasd",
			valid:    false,
		},
		{
			desc: "subdenom and creator pair already exists",
			setup: func() {
				_, err := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[0].String(), "bitcoin"))
				suite.Require().NoError(err)
			},
			subdenom: "bitcoin",
			valid:    false,
		},
		{
			desc:     "subdenom having invalid characters",
			subdenom: "bit/***///&&&/coin",
			valid:    false,
		},
	} {
		suite.SetupTest()
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			if tc.setup != nil {
				tc.setup()
			}
			coinFactoryKeeper := suite.app.CoinFactoryKeeper
			bankKeeper := suite.app.BankKeeper
			// Set denom creation fee in params
			coinFactoryKeeper.SetParams(suite.ctx, tc.denomCreationFee)

			// note balance, create a coinfactory denom, then note balance again
			preCreateBalance := bankKeeper.GetAllBalances(suite.ctx, suite.addrs[0])
			res, err := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[0].String(), tc.subdenom))
			postCreateBalance := bankKeeper.GetAllBalances(suite.ctx, suite.addrs[0])
			if tc.valid {
				suite.Require().NoError(err)

				// Make sure that the admin is set correctly
				queryRes, err := suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
					Denom: res.GetNewTokenDenom(),
				})

				suite.Require().NoError(err)
				suite.Require().Equal(suite.addrs[0].String(), queryRes.AuthorityMetadata.Admin)

			} else {
				suite.Require().Error(err)
				// Ensure we don't charge if we expect an error
				suite.Require().True(preCreateBalance.IsEqual(postCreateBalance))
			}
		})
	}
}
