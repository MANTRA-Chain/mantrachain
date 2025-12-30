package v7

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func migrateDistr(ctx sdk.Context, distrKeeper distrkeeper.Keeper) error {
	// migrate community pool
	feePool, err := distrKeeper.FeePool.Get(ctx)
	if err != nil {
		return err
	}

	newCommunityPool := convertDecCoinsToNewDenom(feePool.CommunityPool)
	feePool.CommunityPool = newCommunityPool

	if err := distrKeeper.FeePool.Set(ctx, feePool); err != nil {
		return err
	}

	// migrate validator outstanding rewards
	distrKeeper.IterateValidatorOutstandingRewards(ctx, func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) (stop bool) {
		rewards.Rewards = convertDecCoinsToNewDenom(rewards.Rewards)
		if err = distrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, rewards); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// migrate validator historical rewards
	distrKeeper.IterateValidatorHistoricalRewards(ctx, func(valAddr sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) (stop bool) {
		rewards.CumulativeRewardRatio = convertDecCoinsToNewDenomWithoutScaling(rewards.CumulativeRewardRatio)
		if err = distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, rewards); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// migrate validator current rewards
	distrKeeper.IterateValidatorCurrentRewards(ctx, func(valAddr sdk.ValAddress, rewards distrtypes.ValidatorCurrentRewards) (stop bool) {
		rewards.Rewards = convertDecCoinsToNewDenom(rewards.Rewards)
		if err = distrKeeper.SetValidatorCurrentRewards(ctx, valAddr, rewards); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// migrate validator accumulated commission
	distrKeeper.IterateValidatorAccumulatedCommissions(ctx, func(valAddr sdk.ValAddress, commission distrtypes.ValidatorAccumulatedCommission) (stop bool) {
		commission.Commission = convertDecCoinsToNewDenom(commission.Commission)
		if err = distrKeeper.SetValidatorAccumulatedCommission(ctx, valAddr, commission); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// migrate delegator starting infos
	distrKeeper.IterateDelegatorStartingInfos(ctx, func(valAddr sdk.ValAddress, delAddr sdk.AccAddress, info distrtypes.DelegatorStartingInfo) (stop bool) {
		info.Stake = info.Stake.MulInt(ScalingFactor)
		if err = distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, delAddr, info); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	return nil
}
