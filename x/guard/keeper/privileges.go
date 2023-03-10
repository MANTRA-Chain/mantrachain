package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/LimeChain/mantrachain/x/guard/types"
)

func (k Keeper) HasPrivileges(ctx sdk.Context, address sdk.AccAddress, requiredPrivileges []byte) (bool, error) {
	conf := k.GetParams(ctx)

	accountPrivileges, _ := k.GetAccountPrivileges(ctx, address, conf.DefaultAccountPrivileges)
	accPr := types.AccPrivilegesFromBytes(accountPrivileges)

	if !accPr.Check(requiredPrivileges) {
		return false, errors.Wrapf(types.ErrInsufficientPrivileges, "insufficient privileges, address %s", address.String())
	}

	return true, nil
}
