package keeper_test

import (
	"testing"

	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/MANTRA-Finance/aumega/app"
	"github.com/MANTRA-Finance/aumega/testutil"
	"github.com/MANTRA-Finance/aumega/x/marketmaker"
	"github.com/MANTRA-Finance/aumega/x/marketmaker/keeper"
	"github.com/MANTRA-Finance/aumega/x/marketmaker/types"
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

	app        *chain.App
	ctx        sdk.Context
	keeper     keeper.Keeper
	querier    keeper.Querier
	msgServer  types.MsgServer
	govHandler govv1beta1.Handler
	addrs      []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	app := testutil.SetupWithGenesisValSet(suite.T())
	ctx := app.BaseApp.NewContext(false, cbproto.Header{})

	suite.app = app
	suite.ctx = ctx
	suite.keeper = suite.app.MarketMakerKeeper
	suite.querier = keeper.Querier{Keeper: suite.keeper}
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)
	suite.govHandler = marketmaker.NewMarketMakerProposalHandler(suite.keeper)
	suite.addrs = testutil.AddTestAddrsAndAdmin(suite.app, suite.ctx, 30, sdk.ZeroInt())
	for _, addr := range suite.addrs {
		err := testutil.FundAccount(suite.app.BankKeeper, suite.ctx, addr, initialBalances)
		suite.Require().NoError(err)
	}
	suite.SetIncentivePairs()
}

func (suite *KeeperTestSuite) AddTestAddrs(num int, coins sdk.Coins) []sdk.AccAddress {
	addrs := testutil.AddTestAddrsAndAdmin(suite.app, suite.ctx, num, sdk.ZeroInt())
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
