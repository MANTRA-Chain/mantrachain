package liquidfarming_test

import (
	"encoding/binary"

	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	cbproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	"github.com/MANTRA-Finance/mantrachain/x/liquidfarming/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/liquidfarming/types"
	liquiditytypes "github.com/MANTRA-Finance/mantrachain/x/liquidity/types"

	_ "github.com/stretchr/testify/suite"
)

type ModuleTestSuite struct {
	suite.Suite

	app       *chain.App
	ctx       sdk.Context
	keeper    keeper.Keeper
	querier   keeper.Querier
	msgServer types.MsgServer
}

func TestModuleTestSuite(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}

func (s *ModuleTestSuite) SetupTest() {
	s.app = testutil.SetupWithGenesisValSet(s.T())
	s.ctx = s.app.BaseApp.NewContext(false, cbproto.Header{})
	s.keeper = s.app.LiquidFarmingKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.ctx = s.ctx.WithBlockTime(time.Now()) // set to current time
}

//
// Below are just shortcuts to frequently-used functions.
//

func (s *ModuleTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) { //nolint
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt)
	s.Require().NoError(err)
}

func (s *ModuleTestSuite) createPair(creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string) liquiditytypes.Pair {
	s.T().Helper()
	s.fundAddr(creator, s.app.LiquidityKeeper.GetPairCreationFee(s.ctx))
	pair, err := s.app.LiquidityKeeper.CreatePair(s.ctx, liquiditytypes.NewMsgCreatePair(creator, baseCoinDenom, quoteCoinDenom))
	s.Require().NoError(err)
	return pair
}

// createPairWithLastPrice is a convenient method to create a pair with last price.
// it is needed as farming plan doesn't distribute farming rewards if the last price is not set.
func (s *ModuleTestSuite) createPairWithLastPrice(creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string, lastPrice sdk.Dec) liquiditytypes.Pair {
	s.T().Helper()
	pair := s.createPair(creator, baseCoinDenom, quoteCoinDenom)
	pair.LastPrice = &lastPrice
	s.app.LiquidityKeeper.SetPair(s.ctx, pair)
	return pair
}

func (s *ModuleTestSuite) createPool(creator sdk.AccAddress, pairId uint64, depositCoins sdk.Coins) liquiditytypes.Pool {
	s.T().Helper()
	s.fundAddr(creator, s.app.LiquidityKeeper.GetPoolCreationFee(s.ctx).Add(depositCoins...))
	pool, err := s.app.LiquidityKeeper.CreatePool(s.ctx, liquiditytypes.NewMsgCreatePool(creator, pairId, depositCoins))
	s.Require().NoError(err)
	return pool
}

func (s *ModuleTestSuite) createLiquidFarm(poolId uint64, minFarmAmt, minBidAmt math.Int, feeRate sdk.Dec) types.LiquidFarm { //nolint
	s.T().Helper()
	liquidFarm := types.NewLiquidFarm(poolId, minFarmAmt, minBidAmt, feeRate)
	params := s.keeper.GetParams(s.ctx)
	params.LiquidFarms = append(params.LiquidFarms, liquidFarm)
	s.keeper.SetParams(s.ctx, params)
	s.keeper.SetLiquidFarm(s.ctx, liquidFarm)
	return liquidFarm
}

//
// Below are helper functions to write test code easily
//

func (s *ModuleTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}
