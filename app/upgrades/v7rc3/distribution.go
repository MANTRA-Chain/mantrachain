package v7rc3

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
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
	OutstandingRewardBeforeUpgrade = map[string]sdk.DecCoins{
		"mantra1...": sdk.NewDecCoins(sdk.NewDecCoinFromDec(AMANTRA, math.LegacyMustNewDecFromStr("4000"))),
		"mantra2...": sdk.NewDecCoins(sdk.NewDecCoinFromDec(AMANTRA, math.LegacyMustNewDecFromStr("2000"))),
	}
	ScalingFactor = math.LegacyNewDec(4_000_000_000_000)
)

func migrateDistr(ctx sdk.Context, distrKeeper distrkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper) (err error) {
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

	// reset outstanding rewards to right after upgrade
	totalOutstandingRewards := math.LegacyZeroDec()
	distrKeeper.IterateValidatorOutstandingRewards(ctx, func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) (stop bool) {
		newOutstandingRewards := OutstandingRewardBeforeUpgrade[valAddr.String()].MulDec(ScalingFactor)
		rewards.Rewards = newOutstandingRewards
		if err = distrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, rewards); err != nil {
			return true
		}
		totalOutstandingRewards = totalOutstandingRewards.Add(newOutstandingRewards.AmountOf(AMANTRA))
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
