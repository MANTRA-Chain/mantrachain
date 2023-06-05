package cli_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

type QueryTestSuite struct {
	testutil.IBCConnectionTestSuite

	app         *app.App
	ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	queryClient types.QueryClient
}

func (s *QueryTestSuite) SetupSuite() {
	s.IBCConnectionTestSuite.SetupTest()
	s.app = s.GetApp(s.Chain)
	hdr := tmproto.Header{
		Height: s.app.LastBlockHeight() + 1,
		Time:   utils.ParseTime("2022-01-01T00:00:00Z"),
	}
	s.app.BeginBlock(abci.RequestBeginBlock{Header: hdr})
	s.ctx = s.app.BaseApp.NewContext(false, hdr)
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.app.GRPCQueryRouter(),
		Ctx:             s.ctx,
	}
	s.queryClient = types.NewQueryClient(s.QueryHelper)

	// fund acc
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(testutil.SecondaryDenom, testutil.SecondaryAmount))
	testutil.FundAccount(s.app.BankKeeper, s.ctx, s.TestAccs[0], fundAccsAmount)
	// create new token
	_, err := s.app.CoinFactoryKeeper.CreateDenom(s.ctx, s.TestAccs[0].String(), "coinfactory")
	s.Require().NoError(err)

	oldHeight := s.ctx.BlockHeight()
	oldHeader := s.ctx.BlockHeader()
	s.app.Commit()
	newHeader := tmproto.Header{Height: oldHeight + 1, ChainID: oldHeader.ChainID, Time: oldHeader.Time.Add(time.Second)}
	s.app.BeginBlock(abci.RequestBeginBlock{Header: newHeader})
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
