package v7rc0

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	blocked_account = "mantra14t56rzvxzw0yp9plcf9dy6rr53chyvxt4cqtt5"
)

func migrateBlockedAccount(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, bankKeeper bankkeeper.Keeper, distrKeeper distributionkeeper.Keeper) error {
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

	// Unbond all delegations immediately.
	if err := unbondAllDelegationsImmediately(ctx, stakingKeeper, bankKeeper, delAddr); err != nil {
		return errorsmod.Wrap(err, "failed to unbond all delegations from blocked account")
	}

	return nil
}

// This is conceptual code for a migration. Do not use as-is without careful testing.
func performImmediateUnbond(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, bankKeeper bankkeeper.Keeper, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesToUnbond math.LegacyDec) error {
	// 1. Call the staking keeper's Unbond method. This correctly handles the delegation
	//    and validator state, but moves tokens to the NotBondedPool.
	unbondedTokens, err := stakingKeeper.Unbond(ctx, delAddr, valAddr, sharesToUnbond)
	if err != nil {
		return err
	}

	validator, _ := stakingKeeper.GetValidator(ctx, valAddr)

	// 2. Manually move the tokens from BondedPool to NotBondedPool if the validator was bonded.
	//    The Unbond method above only returns the amount; the pool transfer happens in Undelegate.
	if validator.IsBonded() {
		if err := bankKeeper.SendCoinsFromModuleToModule(
			ctx, types.BondedPoolName, types.NotBondedPoolName, sdk.NewCoins(sdk.NewCoin(UOM, unbondedTokens)),
		); err != nil {
			return err
		}
	}

	// 3. Now, immediately move the tokens from the NotBondedPool to the user's account.
	//    This bypasses the waiting period.
	coins := sdk.NewCoins(sdk.NewCoin(UOM, unbondedTokens))
	if err := bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.NotBondedPoolName, delAddr, coins); err != nil {
		return err
	}

	return nil
}

// unbondAllDelegationsImmediately iterates through all of a delegator's delegations and performs an immediate unbond for each.
func unbondAllDelegationsImmediately(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, bankKeeper bankkeeper.Keeper, delAddr sdk.AccAddress) error {
	delegations, err := stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
	if err != nil {
		return err
	}

	if len(delegations) == 0 {
		return nil // No delegations to unbond.
	}

	for _, delegation := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			return errorsmod.Wrap(err, "failed to parse validator address")
		}

		if err := performImmediateUnbond(ctx, stakingKeeper, bankKeeper, delAddr, valAddr, delegation.Shares); err != nil {
			return errorsmod.Wrapf(err, "failed to immediately unbond from validator %s", delegation.ValidatorAddress)
		}
	}

	return nil
}
