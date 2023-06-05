package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

type KeeperTestSuite struct {
	testutil.IBCConnectionTestSuite

	app         *app.App
	ctx         sdk.Context
	keeper      keeper.Keeper
	QueryHelper *baseapp.QueryServiceTestHelper
	queryClient types.QueryClient
	msgServer   types.MsgServer
	// defaultDenom is on the suite, as it depends on the creator test address.
	defaultDenom string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.IBCConnectionTestSuite.SetupTest()
	suite.app = suite.GetApp(suite.Chain)
	hdr := tmproto.Header{
		Height: suite.app.LastBlockHeight() + 1,
		Time:   utils.ParseTime("2022-01-01T00:00:00Z"),
	}
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: hdr})
	suite.ctx = suite.app.BaseApp.NewContext(false, hdr)
	suite.app.BeginBlocker(suite.ctx, abci.RequestBeginBlock{Header: hdr})
	suite.keeper = suite.app.CoinFactoryKeeper
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(testutil.SecondaryDenom, testutil.SecondaryAmount))
	for _, acc := range suite.TestAccs {
		testutil.FundAccount(suite.app.BankKeeper, suite.ctx, acc, fundAccsAmount)
	}
	suite.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: suite.app.GRPCQueryRouter(),
		Ctx:             suite.ctx,
	}
	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)
}

func (suite *KeeperTestSuite) CreateDefaultDenom() {
	res, _ := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.TestAccs[0].String(), "bitcoin"))
	suite.defaultDenom = res.GetNewTokenDenom()
}

func (suite *KeeperTestSuite) TestCreateModuleAccount() {
	app := suite.app

	// remove module account
	coinfactoryModuleAccount := app.AccountKeeper.GetAccount(suite.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	app.AccountKeeper.RemoveAccount(suite.ctx, coinfactoryModuleAccount)

	// ensure module account was removed
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	coinfactoryModuleAccount = app.AccountKeeper.GetAccount(suite.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().Nil(coinfactoryModuleAccount)

	// create module account
	app.CoinFactoryKeeper.CreateModuleAccount(suite.ctx)

	// check that the module account is now initialized
	coinfactoryModuleAccount = app.AccountKeeper.GetAccount(suite.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().NotNil(coinfactoryModuleAccount)
}
