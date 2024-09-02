package keeper_test

import (
	"encoding/binary"
	"testing"
	"time"

	"cosmossdk.io/math"
	chain "github.com/MANTRA-Finance/mantrachain/app"
	"github.com/MANTRA-Finance/mantrachain/testutil"
	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/amm"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *chain.App
	ctx         sdk.Context
	keeper      *keeper.Keeper
	QueryHelper *baseapp.QueryServiceTestHelper
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

const expectedStr = "expected:\t%v\ngot:\t\t%v"

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app, err := testutil.Setup()
	s.NoError(err)

	s.app = app
	s.ctx = app.BaseApp.NewContextLegacy(false, cmtproto.Header{
		Height: 1,
		Time:   utils.ParseTime("2022-01-01T00:00:05Z"),
	})
	s.keeper = s.app.LiquidityKeeper
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.app.GRPCQueryRouter(),
		Ctx:             s.ctx,
	}

	s.queryClient = types.NewQueryClient(s.QueryHelper)
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
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
	ctx, err := testutil.NextBlock(s.app, s.ctx, time.Second*5)
	s.Require().NoError(err)
	s.ctx = ctx
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.app.GRPCQueryRouter(),
		Ctx:             s.ctx,
	}
	s.queryClient = types.NewQueryClient(s.QueryHelper)
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

func (s *KeeperTestSuite) createPair(creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string) types.Pair {
	s.T().Helper()
	s.fundAddr(creator, s.keeper.GetPairCreationFee(s.ctx))
	msg := types.NewMsgCreatePair(creator, baseCoinDenom, quoteCoinDenom, &math.LegacyDec{}, &math.LegacyDec{})
	s.Require().NoError(msg.ValidateBasic())
	pair, err := s.keeper.CreatePair(s.ctx, msg)
	s.Require().NoError(err)
	return pair
}

func (s *KeeperTestSuite) createPool(creator sdk.AccAddress, pairId uint64, depositCoins sdk.Coins) types.Pool {
	s.T().Helper()
	s.fundAddr(creator, depositCoins.Add(s.keeper.GetPoolCreationFee(s.ctx)...))
	msg := types.NewMsgCreatePool(creator, pairId, depositCoins)
	s.Require().NoError(msg.ValidateBasic())
	pool, err := s.keeper.CreatePool(s.ctx, msg)
	s.Require().NoError(err)
	return pool
}

func (s *KeeperTestSuite) createRangedPool(creator sdk.AccAddress, pairId uint64, depositCoins sdk.Coins, minPrice, maxPrice, initialPrice math.LegacyDec) types.Pool {
	s.T().Helper()
	s.fundAddr(creator, depositCoins.Add(s.keeper.GetPoolCreationFee(s.ctx)...))
	msg := types.NewMsgCreateRangedPool(creator, pairId, depositCoins, minPrice, maxPrice, initialPrice)
	s.Require().NoError(msg.ValidateBasic())
	pool, err := s.keeper.CreateRangedPool(s.ctx, msg)
	s.Require().NoError(err)
	return pool
}

func (s *KeeperTestSuite) deposit(depositor sdk.AccAddress, poolId uint64, depositCoins sdk.Coins, fund bool) types.DepositRequest {
	s.T().Helper()
	if fund {
		s.fundAddr(depositor, depositCoins)
	}
	req, err := s.keeper.Deposit(s.ctx, types.NewMsgDeposit(depositor, poolId, depositCoins))
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) withdraw(withdrawer sdk.AccAddress, poolId uint64, poolCoin sdk.Coin) types.WithdrawRequest {
	s.T().Helper()
	req, err := s.keeper.Withdraw(s.ctx, types.NewMsgWithdraw(withdrawer, poolId, poolCoin))
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) limitOrder(
	orderer sdk.AccAddress, pairId uint64, dir types.OrderDirection,
	price math.LegacyDec, amt math.Int, orderLifespan time.Duration, fund bool,
) types.Order {
	s.T().Helper()
	pair, found := s.keeper.GetPair(s.ctx, pairId)
	s.Require().True(found)
	var ammDir amm.OrderDirection
	var offerCoinDenom, demandCoinDenom string
	switch dir {
	case types.OrderDirectionBuy:
		ammDir = amm.Buy
		offerCoinDenom, demandCoinDenom = pair.QuoteCoinDenom, pair.BaseCoinDenom
	case types.OrderDirectionSell:
		ammDir = amm.Sell
		offerCoinDenom, demandCoinDenom = pair.BaseCoinDenom, pair.QuoteCoinDenom
	}
	offerCoin := sdk.NewCoin(offerCoinDenom, amm.OfferCoinAmount(ammDir, price, amt))
	if fund {
		s.fundAddr(orderer, sdk.NewCoins(offerCoin))
	}
	msg := types.NewMsgLimitOrder(
		orderer, pairId, dir, offerCoin, demandCoinDenom,
		price, amt, orderLifespan)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.LimitOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) buyLimitOrder(
	orderer sdk.AccAddress, pairId uint64, price math.LegacyDec,
	amt math.Int, orderLifespan time.Duration, fund bool,
) types.Order {
	s.T().Helper()
	return s.limitOrder(
		orderer, pairId, types.OrderDirectionBuy, price, amt, orderLifespan, fund)
}

func (s *KeeperTestSuite) sellLimitOrder(
	orderer sdk.AccAddress, pairId uint64, price math.LegacyDec,
	amt math.Int, orderLifespan time.Duration, fund bool,
) types.Order {
	s.T().Helper()
	return s.limitOrder(
		orderer, pairId, types.OrderDirectionSell, price, amt, orderLifespan, fund)
}

func (s *KeeperTestSuite) marketOrder(
	orderer sdk.AccAddress, pairId uint64, dir types.OrderDirection,
	amt math.Int, orderLifespan time.Duration, fund bool,
) types.Order {
	s.T().Helper()
	pair, found := s.keeper.GetPair(s.ctx, pairId)
	s.Require().True(found)
	s.Require().NotNil(pair.LastPrice)
	lastPrice := *pair.LastPrice
	var offerCoin sdk.Coin
	var demandCoinDenom string
	switch dir {
	case types.OrderDirectionBuy:
		maxPrice := lastPrice.Mul(math.LegacyOneDec().Add(s.keeper.GetMaxPriceLimitRatio(s.ctx)))
		offerCoin = sdk.NewCoin(pair.QuoteCoinDenom, amm.OfferCoinAmount(amm.Buy, maxPrice, amt))
		demandCoinDenom = pair.BaseCoinDenom
	case types.OrderDirectionSell:
		offerCoin = sdk.NewCoin(pair.BaseCoinDenom, amt)
		demandCoinDenom = pair.QuoteCoinDenom
	}
	if fund {
		s.fundAddr(orderer, sdk.NewCoins(offerCoin))
	}
	msg := types.NewMsgMarketOrder(
		orderer, pairId, dir, offerCoin, demandCoinDenom,
		amt, orderLifespan)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.MarketOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) buyMarketOrder(
	orderer sdk.AccAddress, pairId uint64,
	amt math.Int, orderLifespan time.Duration, fund bool,
) {
	s.T().Helper()
	s.marketOrder(
		orderer, pairId, types.OrderDirectionBuy, amt, orderLifespan, fund)
}

func (s *KeeperTestSuite) sellMarketOrder(
	orderer sdk.AccAddress, pairId uint64,
	amt math.Int, orderLifespan time.Duration, fund bool,
) {
	s.T().Helper()
	s.marketOrder(
		orderer, pairId, types.OrderDirectionSell, amt, orderLifespan, fund)
}

func (s *KeeperTestSuite) mmOrder(
	orderer sdk.AccAddress, pairId uint64,
	maxSellPrice, minSellPrice math.LegacyDec, sellAmt math.Int,
	maxBuyPrice, minBuyPrice math.LegacyDec, buyAmt math.Int,
	orderLifespan time.Duration, fund bool,
) []types.Order {
	s.T().Helper()
	if fund {
		pair, found := s.keeper.GetPair(s.ctx, pairId)
		s.Require().True(found)

		maxNumTicks := int(s.keeper.GetMaxNumMarketMakingOrderTicks(s.ctx))
		tickPrec := int(s.keeper.GetTickPrecision(s.ctx))

		var buyTicks, sellTicks []types.MMOrderTick
		offerBaseCoin := sdk.NewInt64Coin(pair.BaseCoinDenom, 0)
		offerQuoteCoin := sdk.NewInt64Coin(pair.QuoteCoinDenom, 0)
		if buyAmt.IsPositive() {
			buyTicks = types.MMOrderTicks(
				types.OrderDirectionBuy, minBuyPrice, maxBuyPrice, buyAmt, maxNumTicks, tickPrec)
			for _, tick := range buyTicks {
				offerQuoteCoin = offerQuoteCoin.AddAmount(tick.OfferCoinAmount)
			}
		}
		if sellAmt.IsPositive() {
			sellTicks = types.MMOrderTicks(
				types.OrderDirectionSell, minSellPrice, maxSellPrice, sellAmt, maxNumTicks, tickPrec)
			for _, tick := range sellTicks {
				offerBaseCoin = offerBaseCoin.AddAmount(tick.OfferCoinAmount)
			}
		}
		s.fundAddr(orderer, sdk.NewCoins(offerBaseCoin, offerQuoteCoin))
	}
	msg := types.NewMsgMMOrder(
		orderer, pairId,
		maxSellPrice, minSellPrice, sellAmt,
		maxBuyPrice, minBuyPrice, buyAmt,
		orderLifespan)
	s.Require().NoError(msg.ValidateBasic())
	orders, err := s.keeper.MMOrder(s.ctx, msg)
	s.Require().NoError(err)

	index, found := s.keeper.GetMMOrderIndex(s.ctx, orderer, pairId)
	maxNumTicks := int(s.keeper.GetMaxNumMarketMakingOrderTicks(s.ctx))
	s.Require().True(found)
	s.Require().Equal(orderer.String(), index.Orderer)
	s.Require().Equal(pairId, index.PairId)
	s.Require().True(len(index.OrderIds) <= maxNumTicks*2)
	s.Require().True(len(index.OrderIds) == len(orders))
	for i, order := range orders {
		s.Require().Equal(order.Id, index.OrderIds[i])
	}
	return orders
}

// nolint
func (s *KeeperTestSuite) cancelOrder(orderer sdk.AccAddress, pairId, orderId uint64) {
	s.T().Helper()
	err := s.keeper.CancelOrder(s.ctx, types.NewMsgCancelOrder(orderer, pairId, orderId))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) cancelAllOrders(orderer sdk.AccAddress, pairIds []uint64) {
	s.T().Helper()
	err := s.keeper.CancelAllOrders(s.ctx, types.NewMsgCancelAllOrders(orderer, pairIds))
	s.Require().NoError(err)
}

func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	if exp.IsEqual(got) {
		return true, "", "", ""
	}
	return false, "Coins not equal", exp.String(), got.String()
}

func coinsEq(exp, got sdk.Coins) (bool, string, string, string) {
	return exp.Equal(got), expectedStr, exp.String(), got.String()
}

func intEq(exp, got math.Int) (bool, string, string, string) {
	return exp.Equal(got), expectedStr, exp.String(), got.String()
}

func decEq(exp, got math.LegacyDec) (bool, string, string, string) {
	return exp.Equal(got), expectedStr, exp.String(), got.String()
}

func newInt(i int64) math.Int {
	return math.NewInt(i)
}
