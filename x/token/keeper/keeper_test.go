package keeper_test

import (
	"context"
	"encoding/binary"
	"testing"
	"time"

	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/token/keeper"
	tokentestutil "github.com/MANTRA-Finance/mantrachain/x/token/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
)

type KeeperTestSuite struct {
	testutil.IBCConnectionTestSuite

	app       *app.App
	ctx       sdk.Context
	keeper    keeper.Keeper
	querier   keeper.Querier
	msgServer types.MsgServer
	addrs     []sdk.AccAddress
	nftKeeper *tokentestutil.MockNFTKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.IBCConnectionTestSuite.SetupTest()
	s.app = s.GetApp(s.Chain)
	hdr := tmproto.Header{
		Height: s.app.LastBlockHeight() + 1,
		Time:   utils.ParseTime("2022-01-01T00:00:00Z"),
	}
	ctrl := gomock.NewController(s.T())
	s.addrs = testutil.CreateIncrementalAccounts(3)
	s.app.BeginBlock(abci.RequestBeginBlock{Header: hdr})
	s.ctx = s.app.BaseApp.NewContext(false, hdr)
	s.app.BeginBlocker(s.ctx, abci.RequestBeginBlock{Header: hdr})
	s.keeper = s.app.TokenKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.nftKeeper = tokentestutil.NewMockNFTKeeper(ctrl)

	guardKeeper := tokentestutil.NewMockGuardKeeper(ctrl)
	guardKeeper.EXPECT().CheckIsAdmin(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	guardKeeper.EXPECT().CheckRestrictedNftsCollection(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	guardKeeper.EXPECT().CheckNewRestrictedNftsCollection(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
}

// Below are just shortcuts to frequently-used functions.
func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
	return s.app.BankKeeper.GetAllBalances(s.ctx, addr)
}

func (s *KeeperTestSuite) getBalance(addr sdk.AccAddress, denom string) sdk.Coin {
	return s.app.BankKeeper.GetBalance(s.ctx, addr, denom)
}

func (s *KeeperTestSuite) sendCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.SendCoins(s.ctx, fromAddr, toAddr, amt)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) nextBlock() {
	s.T().Helper()
	s.app.EndBlock(abci.RequestEndBlock{})
	s.app.Commit()
	hdr := tmproto.Header{
		Height: s.app.LastBlockHeight() + 1,
		Time:   s.ctx.BlockTime().Add(5 * time.Second),
	}
	s.app.BeginBlock(abci.RequestBeginBlock{Header: hdr})
	s.ctx = s.app.BaseApp.NewContext(false, hdr)
	s.app.BeginBlocker(s.ctx, abci.RequestBeginBlock{Header: hdr})
}

// Below are useful helpers to write test code easily.
func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) mintNft(creator sdk.AccAddress) {
	s.T().Helper()
	msgMetadata := &types.MsgCreateNftCollectionMetadata{Id: "1"}
	msg := types.NewMsgCreateNftCollection(creator.String(), msgMetadata)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.msgServer.CreateNftCollection(context.Background(), msg)
	s.Require().NoError(err)
}
