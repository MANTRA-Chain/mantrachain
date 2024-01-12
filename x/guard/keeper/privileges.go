package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/AumegaChain/aumega/x/guard/types"
)

func (k Keeper) CheckAccountFulfillsRequiredPrivileges(ctx sdk.Context, address sdk.AccAddress, requiredPrivilegesRaw []byte) (bool, error) {
	conf := k.GetParams(ctx)

	accountPrivileges, found := k.GetAccountPrivileges(ctx, address, conf.DefaultPrivileges)

	if !found {
		return false, errors.Wrapf(types.ErrAccountRequiredPrivilegesNotFound, "account privileges not found, address %s", address.String())
	}

	accPr := types.PrivilegesFromBytes(accountPrivileges)

	defaultPrivileges := big.NewInt(0).SetBytes(conf.DefaultPrivileges)
	inverseDefaultPrilileges := big.NewInt(0).Not(defaultPrivileges)
	accountPrivilegesWithoutDefault := big.NewInt(0).And(inverseDefaultPrilileges, accPr.BigInt())

	if accountPrivilegesWithoutDefault.Cmp(big.NewInt(0)) == 0 {
		return false, errors.Wrapf(types.ErrAccountRequiredPrivilegesNotSet, "account privileges not set, address %s", address.String())
	}

	requiredPrivileges := types.PrivilegesFromBytes(requiredPrivilegesRaw)

	if !accPr.CheckPrivileges(requiredPrivileges, conf.DefaultPrivileges) {
		return false, nil
	}

	return true, nil
}
