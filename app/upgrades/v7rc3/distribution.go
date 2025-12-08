package v7rc3

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const (
	UOM     = "uom"
	AMANTRA = "amantra"
	MANTRA  = "mantra"
)

var (
	lastPeriodBeforeUpgrade = map[string]uint64{
		"mantra1...": 15,
		"mantra2...": 20,
	}
	lastCumulativeRewardRationBeforeUpgrade = map[string]math.LegacyDec{
		"mantra1...": math.LegacyMustNewDecFromStr("1.2"), // 1.0
		"mantra2...": math.LegacyMustNewDecFromStr("0.5"), // 0.5
	}
	ScalingFactor = math.LegacyNewDec(4_000_000_000_000)
)

func migrateDistr(ctx sdk.Context, distrKeeper distrkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, stakingKeeper stakingkeeper.Keeper) (err error) {
	var reductionAmountPostUpgradePeriod map[string]sdk.DecCoins
	for valAddr, lastCumulativeReward := range lastCumulativeRewardRationBeforeUpgrade {
		newAmount := lastCumulativeReward.Mul(ScalingFactor.Sub(math.LegacyOneDec()))
		reductionAmountPostUpgradePeriod[valAddr] = sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(AMANTRA, newAmount),
		)
	}

	// scale down historical rewards before upgrade by scaling factor
	// scale down historical rewards after upgrade by subtracting the last pre-upgrade period scale amount
	distrKeeper.IterateValidatorHistoricalRewards(ctx, func(valAddr sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) (stop bool) {
		lastPeriod, ok := lastPeriodBeforeUpgrade[valAddr.String()]
		if !ok {
			return false
		}
		if period <= lastPeriod {
			rewards.CumulativeRewardRatio = rewards.CumulativeRewardRatio.QuoDec(ScalingFactor)
		} else {
			rewards.CumulativeRewardRatio = rewards.CumulativeRewardRatio.Sub(reductionAmountPostUpgradePeriod[valAddr.String()])
		}
		if err = distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, rewards); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// recalculate all outstanding rewards
	var endingPeriods map[string]uint64
	var newOutstandingRewards map[string]sdk.DecCoins
	distrKeeper.IterateDelegatorStartingInfos(ctx, func(valAddr sdk.ValAddress, delAddr sdk.AccAddress, info distrtypes.DelegatorStartingInfo) (stop bool) {
		val, err := stakingKeeper.Validator(ctx, valAddr)
		if err != nil {
			return true
		}
		// end current period and calculate rewards
		endingPeriod := uint64(0)
		if endingPeriods[valAddr.String()] == 0 {
			endingPeriod, err := distrKeeper.IncrementValidatorPeriod(ctx, val)
			if err != nil {
				return true
			}
			endingPeriods[valAddr.String()] = endingPeriod
		} else {
			endingPeriod = endingPeriods[valAddr.String()]
		}

		del, err := stakingKeeper.Delegation(ctx, delAddr, valAddr)
		if err != nil {
			return true
		}

		rewards, err := distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
		if err != nil {
			return true
		}
		if newOutstandingRewards[valAddr.String()] == nil {
			newOutstandingRewards[valAddr.String()] = rewards
		} else {
			newOutstandingRewards[valAddr.String()] = newOutstandingRewards[valAddr.String()].Add(rewards...)
		}
		return false
	})
	if err != nil {
		return err
	}

	// reset outstanding rewards to new outstanding rewards
	totalOutstandingRewards := math.LegacyZeroDec()
	distrKeeper.IterateValidatorOutstandingRewards(ctx, func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) (stop bool) {
		newOutstandingReward := newOutstandingRewards[valAddr.String()]
		rewards.Rewards = newOutstandingReward
		if err = distrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, rewards); err != nil {
			return true
		}
		totalOutstandingRewards = totalOutstandingRewards.Add(newOutstandingReward.AmountOf(AMANTRA))
		return false
	})
	if err != nil {
		return err
	}

	// top up x/distribution module account if needed
	currDistrBalance := bankKeeper.GetBalance(ctx, accountKeeper.GetModuleAddress(distrtypes.ModuleName), AMANTRA)
	if currDistrBalance.Amount.LT(totalOutstandingRewards.TruncateInt()) {
		topUpAmount := totalOutstandingRewards.TruncateInt().Sub(currDistrBalance.Amount)
		topUpCoins := sdk.NewCoins(sdk.NewCoin(AMANTRA, topUpAmount))
		err = bankKeeper.MintCoins(ctx, minttypes.ModuleName, topUpCoins)
		if err != nil {
			return err
		}
		err = bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, distrtypes.ModuleName, topUpCoins)
		if err != nil {
			return err
		}
	}

	ctx.Logger().Info("Topping up x/rewards module account", "amount", totalOutstandingRewards.String())

	return nil
}
