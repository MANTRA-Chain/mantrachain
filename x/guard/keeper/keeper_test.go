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
	"github.com/MANTRA-Finance/mantrachain/x/guard/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
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
	defaultPrivileges []byte
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
	suite.keeper = suite.app.GuardKeeper
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(testutil.SecondaryDenom, testutil.SecondaryAmount))
	for _, acc := range suite.TestAccs {
		testutil.FundAccount(suite.app.BankKeeper, suite.ctx, acc, fundAccsAmount)
	}
	suite.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: suite.app.GRPCQueryRouter(),
		Ctx:             suite.ctx,
	}
	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(&suite.keeper)
	suite.defaultPrivileges = []byte{0x01}
}
