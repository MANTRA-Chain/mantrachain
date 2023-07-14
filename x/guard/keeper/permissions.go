package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CheckIsAdmin(ctx sdk.Context, address string) error {
	conf := k.GetParams(ctx)
	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	if !admin.Equals(acc) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized account %s", address)
	}

	return nil
}
