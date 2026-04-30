package v8_1

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type stakingView interface {
	Validator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.ValidatorI, error)
	Delegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.DelegationI, error)
}

// fixSilentlySkippedSlashes patches DelegatorStartingInfo.Stake by the residue
// not explained by recorded ValidatorSlashEvents. Recorded slashes are handled
// at withdraw time, so info.Stake > currentStake alone is not enough to clamp.
func fixSilentlySkippedSlashes(
	ctx sdk.Context,
	stakingKeeper stakingView,
	distrKeeper distrkeeper.Keeper,
) error {
	logger := ctx.Logger()
	endingHeight := uint64(ctx.BlockHeight())

	type fix struct {
		val      sdk.ValAddress
		del      sdk.AccAddress
		info     distrtypes.DelegatorStartingInfo
		newStake math.LegacyDec
	}
	var fixes []fix
	var iterErr error

	distrKeeper.IterateDelegatorStartingInfos(ctx, func(
		val sdk.ValAddress,
		del sdk.AccAddress,
		info distrtypes.DelegatorStartingInfo,
	) (stop bool) {
		validator, err := stakingKeeper.Validator(ctx, val)
		if err != nil {
			iterErr = fmt.Errorf("read validator %s: %w", val, err)
			return true
		}
		if validator == nil {
			logger.Warn(
				"v8.1.0: validator missing for delegator starting info; skipping",
				"validator", val.String(),
				"delegator", del.String(),
			)
			return false
		}
		delegation, err := stakingKeeper.Delegation(ctx, del, val)
		if err != nil {
			iterErr = fmt.Errorf("read delegation %s/%s: %w", del, val, err)
			return true
		}
		if delegation == nil {
			logger.Warn(
				"v8.1.0: delegation missing for delegator starting info; skipping",
				"validator", val.String(),
				"delegator", del.String(),
			)
			return false
		}
		currentStake := validator.TokensFromShares(delegation.GetShares())

		expectedStake := info.Stake
		hasEvents := false
		distrKeeper.IterateValidatorSlashEventsBetween(ctx, val, info.Height, endingHeight,
			func(_ uint64, ev distrtypes.ValidatorSlashEvent) bool {
				hasEvents = true
				expectedStake = expectedStake.MulTruncate(
					math.LegacyOneDec().Sub(ev.Fraction))
				return false
			})

		if !expectedStake.GT(currentStake) {
			return false
		}

		var newStake math.LegacyDec
		if hasEvents {
			ratio := currentStake.Quo(expectedStake)
			newStake = info.Stake.MulTruncate(ratio)
		} else {
			// silent-only residue — exact assignment, no rounding drift.
			newStake = currentStake
		}
		fixes = append(fixes, fix{val: val, del: del, info: info, newStake: newStake})
		return false
	})
	if iterErr != nil {
		return iterErr
	}

	for _, f := range fixes {
		logger.Info(
			"v8.1.0: clamping silent-skip residue on delegator starting info",
			"validator", f.val.String(),
			"delegator", f.del.String(),
			"old_stake", f.info.Stake.String(),
			"new_stake", f.newStake.String(),
		)
		f.info.Stake = f.newStake
		if err := distrKeeper.SetDelegatorStartingInfo(ctx, f.val, f.del, f.info); err != nil {
			return err
		}
	}

	logger.Info("v8.1.0: fixed silently-skipped slashes", "entries_updated", len(fixes))
	return nil
}
