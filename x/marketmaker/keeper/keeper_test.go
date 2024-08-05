package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	utils "github.com/MANTRA-Finance/mantrachain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/keeper"
	marketmaker "github.com/MANTRA-Finance/mantrachain/x/marketmaker/module"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

const (
	denom1 = "denom1"
	denom2 = "denom2"
	denom3 = "denom3"
)

var (
	initialBalances = sdk.NewCoins(
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 1_000_000_000_000),
		sdk.NewInt64Coin(denom1, 1_000_000_000),
		sdk.NewInt64Coin(denom2, 1_000_000_000),
		sdk.NewInt64Coin(denom3, 1_000_000_000))
)

type KeeperTestSuite struct {
	suite.Suite

	app         *chain.App
	ctx         sdk.Context
	keeper      keeper.Keeper
	QueryHelper *baseapp.QueryServiceTestHelper
	queryClient types.QueryClient
	msgServer   types.MsgServer
	govHandler  govv1beta1.Handler
	addrs       []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	app, err := testutil.Setup()
	suite.NoError(err)

	suite.app = app
	suite.ctx = app.BaseApp.NewContextLegacy(false, cmtproto.Header{
		Height: 1,
		Time:   utils.ParseTime("2022-01-01T00:00:05Z"),
	})
	suite.keeper = suite.app.MarketmakerKeeper
	suite.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: suite.app.GRPCQueryRouter(),
		Ctx:             suite.ctx,
	}

	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)

	suite.govHandler = marketmaker.NewMarketMakerProposalHandler(suite.keeper)
	suite.addrs = testutil.AddTestAddrsAndAdmin(suite.app, suite.ctx, 30, sdkmath.ZeroInt())
	for _, addr := range suite.addrs {
		err := testutil.FundAccount(suite.app.BankKeeper, suite.ctx, addr, initialBalances)
		suite.Require().NoError(err)
	}
	suite.SetIncentivePairs()
}

func (suite *KeeperTestSuite) AddTestAddrs(num int, coins sdk.Coins) []sdk.AccAddress {
	addrs := testutil.AddTestAddrsAndAdmin(suite.app, suite.ctx, num, sdkmath.ZeroInt())
	for _, addr := range addrs {
		err := testutil.FundAccount(suite.app.BankKeeper, suite.ctx, addr, coins)
		suite.Require().NoError(err)
	}
	return addrs
}

func (suite *KeeperTestSuite) SetIncentivePairs() {
	params := suite.keeper.GetParams(suite.ctx)
	params.IncentivePairs = []types.IncentivePair{
		{
			PairId: uint64(1),
		},
		{
			PairId: uint64(2),
		},
		{
			PairId: uint64(3),
		},
		{
			PairId: uint64(4),
		},
		{
			PairId: uint64(5),
		},
		{
			PairId: uint64(6),
		},
		{
			PairId: uint64(7),
		},
	}
	suite.keeper.SetParams(suite.ctx, params)
}

func (suite *KeeperTestSuite) ResetIncentivePairs() {
	params := suite.keeper.GetParams(suite.ctx)
	params.IncentivePairs = []types.IncentivePair{}
	suite.keeper.SetParams(suite.ctx, params)
}

func (suite *KeeperTestSuite) handleProposal(content govv1beta1.Content) {
	suite.T().Helper()
	err := content.ValidateBasic()
	suite.Require().NoError(err)
	err = suite.govHandler(suite.ctx, content)
	suite.Require().NoError(err)
}
