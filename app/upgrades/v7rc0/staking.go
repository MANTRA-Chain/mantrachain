package v7rc0

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func migrateStaking(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper) error {
	stakingParams, err := stakingKeeper.GetParams(ctx)
	if err != nil {
		return err
	}
	stakingParams.BondDenom = AMANTRA
	err = stakingKeeper.SetParams(ctx, stakingParams)
	if err != nil {
		return err
	}

	// migrate validators
	err = stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		val, ok := validator.(stakingtypes.Validator)
		if !ok {
			// this should not happen
			err = errorsmod.Wrapf(sdkerrors.ErrInvalidType, "expected validator, got %T", validator)
			return true
		}

		val.Tokens = val.Tokens.Mul(ScalingFactor)
		val.DelegatorShares = val.DelegatorShares.MulInt(ScalingFactor)

		if err = stakingKeeper.SetValidator(ctx, val); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// migrate delegations
	err = stakingKeeper.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) (stop bool) {
		delegation.Shares = delegation.Shares.MulInt(ScalingFactor)
		if err = stakingKeeper.SetDelegation(ctx, delegation); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// migrate unbonding delegations
	err = stakingKeeper.IterateUnbondingDelegations(ctx, func(index int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].Balance = ubd.Entries[i].Balance.Mul(ScalingFactor)
			ubd.Entries[i].InitialBalance = ubd.Entries[i].InitialBalance.Mul(ScalingFactor)
		}
		if err = stakingKeeper.SetUnbondingDelegation(ctx, ubd); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	// migrate redelegations
	err = stakingKeeper.IterateRedelegations(ctx, func(index int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].InitialBalance = red.Entries[i].InitialBalance.Mul(ScalingFactor)
			red.Entries[i].SharesDst = red.Entries[i].SharesDst.MulInt(ScalingFactor)
		}
		if err = stakingKeeper.SetRedelegation(ctx, red); err != nil {
			return true
		}
		return false
	})
	if err != nil {
		return err
	}

	return nil
}
