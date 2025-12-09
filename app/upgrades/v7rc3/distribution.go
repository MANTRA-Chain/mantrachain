package v7rc3

import (
	"encoding/json"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	AMANTRA = "amantra"
)

var ScalingFactor = math.LegacyNewDec(4_000_000_000_000)

type Period struct {
	Period                uint64 `json:"period,string"`
	CumulativeRewardRatio string `json:"cumulative_reward_ratio,omitempty"`
}

func migrateDistr(ctx sdk.Context, distrKeeper distrkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, stakingKeeper stakingkeeper.Keeper) (err error) {
	ctx.Logger().Info("Starting v7rc3 distribution migration...")

	var dataBeforeUpgrade map[string]Period
	switch ctx.ChainID() {
	case "mantra-dukong-1":
		if err := json.Unmarshal(DukongBeforeUpgrade, &dataBeforeUpgrade); err != nil {
			return err
		}
	case "mantra-canary-net-1":
		if err := json.Unmarshal(CanaryBeforeUpgrade, &dataBeforeUpgrade); err != nil {
			return err
		}
	case "mantra-dryrun-1":
		if err := json.Unmarshal(DryrunBeforeUpgrade, &dataBeforeUpgrade); err != nil {
			return err
		}
	default:
		ctx.Logger().Info("No distribution migration data for this chain ID; skipping migration.")
		return nil
	}

	lastPeriodBeforeUpgrade := make(map[string]uint64)
	lastCumulativeRewardRatioBeforeUpgrade := make(map[string]math.LegacyDec)
	for valAddr, periodData := range dataBeforeUpgrade {
		lastPeriodBeforeUpgrade[valAddr] = periodData.Period
		dec, err := math.LegacyNewDecFromStr(periodData.CumulativeRewardRatio)
		if err != nil {
			return err
		}
		lastCumulativeRewardRatioBeforeUpgrade[valAddr] = dec
	}

	// Initialize maps to prevent nil map panics.
	reductionAmountPostUpgradePeriod := make(map[string]sdk.DecCoins)
	for valAddr, lastCumulativeReward := range lastCumulativeRewardRatioBeforeUpgrade {
		newAmount := lastCumulativeReward.Mul(ScalingFactor.Sub(math.LegacyOneDec()))
		reductionAmountPostUpgradePeriod[valAddr] = sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(AMANTRA, newAmount),
		)
	}

	ctx.Logger().Info("Step 1: Scaling down historical rewards...")
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
			// no cumulative reward pre upgrade so we leave it as is, nothing to subtract
			if _, ok := reductionAmountPostUpgradePeriod[valAddr.String()]; !ok {
				return false
			}
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

	ctx.Logger().Info("Step 2: Re-calculating all outstanding rewards from delegator state...")
	// Use a more efficient iteration pattern: iterate validators, then their delegations.
	newOutstandingRewards := make(map[string]sdk.DecCoins)
	err = stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		valAddrStr := validator.GetOperator()
		valAddr, err := stakingKeeper.ValidatorAddressCodec().StringToBytes(valAddrStr)
		if err != nil {
			return true
		}

		val, err := stakingKeeper.Validator(ctx, valAddr)
		if err != nil {
			// Returning true from the iterator will stop the process and propagate the error.
			return true
		}

		// Increment period ONCE per validator before calculations.
		endingPeriod, err := distrKeeper.IncrementValidatorPeriod(ctx, val)
		if err != nil {
			return true
		}

		// Initialize the rewards for this validator.
		newOutstandingRewards[valAddrStr] = sdk.NewDecCoins()

		// Iterate through all delegators of this specific validator.
		delegations, err := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
		if err != nil {
			return true
		}

		for _, delegation := range delegations {
			delAddr := delegation.GetDelegatorAddr()
			rewards, err := distrKeeper.CalculateDelegationRewards(ctx, val, delegation, endingPeriod)
			if err != nil {
				ctx.Logger().Error("failed to calculate delegation rewards", "validator", valAddrStr, "delegator", delAddr, "error", err)
				// Continue to the next delegator, but log the error.
				continue
			}
			newOutstandingRewards[valAddrStr] = newOutstandingRewards[valAddrStr].Add(rewards...)
		}

		return false // Continue to the next validator.
	})
	if err != nil {
		return err
	}

	ctx.Logger().Info("Step 3: Resetting outstanding rewards and calculating total top-up amount...")
	totalDistrBalanceNeeded := math.LegacyZeroDec()
	distrKeeper.IterateValidatorOutstandingRewards(ctx, func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) (stop bool) {
		// It's possible a validator has no delegations, so we check existence in the map.
		if newReward, ok := newOutstandingRewards[valAddr.String()]; ok {
			rewards.Rewards = newReward
			if err = distrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, rewards); err != nil {
				return true
			}
			totalDistrBalanceNeeded = totalDistrBalanceNeeded.Add(newReward.AmountOf(AMANTRA))
		}
		return false
	})
	if err != nil {
		return err
	}

	// add commissions
	distrKeeper.IterateValidatorAccumulatedCommissions(ctx, func(val sdk.ValAddress, commission distrtypes.ValidatorAccumulatedCommission) (stop bool) {
		totalDistrBalanceNeeded = totalDistrBalanceNeeded.Add(commission.Commission.AmountOf(AMANTRA))
		return false
	})

	ctx.Logger().Info("Step 4: Topping up x/distribution module account if needed...")
	currDistrBalance := bankKeeper.GetBalance(ctx, accountKeeper.GetModuleAddress(distrtypes.ModuleName), AMANTRA)
	requiredBalance := totalDistrBalanceNeeded.TruncateInt()

	if currDistrBalance.Amount.LT(requiredBalance) {
		topUpAmount := requiredBalance.Sub(currDistrBalance.Amount)
		topUpCoins := sdk.NewCoins(sdk.NewCoin(AMANTRA, topUpAmount))
		ctx.Logger().Info("Minting and sending funds to distribution module", "amount", topUpCoins.String())

		err = bankKeeper.MintCoins(ctx, minttypes.ModuleName, topUpCoins)
		if err != nil {
			return err
		}
		err = bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, distrtypes.ModuleName, topUpCoins)
		if err != nil {
			return err
		}
	} else {
		ctx.Logger().Info("Distribution module account balance is sufficient.", "balance", currDistrBalance.Amount, "required", requiredBalance)
	}

	ctx.Logger().Info("Distribution migration complete.", "total_balance_needed", totalDistrBalanceNeeded.String())
	panic("checking logs")
	return err
}
