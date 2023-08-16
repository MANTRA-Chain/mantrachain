package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/liquidfarming/types"
)

func TestLiquidFarm(t *testing.T) {
	liquidFarm := types.LiquidFarm{
		PoolId:        1,
		MinFarmAmount: sdk.ZeroInt(),
		MinBidAmount:  sdk.ZeroInt(),
		FeeRate:       sdk.ZeroDec(),
	}
	require.Equal(t, `fee_rate: "0.000000000000000000"
min_bid_amount: "0"
min_farm_amount: "0"
pool_id: "1"
`, liquidFarm.String())
}

func TestLiquidFarmCoinDenom(t *testing.T) {
	for _, tc := range []struct {
		denom      string
		expectsErr bool
	}{
		{"lf1", false},
		{"lf10", false},
		{"lf18446744073709551615", false},
		{"lf18446744073709551616", true},
		{"lfabc", true},
		{"lf01", true},
		{"lf-10", true},
		{"lf+10", true},
		{"ucre", true},
		{"denom1", true},
	} {
		t.Run("", func(t *testing.T) {
			poolId, err := types.ParseLiquidFarmCoinDenom(tc.denom)
			if tc.expectsErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.denom, types.LiquidFarmCoinDenom(poolId))
			}
		})
	}
}

func TestLiquidFarmReserveAddress(t *testing.T) {
	config := sdk.GetConfig()
	addrPrefix := config.GetBech32AccountAddrPrefix()

	for _, tc := range []struct {
		poolId   uint64
		expected string
	}{
		{1, addrPrefix + "1zyyf855slxure4c8dr06p00qjnkem95d2lgv8wgvry2rt437x6tsaf9tcf"},
		{2, addrPrefix + "1d2csu4ynxpuxll8wk72n9z98ytm649u78paj9efskjwrlc2wyhpq8h886j"},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.expected, types.LiquidFarmReserveAddress(tc.poolId).String())
		})
	}
}

func TestCalculateLiquidFarmAmount(t *testing.T) {
	for _, tc := range []struct {
		name              string
		lfTotalSupplyAmt  math.Int
		lpTotalFarmingAmt math.Int
		newFarmingAmt     math.Int
		expectedAmt       math.Int
	}{
		{
			name:              "initial minting",
			lfTotalSupplyAmt:  sdk.ZeroInt(),
			lpTotalFarmingAmt: sdk.ZeroInt(),
			newFarmingAmt:     math.NewInt(1_000_00_000),
			expectedAmt:       math.NewInt(1_000_00_000),
		},
		{
			name:              "normal",
			lfTotalSupplyAmt:  math.NewInt(1_000_000_000),
			lpTotalFarmingAmt: math.NewInt(1_000_000_000),
			newFarmingAmt:     math.NewInt(250_000_000),
			expectedAmt:       math.NewInt(250_000_000),
		},
		{
			name:              "rewards are auto compounded",
			lfTotalSupplyAmt:  math.NewInt(1_000_000_000),
			lpTotalFarmingAmt: math.NewInt(1_100_000_000),
			newFarmingAmt:     math.NewInt(100_000_000),
			expectedAmt:       math.NewInt(90_909_090),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mintingAmt := types.CalculateLiquidFarmAmount(
				tc.lfTotalSupplyAmt,
				tc.lpTotalFarmingAmt,
				tc.newFarmingAmt,
			)
			require.Equal(t, tc.expectedAmt, mintingAmt)
		})
	}
}

func TestCalculateLiquidUnfarmAmount(t *testing.T) {
	for _, tc := range []struct {
		name               string
		lfTotalSupplyAmt   math.Int
		lpTotalFarmingAmt  math.Int
		unfarmingAmt       math.Int
		compoundingRewards math.Int
		expectedAmt        math.Int
	}{
		{
			name:               "unfarm all",
			lfTotalSupplyAmt:   math.NewInt(100_000_000),
			lpTotalFarmingAmt:  math.NewInt(100_000_000),
			unfarmingAmt:       math.NewInt(100_000_000),
			compoundingRewards: sdk.ZeroInt(),
			expectedAmt:        math.NewInt(100_000_000),
		},
		{
			name:               "unfarming small amount #1: no compounding rewards",
			lfTotalSupplyAmt:   math.NewInt(100_000_000),
			lpTotalFarmingAmt:  math.NewInt(100_000_000),
			unfarmingAmt:       math.NewInt(1),
			compoundingRewards: sdk.ZeroInt(),
			expectedAmt:        math.NewInt(1),
		},
		{
			name:               "unfarming small amount #2: with compounding rewards",
			lfTotalSupplyAmt:   math.NewInt(100_000_000),
			lpTotalFarmingAmt:  math.NewInt(100_000_100),
			unfarmingAmt:       math.NewInt(1),
			compoundingRewards: math.NewInt(100),
			expectedAmt:        math.NewInt(1),
		},
		{
			name:               "rewards are auto compounded #1: no compouding rewards",
			lfTotalSupplyAmt:   math.NewInt(1_000_000_000),
			lpTotalFarmingAmt:  math.NewInt(1_100_000_000),
			unfarmingAmt:       math.NewInt(100_000_000),
			compoundingRewards: sdk.ZeroInt(),
			expectedAmt:        math.NewInt(110_000_000),
		},
		{
			name:               "rewards are auto compounded #1: with compouding rewards",
			lfTotalSupplyAmt:   math.NewInt(1_000_000_000),
			lpTotalFarmingAmt:  math.NewInt(1_100_000_000),
			unfarmingAmt:       math.NewInt(100_000_000),
			compoundingRewards: math.NewInt(100_000),
			expectedAmt:        math.NewInt(109_990_000),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			unfarmingAmt := types.CalculateLiquidUnfarmAmount(
				tc.lfTotalSupplyAmt,
				tc.lpTotalFarmingAmt,
				tc.unfarmingAmt,
				tc.compoundingRewards,
			)
			require.Equal(t, tc.expectedAmt, unfarmingAmt)
		})
	}
}

func TestDeductFees(t *testing.T) {
	for _, tc := range []struct {
		name     string
		feeRate  sdk.Dec
		rewards  sdk.Coins
		deducted sdk.Coins
		fees     sdk.Coins
	}{
		{
			name:     "zero fee rate",
			feeRate:  sdk.ZeroDec(),
			rewards:  utils.ParseCoins("100denom1"),
			deducted: utils.ParseCoins("100denom1"),
			fees:     sdk.Coins{},
		},
		{
			name:     "fee rate - 10%",
			feeRate:  math.LegacyMustNewDecFromStr("0.1"),
			rewards:  utils.ParseCoins("100denom1"),
			deducted: utils.ParseCoins("90denom1"),
			fees:     utils.ParseCoins("10denom1"),
		},
		{
			name:     "fee rate - 6.666666666666%",
			feeRate:  math.LegacyMustNewDecFromStr("0.066666666666666"),
			rewards:  utils.ParseCoins("100000denom1"),
			deducted: utils.ParseCoins("93333denom1"),
			fees:     utils.ParseCoins("6667denom1"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			deducted, fees := types.DeductFees(tc.rewards, tc.feeRate)
			require.Equal(t, tc.deducted, deducted)
			require.Equal(t, tc.fees, fees)
		})
	}
}
