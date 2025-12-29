package v7

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

func migrateAuth(ctx sdk.Context, authKeeper authkeeper.AccountKeeper) error {
	iter, err := authKeeper.Accounts.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.Valid() {
		acc, err := iter.Value()
		if err != nil {
			return err
		}

		var baseAcc *vestingtypes.BaseVestingAccount
		switch vestingAcc := acc.(type) {
		case *vestingtypes.ContinuousVestingAccount:
			baseAcc = vestingAcc.BaseVestingAccount
		case *vestingtypes.DelayedVestingAccount:
			baseAcc = vestingAcc.BaseVestingAccount
		case *vestingtypes.PermanentLockedAccount:
			baseAcc = vestingAcc.BaseVestingAccount
		case *vestingtypes.PeriodicVestingAccount:
			baseAcc = vestingAcc.BaseVestingAccount
			for i := range vestingAcc.VestingPeriods {
				vestingAcc.VestingPeriods[i].Amount = convertCoinsToNewDenom(vestingAcc.VestingPeriods[i].Amount)
			}
		}
		if baseAcc != nil {
			baseAcc.OriginalVesting = convertCoinsToNewDenom(baseAcc.OriginalVesting)
			baseAcc.DelegatedFree = convertCoinsToNewDenom(baseAcc.DelegatedFree)
			baseAcc.DelegatedVesting = convertCoinsToNewDenom(baseAcc.DelegatedVesting)
			authKeeper.SetAccount(ctx, acc)
		}
		iter.Next()
	}

	return nil
}
