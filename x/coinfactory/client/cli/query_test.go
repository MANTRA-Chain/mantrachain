package cli_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	abci "github.com/cometbft/cometbft/abci/types"
	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
)

type QueryTestSuite struct {
	suite.Suite

	app         *app.App
	addrs       []sdk.AccAddress
	ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	queryClient types.QueryClient
}

var (
	SecondaryDenom  = "ucoin"
	SecondaryAmount = math.NewInt(100000000)
)

func (s *QueryTestSuite) SetupSuite() {
	s.app = testutil.Setup(false)
	hdr := cbproto.Header{
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
	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(SecondaryDenom, SecondaryAmount))
	testutil.FundAccount(s.app.BankKeeper, s.ctx, s.addrs[0], fundAccsAmount)
	// create new token
	_, err := s.app.CoinFactoryKeeper.CreateDenom(s.ctx, s.addrs[0].String(), "coinfactory")
	s.Require().NoError(err)

	oldHeight := s.ctx.BlockHeight()
	oldHeader := s.ctx.BlockHeader()
	s.app.Commit()
	newHeader := cbproto.Header{Height: oldHeight + 1, ChainID: oldHeader.ChainID, Time: oldHeader.Time.Add(time.Second)}
	s.app.BeginBlock(abci.RequestBeginBlock{Header: newHeader})
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
