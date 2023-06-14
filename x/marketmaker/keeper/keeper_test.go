package keeper_test

import (
	"testing"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"

	"github.com/MANTRA-Finance/mantrachain/x/marketmaker"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/marketmaker/types"
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
	testutil.IBCConnectionTestSuite

	app        *chain.App
	ctx        sdk.Context
	keeper     keeper.Keeper
	querier    keeper.Querier
	msgServer  types.MsgServer
	govHandler govtypes.Handler
	addrs      []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.IBCConnectionTestSuite.SetupTest()
	suite.app = suite.GetApp(suite.Chain)
	hdr := tmproto.Header{}
	suite.ctx = suite.app.BaseApp.NewContext(false, hdr)
	suite.keeper = suite.app.MarketMakerKeeper
	suite.querier = keeper.Querier{Keeper: suite.keeper}
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)
	suite.govHandler = marketmaker.NewMarketMakerProposalHandler(suite.keeper)
	suite.addrs = testutil.CreateRandomAccounts(30)
	for _, addr := range suite.addrs {
		suite.FundAAccount(addr, initialBalances)
	}
	suite.SetIncentivePairs()
}

func (suite *KeeperTestSuite) FundAAccount(addr sdk.AccAddress, amt sdk.Coins) {
	suite.T().Helper()
	suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, liquiditytypes.ModuleName, amt))
	suite.Require().NoError(
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, liquiditytypes.ModuleName, addr, amt))
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

func (suite *KeeperTestSuite) handleProposal(content govtypes.Content) {
	suite.T().Helper()
	err := content.ValidateBasic()
	suite.Require().NoError(err)
	err = suite.govHandler(suite.ctx, content)
	suite.Require().NoError(err)
}
