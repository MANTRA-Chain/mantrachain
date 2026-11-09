package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SetSanctioners replaces the whole sanctioner list. The list must already be
// validated with types.ValidateSanctioners.
func (k Keeper) SetSanctioners(ctx context.Context, sanctioners []types.Sanctioner) error {
	if err := k.Sanctioners.Clear(ctx, nil); err != nil {
		return err
	}
	for _, s := range sanctioners {
		addr, err := types.NormalizeAddress(s.Address)
		if err != nil {
			return err
		}
		s.Address = addr
		if err := k.Sanctioners.Set(ctx, addr, s); err != nil {
			return err
		}
	}
	return nil
}

// GetSanctioners returns the whole sanctioner list.
func (k Keeper) GetSanctioners(ctx context.Context) ([]types.Sanctioner, error) {
	iter, err := k.Sanctioners.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	return iter.Values()
}

// actor is whoever signed a sanction message: the gov authority or a
// sanctioner.
type actor struct {
	address string
	role    types.SanctionerRole // UNSPECIFIED for governance
	gov     bool
}

// resolveActor maps a signer to the gov authority or a sanctioner, failing
// for anyone else.
func (k Keeper) resolveActor(ctx context.Context, signer string) (actor, error) {
	addr, err := types.NormalizeAddress(signer)
	if err != nil {
		return actor{}, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid signer: %s", signer)
	}
	if addr == k.authority {
		return actor{address: addr, role: types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED, gov: true}, nil
	}
	s, err := k.getSanctioner(ctx, addr)
	if err != nil {
		return actor{}, err
	}
	return actor{address: s.Address, role: s.Role}, nil
}

// SanctionAddresses adds addresses to the sanction list on behalf of the gov
// authority or a sanctioner with any role. Addresses already listed are
// skipped. For sanctioners it enforces the caps and fails if any address is
// protected; governance is subject to neither.
func (k Keeper) SanctionAddresses(ctx context.Context, signer string, addresses []string, reason string) error {
	a, err := k.resolveActor(ctx, signer)
	if err != nil {
		return err
	}
	if !a.gov {
		if err := k.checkLimits(ctx, len(addresses), reason); err != nil {
			return err
		}
	}

	height := sdk.UnwrapSDKContext(ctx).BlockHeight()
	for _, address := range addresses {
		addr, err := types.NormalizeAddress(address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", address)
		}
		if !a.gov {
			if err := k.checkNotProtected(ctx, addr); err != nil {
				return err
			}
		}
		has, err := k.BlacklistAccounts.Has(ctx, addr)
		if err != nil {
			return err
		}
		if has {
			continue
		}

		record := types.SanctionRecord{
			Actor:     a.address,
			ActorRole: a.role,
			Reason:    reason,
			Height:    height,
		}
		if err := k.addEntry(ctx, addr, record); err != nil {
			return err
		}
	}
	return nil
}

// UnsanctionAddresses removes addresses from the sanction list on behalf of
// the gov authority or a full sanctioner. Addresses not listed are skipped.
// Governance can remove any entry. A full sanctioner is capped and fails if
// any address was added by governance, including entries that predate
// sanctioners.
func (k Keeper) UnsanctionAddresses(ctx context.Context, signer string, addresses []string, reason string) error {
	a, err := k.resolveActor(ctx, signer)
	if err != nil {
		return err
	}
	if !a.gov {
		if !a.role.CanUnsanction() {
			return errorsmod.Wrapf(types.ErrRoleNotPermitted, "%s has role %s and cannot remove addresses", a.address, a.role)
		}
		if err := k.checkLimits(ctx, len(addresses), reason); err != nil {
			return err
		}
	}

	for _, address := range addresses {
		addr, err := types.NormalizeAddress(address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", address)
		}
		has, err := k.BlacklistAccounts.Has(ctx, addr)
		if err != nil {
			return err
		}
		if !has {
			continue
		}

		if !a.gov {
			if err := k.checkNotGovernanceEntry(ctx, addr); err != nil {
				return err
			}
		}

		if err := k.removeEntry(ctx, addr, a.address, a.role, reason); err != nil {
			return err
		}
	}
	return nil
}

// checkNotGovernanceEntry fails if addr was added by governance or has no
// record because it predates sanctioners.
func (k Keeper) checkNotGovernanceEntry(ctx context.Context, addr string) error {
	record, err := k.SanctionRecords.Get(ctx, addr)
	if errors.Is(err, collections.ErrNotFound) {
		return errorsmod.Wrapf(types.ErrGovernanceEntry, "%s", addr)
	}
	if err != nil {
		return err
	}
	if record.ActorRole == types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED {
		return errorsmod.Wrapf(types.ErrGovernanceEntry, "%s", addr)
	}
	return nil
}

// RevokeSanctioner removes an add-only sanctioner from the list on behalf of
// a full sanctioner. Full sanctioners can only be removed by governance.
func (k Keeper) RevokeSanctioner(ctx context.Context, signer, target, reason string) error {
	sanctioner, err := k.getSanctioner(ctx, signer)
	if err != nil {
		return err
	}
	if !sanctioner.Role.CanUnsanction() {
		return errorsmod.Wrapf(types.ErrRoleNotPermitted, "%s has role %s and cannot revoke sanctioners", sanctioner.Address, sanctioner.Role)
	}
	if err := k.checkLimits(ctx, 1, reason); err != nil {
		return err
	}

	addr, err := types.NormalizeAddress(target)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", target)
	}
	revoked, err := k.Sanctioners.Get(ctx, addr)
	if errors.Is(err, collections.ErrNotFound) {
		return errorsmod.Wrapf(types.ErrNotSanctioner, "%s", addr)
	}
	if err != nil {
		return err
	}
	if revoked.Role != types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY {
		return errorsmod.Wrapf(types.ErrRoleNotPermitted, "%s has role %s; only add-only sanctioners can be revoked", addr, revoked.Role)
	}

	if err := k.Sanctioners.Remove(ctx, addr); err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.EventManager().EmitTypedEvent(&types.EventSanctionerRevoked{
		Address:   addr,
		RevokedBy: sanctioner.Address,
		Reason:    reason,
		Height:    sdkCtx.BlockHeight(),
	})
}

func (k Keeper) getSanctioner(ctx context.Context, signer string) (types.Sanctioner, error) {
	addr, err := types.NormalizeAddress(signer)
	if err != nil {
		return types.Sanctioner{}, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid signer: %s", signer)
	}
	s, err := k.Sanctioners.Get(ctx, addr)
	if errors.Is(err, collections.ErrNotFound) {
		return types.Sanctioner{}, errorsmod.Wrapf(types.ErrNotSanctioner, "%s", addr)
	}
	return s, err
}

func (k Keeper) checkLimits(ctx context.Context, numAddresses int, reason string) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	if uint64(numAddresses) > uint64(params.MaxAddressesPerMsg) {
		return errorsmod.Wrapf(types.ErrTooManyAddresses, "%d addresses, max %d", numAddresses, params.MaxAddressesPerMsg)
	}
	if uint64(len(reason)) > uint64(params.MaxReasonLength) {
		return errorsmod.Wrapf(types.ErrReasonTooLong, "%d bytes, max %d", len(reason), params.MaxReasonLength)
	}
	return nil
}

// checkNotProtected fails if a sanctioner may not sanction addr: the gov
// authority, any module account, or anyone on the sanctioner list.
func (k Keeper) checkNotProtected(ctx context.Context, addr string) error {
	if _, ok := k.protectedAddrs[addr]; ok {
		return errorsmod.Wrapf(types.ErrProtectedAddress, "%s is the gov authority or a module account", addr)
	}
	isSanctioner, err := k.Sanctioners.Has(ctx, addr)
	if err != nil {
		return err
	}
	if isSanctioner {
		return errorsmod.Wrapf(types.ErrProtectedAddress, "%s is a sanctioner", addr)
	}
	if k.accountKeeper != nil {
		acc := k.accountKeeper.GetAccount(ctx, sdk.MustAccAddressFromBech32(addr))
		if _, ok := acc.(sdk.ModuleAccountI); ok {
			return errorsmod.Wrapf(types.ErrProtectedAddress, "%s is a module account", addr)
		}
	}
	return nil
}
