package v7rc3

import (
	"context"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
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
	AMANTRA       = "amantra"
	V7RC2_UPGRADE = "v7.0.0-rc2"
)

var ScalingFactor = math.LegacyNewDec(4_000_000_000_000)

type Period struct {
	Period                uint64 `json:"period"`
	CumulativeRewardRatio string `json:"cumulative_reward_ratio,omitempty"`
}

func delegationTotalRewards(ctx context.Context, delegatorAddress string, distrKeeper distrkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, stakingKeeper stakingkeeper.Keeper) (total sdk.DecCoins, err error) {
	delAdr, err := accountKeeper.AddressCodec().StringToBytes(delegatorAddress)
	if err != nil {
		return nil, err
	}

	total = sdk.DecCoins{}
	err = stakingKeeper.IterateDelegations(
		ctx, delAdr,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr, err := stakingKeeper.ValidatorAddressCodec().StringToBytes(del.GetValidatorAddr())
			if err != nil {
				panic(err)
			}

			val, err := stakingKeeper.Validator(ctx, valAddr)
			if err != nil {
				panic(err)
			}

			endingPeriod, err := distrKeeper.IncrementValidatorPeriod(ctx, val)
			if err != nil {
				panic(err)
			}

			delReward, err := distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
			if err != nil {
				panic(err)
			}
			total = total.Add(delReward...)
			return false
		},
	)
	if err != nil {
		return nil, err
	}
	return total, nil
}

func migrateDistr(ctx sdk.Context, cms storetypes.CommitMultiStore, distrKeeper distrkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper, stakingKeeper stakingkeeper.Keeper, upgradeKeeper *upgradekeeper.Keeper) (err error) {
	ctx.Logger().Info("Starting v7rc3 distribution migration...")

	v7rc2Height, err := upgradeKeeper.GetDoneHeight(ctx, V7RC2_UPGRADE)
	if err != nil {
		ctx.Logger().Warn("Could not query v7rc2 upgrade height", "error", err)
	} else if v7rc2Height > 0 {
		ctx.Logger().Info("Found v7rc2 upgrade height", "height", v7rc2Height)
	} else {
		ctx.Logger().Warn("v7rc2 upgrade not found in done upgrades")
	}

	corruptedValidators := make(map[string]uint64)

	if v7rc2Height > 1 {
		ctx.Logger().Info("Checking for corrupted validators by comparing rewards across v7rc2 upgrade boundary")

		heightBeforeUpgrade := v7rc2Height - 1

		cmsBeforeVersion, err := cms.CacheMultiStoreWithVersion(heightBeforeUpgrade)
		if err != nil {
			distrKeeper.IterateValidatorHistoricalRewards(ctx, func(valAddr sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) (stop bool) {
				valAddrStr := valAddr.String()
				if len(rewards.CumulativeRewardRatio) == 0 {
					return false
				}
				ratio, err := math.LegacyNewDecFromStr(rewards.CumulativeRewardRatio[0].Amount.String())
				if err != nil {
					return false
				}
				if ratio.GT(math.LegacyOneDec()) {
					if _, exists := corruptedValidators[valAddrStr]; !exists {
						if period > 0 {
							corruptedValidators[valAddrStr] = period - 1
						} else {
							corruptedValidators[valAddrStr] = 0
						}
					}
				}
				return false
			})
		} else {
			cmsAtVersion, err := cms.CacheMultiStoreWithVersion(v7rc2Height)
			if err != nil {
				ctx.Logger().Error("Failed to load multistore at upgrade height", "height", v7rc2Height, "error", err)
			} else {
				ctxBefore := ctx.WithMultiStore(cmsBeforeVersion).WithBlockHeight(heightBeforeUpgrade)
				ctxAtUpgrade := ctx.WithMultiStore(cmsAtVersion).WithBlockHeight(v7rc2Height)

				delegatorSet := make(map[string]bool)
				stakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
					valAddrStr := val.GetOperator()
					valAddr, err := stakingKeeper.ValidatorAddressCodec().StringToBytes(valAddrStr)
					if err != nil {
						return false
					}

					delegations, err := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
					if err != nil {
						return false
					}

					for _, del := range delegations {
						delegatorSet[del.GetDelegatorAddr()] = true
					}
					return false
				})

				ctx.Logger().Info("Checking rewards for delegators", "count", len(delegatorSet))

				for delAddr := range delegatorSet {
					rewardsBefore, err := delegationTotalRewards(ctxBefore, delAddr, distrKeeper, accountKeeper, stakingKeeper)
					if err != nil {
						ctx.Logger().Debug("Failed to calculate rewards before upgrade", "delegator", delAddr, "error", err)
						continue
					}

					rewardsAt, err := delegationTotalRewards(ctxAtUpgrade, delAddr, distrKeeper, accountKeeper, stakingKeeper)
					if err != nil {
						ctx.Logger().Debug("Failed to calculate rewards at upgrade", "delegator", delAddr, "error", err)
						continue
					}

					rewardsBeforeAmt := math.LegacyZeroDec()
					for _, coin := range rewardsBefore {
						if coin.Denom == "uom" {
							rewardsBeforeAmt = rewardsBeforeAmt.Add(coin.Amount)
						}
					}
					rewardsAtAmt := rewardsAt.AmountOf(AMANTRA)

					if rewardsBeforeAmt.IsPositive() {
						scaledRewardsBefore := rewardsBeforeAmt.Mul(ScalingFactor)
						diff := rewardsAtAmt.Quo(scaledRewardsBefore)

						ctx.Logger().Debug("Rewards comparison across v7rc2 upgrade",
							"delegator", delAddr,
							"rewards_before", rewardsBefore.String(),
							"rewards_at", rewardsAt.String(),
							"diff", diff.String())

						threshold := math.LegacyNewDec(2)
						if diff.GT(threshold) {
							ctx.Logger().Info("Detected corruption via delegator rewards",
								"delegator", delAddr,
								"rewards_before", rewardsBefore.String(),
								"rewards_at", rewardsAt.String(),
								"diff", diff.String())

							delAddrBytes, _ := accountKeeper.AddressCodec().StringToBytes(delAddr)
							stakingKeeper.IterateDelegations(ctx, delAddrBytes, func(_ int64, del stakingtypes.DelegationI) (stop bool) {
								valAddrStr := del.GetValidatorAddr()
								if _, exists := corruptedValidators[valAddrStr]; !exists {
									valAddr, _ := stakingKeeper.ValidatorAddressCodec().StringToBytes(valAddrStr)
									currentRewardsBefore, err := distrKeeper.GetValidatorCurrentRewards(ctxBefore, valAddr)
									if err == nil {
										corruptedValidators[valAddrStr] = currentRewardsBefore.Period
										ctx.Logger().Info("Marked validator as corrupted", "validator", valAddrStr, "lastGoodPeriod", currentRewardsBefore.Period)
									}
								}
								return false
							})
							break
						}
					}
				}
			}
		}
	} else {
		ctx.Logger().Warn("v7rc2 upgrade height not found, skipping corruption detection")
	}
	ctx.Logger().Info("Scaling down corrupted historical rewards...", "count", len(corruptedValidators))

	reductionAmountPostUpgradePeriod := make(map[string]sdk.DecCoins)

	for valAddrStr, lastGoodPeriod := range corruptedValidators {
		valAddr, _ := sdk.ValAddressFromBech32(valAddrStr)

		var lastGoodRatio math.LegacyDec
		var foundLastGood bool

		if lastGoodPeriod > 0 {
			histRewards, err := distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, lastGoodPeriod)
			if err == nil && len(histRewards.CumulativeRewardRatio) > 0 {
				lastGoodRatio, err = math.LegacyNewDecFromStr(histRewards.CumulativeRewardRatio[0].Amount.String())
				if err == nil {
					foundLastGood = true
				}
			}
		}

		if !foundLastGood {
			lastGoodRatio = math.LegacyZeroDec()
		}

		scaledDiff := lastGoodRatio.Mul(ScalingFactor.Sub(math.LegacyOneDec()))
		reductionAmountPostUpgradePeriod[valAddrStr] = sdk.NewDecCoins(sdk.NewDecCoinFromDec(AMANTRA, scaledDiff))

		distrKeeper.IterateValidatorHistoricalRewards(ctx, func(iterValAddr sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) (stop bool) {
			if !iterValAddr.Equals(valAddr) {
				return false
			}

			if period <= lastGoodPeriod {
				for i := range rewards.CumulativeRewardRatio {
					rewards.CumulativeRewardRatio[i].Amount = rewards.CumulativeRewardRatio[i].Amount.Quo(ScalingFactor)
				}
			} else {
				rewards.CumulativeRewardRatio = rewards.CumulativeRewardRatio.Sub(reductionAmountPostUpgradePeriod[valAddrStr])
			}

			err = distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, rewards)
			return err != nil
		})
	}

	ctx.Logger().Info("Re-calculating all outstanding rewards from delegator state...")
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
				continue
			}
			newOutstandingRewards[valAddrStr] = newOutstandingRewards[valAddrStr].Add(rewards...)
		}

		return false
	})
	if err != nil {
		return err
	}

	ctx.Logger().Info("Resetting outstanding rewards and calculating total top-up amount...")
	totalOutstandingRewards := math.LegacyZeroDec()
	distrKeeper.IterateValidatorOutstandingRewards(ctx, func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) (stop bool) {
		if newReward, ok := newOutstandingRewards[valAddr.String()]; ok {
			rewards.Rewards = newReward
			if err = distrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, rewards); err != nil {
				return true
			}
			totalOutstandingRewards = totalOutstandingRewards.Add(newReward.AmountOf(AMANTRA))
		}
		return false
	})
	if err != nil {
		return err
	}

	ctx.Logger().Info("Topping up x/distribution module account if needed...")
	currDistrBalance := bankKeeper.GetBalance(ctx, accountKeeper.GetModuleAddress(distrtypes.ModuleName), AMANTRA)
	requiredBalance := totalOutstandingRewards.TruncateInt()

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

	ctx.Logger().Info("Distribution migration complete.", "final_outstanding_rewards", totalOutstandingRewards.String())

	return nil
}
