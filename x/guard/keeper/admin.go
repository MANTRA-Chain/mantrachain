package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetAdmin(ctx sdk.Context) sdk.AccAddress {
	conf := k.GetParams(ctx)
	return sdk.MustAccAddressFromBech32(conf.AdminAccount)
}

func (k Keeper) CheckIsAdmin(ctx sdk.Context, address string) error {
	admin := k.GetAdmin(ctx)

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	if !admin.Equals(acc) {
		return errors.Wrapf(errorstypes.ErrUnauthorized, "unauthorized account %s", address)
	}

	return nil
}
