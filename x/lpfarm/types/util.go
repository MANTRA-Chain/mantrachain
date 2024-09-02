package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	utils "github.com/MANTRA-Finance/mantrachain/types"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/amm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func DeriveFarmingPoolAddress(planId uint64) sdk.AccAddress {
	return address.Module(ModuleName, []byte(fmt.Sprintf("FarmingPool/%d", planId)))
}

func DeriveFarmingReserveAddress(denom string) sdk.AccAddress {
	return address.Module(ModuleName, []byte(fmt.Sprintf("FarmingReserve/%s", denom)))
}

func RewardsForBlock(rewardsPerDay sdk.Coins, blockDuration time.Duration) sdk.DecCoins {
	return sdk.NewDecCoinsFromCoins(rewardsPerDay...).
		MulDecTruncate(math.LegacyNewDec(blockDuration.Milliseconds())).
		QuoDecTruncate(math.LegacyNewDec(day.Milliseconds()))
}

// PoolRewardWeight returns given pool's reward weight.
func PoolRewardWeight(pool amm.Pool) (weight math.LegacyDec) {
	rx, ry := pool.Balances()
	sqrt := utils.DecApproxSqrt
	switch pool := pool.(type) {
	case *amm.BasicPool:
		weight = sqrt(math.LegacyNewDecFromInt(rx.Mul(ry)))
	case *amm.RangedPool:
		transX, transY := pool.Translation()
		weight = sqrt(transX.Add(math.LegacyNewDecFromInt(rx))).Mul(sqrt(transY.Add(math.LegacyNewDecFromInt(ry))))
	default:
		panic("invalid pool type")
	}
	return
}
