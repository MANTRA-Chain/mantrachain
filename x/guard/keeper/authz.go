package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MANTRA-Finance/aumega/x/guard/types"
)

func (k Keeper) CheckHasAuthz(ctx sdk.Context, address string, authz string) error {
	conf := k.GetParams(ctx)
	admin := k.GetAdmin(ctx)
	authzBytes := []byte(authz)

	acc, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	if admin.Equals(acc) {
		return nil
	}

	privileges, found := k.GetRequiredPrivileges(ctx, authzBytes, types.RequiredPrivilegesAuthz)

	if !found || privileges == nil || len(privileges) == 0 {
		return sdkerrors.Wrapf(types.ErrAuthzRequiredPrivilegesNotFound, "authz required privileges not found, authz %s", authz)
	}

	defaultPrivileges := big.NewInt(0).SetBytes(conf.DefaultPrivileges)
	inverseDefaultPrilileges := big.NewInt(0).Not(defaultPrivileges)
	requiredPrivileges := types.PrivilegesFromBytes(privileges)
	requiredPrivilegesWithoutDefault := big.NewInt(0).And(inverseDefaultPrilileges, requiredPrivileges.BigInt())

	if requiredPrivilegesWithoutDefault.Cmp(big.NewInt(0)) == 0 {
		return sdkerrors.Wrapf(types.ErrAuthzRequiredPrivilegesNotSet, "authz required privileges not set, authz %s", authz)
	}

	hasPrivileges, err := k.CheckAccountFulfillsRequiredPrivileges(ctx, acc, privileges)

	if err != nil {
		return err
	}

	if !hasPrivileges {
		k.Logger(ctx).Error("insufficient privileges", "address", address, "authz", authz)
		return sdkerrors.Wrapf(types.ErrInsufficientPrivileges, "insufficient privileges, address %s, authz %s", address, authz)
	}

	return nil
}
