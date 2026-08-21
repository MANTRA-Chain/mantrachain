package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddToBlacklist adds each of the given accounts to the blacklist, returning
// an error if any of them is already blacklisted.
func (k Keeper) AddToBlacklist(ctx context.Context, accounts []string) error {
	for _, account := range accounts {
		hasAccount, err := k.BlacklistAccounts.Has(ctx, account)
		if err != nil {
			return err
		}
		if hasAccount {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s has already been blacklisted", account)
		}

		if err := k.BlacklistAccounts.Set(ctx, account); err != nil {
			return err
		}
	}

	return nil
}

// RemoveFromBlacklist removes each of the given accounts from the blacklist,
// returning an error if any of them is not currently blacklisted.
func (k Keeper) RemoveFromBlacklist(ctx context.Context, accounts []string) error {
	for _, account := range accounts {
		hasAccount, err := k.BlacklistAccounts.Has(ctx, account)
		if err != nil {
			return err
		}
		if !hasAccount {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "blacklist account %s is not blacklisted", account)
		}

		if err := k.BlacklistAccounts.Remove(ctx, account); err != nil {
			return err
		}
	}

	return nil
}
