package v7rc0

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

const (
	blocked_account = "mantra14t56rzvxzw0yp9plcf9dy6rr53chyvxt4cqtt5"
)

func migrateBlockedAccount(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, distrKeeper distributionkeeper.Keeper) error {
	delAddr, err := sdk.AccAddressFromBech32(blocked_account)
	if err != nil {
		return err
	}

	// Withdraw all rewards for the blocked account before unbonding.
	validators, err := stakingKeeper.GetDelegatorValidators(ctx, delAddr, 1000)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get delegator validators for blocked account")
	}

	for _, validator := range validators.Validators {
		valAddress := sdk.MustValAddressFromBech32(validator.OperatorAddress)
		_, err := distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddress)
		if err != nil {
			return errorsmod.Wrapf(err, "failed to withdraw rewards from validator %s for blocked account", validator.GetOperator())
		}
	}

	delegations, err := stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get all delegations from blocked account")
	}

	for _, delegation := range delegations {
		valAddress, err := sdk.ValAddressFromBech32(delegation.GetValidatorAddr())
		if err != nil {
			return errorsmod.Wrap(err, "failed to parse validator address from delegation")
		}
		_, amount, err := stakingKeeper.Undelegate(ctx, delAddr, valAddress, delegation.GetShares())
		if err != nil {
			return errorsmod.Wrapf(err, "failed to undelegate from validator %s for blocked account", delegation.GetValidatorAddr())
		}
		ctx.Logger().Info("Undelegated from", "validator", delegation.GetValidatorAddr(), "shares", delegation.GetShares(), "amount", amount.String())
	}

	return nil
}
