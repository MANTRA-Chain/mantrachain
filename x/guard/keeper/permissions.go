package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CheckIsAdmin(ctx sdk.Context, address string) error {
	conf := k.GetParams(ctx)
	admin := sdk.MustAccAddressFromBech32(conf.AdminAccount)

	if !admin.Equals(sdk.MustAccAddressFromBech32(address)) {
		return errors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", address)
	}

	return nil
}
