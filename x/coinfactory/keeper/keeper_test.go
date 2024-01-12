package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"

	"github.com/AumegaChain/aumega/app"
	"github.com/AumegaChain/aumega/testutil"
	utils "github.com/AumegaChain/aumega/types"
	"github.com/AumegaChain/aumega/x/coinfactory/keeper"
	"github.com/AumegaChain/aumega/x/coinfactory/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.App
	addrs       []sdk.AccAddress
	ctx         sdk.Context
	keeper      keeper.Keeper
	QueryHelper *baseapp.QueryServiceTestHelper
	queryClient types.QueryClient
	msgServer   types.MsgServer
	// defaultDenom is on the suite, as it depends on the creator test address.
	defaultDenom string
}

var (
	SecondaryDenom  = "ucoin"
	SecondaryAmount = math.NewInt(100000000)
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = testutil.SetupWithGenesisValSet(suite.T())
	hdr := cbproto.Header{
		Height: suite.app.LastBlockHeight() + 1,
		Time:   utils.ParseTime("2022-01-01T00:00:00Z"),
	}
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: hdr})
	suite.ctx = suite.app.BaseApp.NewContext(false, hdr)
	suite.app.BeginBlocker(suite.ctx, abci.RequestBeginBlock{Header: hdr})
	suite.keeper = suite.app.CoinFactoryKeeper
	initialBalances := sdk.NewCoins(sdk.NewCoin(SecondaryDenom, SecondaryAmount))
	suite.addrs = testutil.AddTestAddrsAndAdmin(suite.app, suite.ctx, 6, sdk.ZeroInt())
	for _, addr := range suite.addrs {
		err := testutil.FundAccount(suite.app.BankKeeper, suite.ctx, addr, initialBalances)
		suite.Require().NoError(err)
	}
	suite.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: suite.app.GRPCQueryRouter(),
		Ctx:             suite.ctx,
	}
	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)
}

func (suite *KeeperTestSuite) CreateDefaultDenom() {
	res, err := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(suite.addrs[0].String(), "bitcoin"))
	require.NoError(suite.T(), err)
	suite.defaultDenom = res.GetNewTokenDenom()
}

func (suite *KeeperTestSuite) TestCreateModuleAccount() {
	app := suite.app

	// remove module account
	coinfactoryModuleAccount := app.AccountKeeper.GetAccount(suite.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	app.AccountKeeper.RemoveAccount(suite.ctx, coinfactoryModuleAccount)

	// ensure module account was removed
	suite.ctx = app.BaseApp.NewContext(false, cbproto.Header{})
	coinfactoryModuleAccount = app.AccountKeeper.GetAccount(suite.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().Nil(coinfactoryModuleAccount)

	// create module account
	app.CoinFactoryKeeper.CreateModuleAccount(suite.ctx)

	// check that the module account is now initialized
	coinfactoryModuleAccount = app.AccountKeeper.GetAccount(suite.ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	suite.Require().NotNil(coinfactoryModuleAccount)
}

// AssertEventEmitted asserts that ctx's event manager has emitted the given number of events
// of the given type.
func (s *KeeperTestSuite) AssertEventEmitted(ctx sdk.Context, eventTypeExpected string, numEventsExpected int) {
	allEvents := ctx.EventManager().Events()
	// filter out other events
	actualEvents := make([]sdk.Event, 0)
	for _, event := range allEvents {
		if event.Type == eventTypeExpected {
			actualEvents = append(actualEvents, event)
		}
	}
	s.Equal(numEventsExpected, len(actualEvents))
}

func (s *KeeperTestSuite) FindEvent(events []sdk.Event, name string) sdk.Event {
	index := slices.IndexFunc(events, func(e sdk.Event) bool { return e.Type == name })
	if index == -1 {
		return sdk.Event{}
	}
	return events[index]
}

func (s *KeeperTestSuite) ExtractAttributes(event sdk.Event) map[string]string {
	attrs := make(map[string]string)
	if event.Attributes == nil {
		return attrs
	}
	for _, a := range event.Attributes {
		attrs[string(a.Key)] = string(a.Value)
	}
	return attrs
}
