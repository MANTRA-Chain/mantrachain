package amm_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	utils "github.com/AumegaChain/aumega/types"
	"github.com/AumegaChain/aumega/x/liquidity/amm"
)

func TestMatchableAmount(t *testing.T) {
	order1 := newOrder(amm.Buy, utils.ParseDec("1.0"), math.NewInt(10000))
	for _, tc := range []struct {
		order    amm.Order
		price    sdk.Dec
		expected math.Int
	}{
		{order1, utils.ParseDec("1"), math.NewInt(10000)},
		{order1, utils.ParseDec("0.01"), math.NewInt(10000)},
		{order1, utils.ParseDec("100"), math.NewInt(100)},
		{order1, utils.ParseDec("100.1"), math.NewInt(99)},
		{order1, utils.ParseDec("9999"), math.NewInt(1)},
		{order1, utils.ParseDec("10001"), math.NewInt(0)},
	} {
		t.Run("", func(t *testing.T) {
			require.True(math.IntEq(t, tc.expected, amm.MatchableAmount(tc.order, tc.price)))
		})
	}
}

type batchIdOrderer struct {
	batchId uint64
}

func (orderer *batchIdOrderer) Order(dir amm.OrderDirection, price sdk.Dec, amt math.Int) amm.Order {
	return &batchIdOrder{newOrder(dir, price, amt), orderer.batchId}
}

type batchIdOrder struct {
	amm.Order
	batchId uint64
}

func (order *batchIdOrder) GetBatchId() uint64 {
	return order.batchId
}

func TestGroupOrdersByBatchId(t *testing.T) {
	price := utils.ParseDec("1.0")
	newOrder := func(amt math.Int, batchId uint64) amm.Order {
		return (&batchIdOrderer{batchId}).Order(amm.Buy, price, amt)
	}
	orders := []amm.Order{
		newOrder(math.NewInt(32000), 0),
		newOrder(math.NewInt(8000), 4),
		newOrder(math.NewInt(1000), 1),
		newOrder(math.NewInt(16000), 4),
		newOrder(math.NewInt(4000), 2),
		newOrder(math.NewInt(2000), 1),
		newOrder(math.NewInt(64000), 0),
	}
	groups := amm.GroupOrdersByBatchId(orders)
	require.EqualValues(t, 1, groups[0].BatchId)
	require.True(math.IntEq(t, math.NewInt(3000), amm.TotalAmount(groups[0].Orders)))
	require.EqualValues(t, 2, groups[1].BatchId)
	require.True(math.IntEq(t, math.NewInt(4000), amm.TotalAmount(groups[1].Orders)))
	require.EqualValues(t, 4, groups[2].BatchId)
	require.True(math.IntEq(t, math.NewInt(24000), amm.TotalAmount(groups[2].Orders)))
	require.EqualValues(t, 0, groups[3].BatchId)
	require.True(math.IntEq(t, math.NewInt(96000), amm.TotalAmount(groups[3].Orders)))
}
