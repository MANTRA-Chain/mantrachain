package keeper_test

import (
	"testing"

	utils "github.com/MANTRA-Finance/aumega/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/suite"

	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/MANTRA-Finance/aumega/app"
	"github.com/MANTRA-Finance/aumega/testutil"
	"github.com/MANTRA-Finance/aumega/x/token/keeper"
	"github.com/MANTRA-Finance/aumega/x/token/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app              *chain.App
	ctx              sdk.Context
	keeper           keeper.Keeper
	querier          keeper.Querier
	msgServer        types.MsgServer
	addrs            []sdk.AccAddress
	testAdminAccount string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = testutil.SetupWithGenesisValSet(s.T())
	hdr := cbproto.Header{
		Height: s.app.LastBlockHeight() + 1,
		Time:   utils.ParseTime("2022-01-01T00:00:00Z"),
	}
	s.addrs = testutil.CreateIncrementalAccounts(3)
	s.app.BeginBlock(abci.RequestBeginBlock{Header: hdr})
	s.ctx = s.app.BaseApp.NewContext(false, hdr)
	s.app.BeginBlocker(s.ctx, abci.RequestBeginBlock{Header: hdr})
	s.keeper = s.app.TokenKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.testAdminAccount = "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"
}
