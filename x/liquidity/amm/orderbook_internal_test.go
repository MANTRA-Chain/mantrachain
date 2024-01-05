package amm

import (
	"sort"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/MANTRA-Finance/aumega/types"
)

func newOrder(dir OrderDirection, price sdk.Dec, amt math.Int) *BaseOrder {
	return NewBaseOrder(dir, price, amt, OfferCoinAmount(dir, price, amt))
}

func TestOrderBookTicks_add(t *testing.T) {
	prices := []sdk.Dec{
		utils.ParseDec("1.0"),
		utils.ParseDec("1.1"),
		utils.ParseDec("1.05"),
		utils.ParseDec("1.1"),
		utils.ParseDec("1.2"),
		utils.ParseDec("0.9"),
		utils.ParseDec("0.9"),
	}
	var ticks orderBookTicks
	for _, price := range prices {
		ticks.addOrder(newOrder(Buy, price, math.NewInt(10000)))
	}
	pricesSet := map[string]struct{}{}
	for _, price := range prices {
		pricesSet[price.String()] = struct{}{}
	}
	prices = nil
	for priceStr := range pricesSet {
		prices = append(prices, utils.ParseDec(priceStr))
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].GT(prices[j])
	})
	for i, price := range prices {
		require.True(math.LegacyDecEq(t, price, ticks.ticks[i].price))
	}
}
