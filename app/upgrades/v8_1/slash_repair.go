package v8_1

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func repairSilentSlashes(
	ctx sdk.Context,
	stakingKeeper stakingView,
	distrKeeper distrkeeper.Keeper,
	records []SilentSlashRecord,
) error {
	affected := affectedValidators(records)
	if len(affected) == 0 {
		ctx.Logger().Info("v8.1.0: no silent-slash records configured; skipping reset repair")
		return nil
	}

	logger := ctx.Logger()
	upgradeHeight := uint64(ctx.BlockHeight())

	for _, r := range records {
		logger.Info("v8.1.0: silent-slash record",
			"validator", r.Operator,
			"height", r.Height,
			"fraction", r.Fraction,
			"reason", r.Reason,
		)
	}

	periodByVal := make(map[string]uint64, len(affected))
	for valStr := range affected {
		valAddr, err := sdk.ValAddressFromBech32(valStr)
		if err != nil {
			return fmt.Errorf("v8.1.0 repair: bad operator %q: %w", valStr, err)
		}
		validator, err := stakingKeeper.Validator(ctx, valAddr)
		if err != nil {
			return fmt.Errorf("v8.1.0 repair: read validator %s: %w", valStr, err)
		}
		if validator == nil {
			logger.Warn("v8.1.0: affected validator not found; skipping reset", "validator", valStr)
			continue
		}
		closedPeriod, err := distrKeeper.IncrementValidatorPeriod(ctx, validator)
		if err != nil {
			return fmt.Errorf("v8.1.0 repair: increment period for %s: %w", valStr, err)
		}
		periodByVal[valStr] = closedPeriod + 1
	}

	type reset struct {
		val  sdk.ValAddress
		del  sdk.AccAddress
		info distrtypes.DelegatorStartingInfo
	}
	var resets []reset
	var iterErr error

	distrKeeper.IterateDelegatorStartingInfos(ctx, func(
		val sdk.ValAddress,
		del sdk.AccAddress,
		info distrtypes.DelegatorStartingInfo,
	) (stop bool) {
		newPeriod, ok := periodByVal[val.String()]
		if !ok {
			return false
		}
		validator, err := stakingKeeper.Validator(ctx, val)
		if err != nil {
			iterErr = fmt.Errorf("read validator %s: %w", val, err)
			return true
		}
		if validator == nil {
			logger.Warn("v8.1.0: affected validator missing during reset; skipping",
				"validator", val.String(), "delegator", del.String())
			return false
		}
		delegation, err := stakingKeeper.Delegation(ctx, del, val)
		if err != nil {
			iterErr = fmt.Errorf("read delegation %s/%s: %w", del, val, err)
			return true
		}
		if delegation == nil {
			logger.Warn("v8.1.0: delegation missing during reset; leaving starting info alone",
				"validator", val.String(), "delegator", del.String())
			return false
		}
		info.Stake = validator.TokensFromShares(delegation.GetShares())
		info.PreviousPeriod = newPeriod
		info.Height = upgradeHeight
		resets = append(resets, reset{val: val, del: del, info: info})
		return false
	})
	if iterErr != nil {
		return iterErr
	}

	for _, r := range resets {
		if err := distrKeeper.SetDelegatorStartingInfo(ctx, r.val, r.del, r.info); err != nil {
			return fmt.Errorf("v8.1.0 repair: set starting info %s/%s: %w", r.val, r.del, err)
		}
	}

	logger.Info(
		"v8.1.0: reset starting info on silent-slash validators",
		"validators", len(periodByVal),
		"entries_reset", len(resets),
		"upgrade_height", upgradeHeight,
	)
	return nil
}
