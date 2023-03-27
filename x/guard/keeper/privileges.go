package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
)

func (k Keeper) CheckAccountFulfillsRequiredPrivileges(ctx sdk.Context, address sdk.AccAddress, requiredPrivilegesList [][]byte) (bool, error) {
	conf := k.GetParams(ctx)

	accountPrivileges, _ := k.GetAccountPrivileges(ctx, address, conf.DefaultPrivileges)
	accPr := types.PrivilegesFromBytes(accountPrivileges)

	requiredPrivileges := types.NewEmptyPrivileges()
	requiredPrivileges.МergeMore(requiredPrivilegesList)

	if !accPr.CheckPrivileges(requiredPrivileges) {
		return false, errors.Wrapf(types.ErrInsufficientPrivileges, "insufficient privileges, address %s", address.String())
	}

	return true, nil
}
