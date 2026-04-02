package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func validateBlacklistMessage(authority string, accounts []string, maxAccounts int) error {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if len(accounts) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "accounts list cannot be empty")
	}

	if maxAccounts > 0 && len(accounts) > maxAccounts {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "accounts list cannot contain more than %d accounts", maxAccounts)
	}

	for _, account := range accounts {
		if _, err := sdk.AccAddressFromBech32(account); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address: %s", account)
		}
	}

	return nil
}
