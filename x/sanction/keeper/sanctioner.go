package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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

// The functions below assume the caller has already checked the actor may
// make the change. They apply the limits and the rules on target addresses,
// so another module can call them directly.

// SanctionAddresses adds addresses to the sanction list on behalf of actor,
// skipping any already listed. Nobody can sanction a module account, and
// only governance can sanction a sanctioner.
func (k Keeper) SanctionAddresses(ctx context.Context, addresses []string, actor, reason string) error {
	if err := k.checkLimits(ctx, len(addresses), reason); err != nil {
		return err
	}
	addrs, err := normalizeAddresses(addresses)
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if err := k.CheckNotModuleAccount(ctx, addr); err != nil {
			return err
		}
		if actor == k.authority {
			continue
		}
		if err := k.checkNotProtected(ctx, addr); err != nil {
			return err
		}
	}
	role, err := k.actorRole(ctx, actor)
	if err != nil {
		return err
	}

	var added []string
	for _, addr := range addrs {
		has, err := k.BlacklistAccounts.Has(ctx, addr)
		if err != nil {
			return err
		}
		if has {
			continue
		}
		if err := k.BlacklistAccounts.Set(ctx, addr); err != nil {
			return err
		}
		if err := k.clearCodeDelegation(ctx, addr); err != nil {
			return err
		}
		added = append(added, addr)
	}
	return k.emitListChanged(ctx, added, true, actor, role, reason)
}

// UnsanctionAddresses removes addresses from the sanction list on behalf of
// actor, skipping any not listed.
func (k Keeper) UnsanctionAddresses(ctx context.Context, addresses []string, actor, reason string) error {
	if err := k.checkLimits(ctx, len(addresses), reason); err != nil {
		return err
	}
	addrs, err := normalizeAddresses(addresses)
	if err != nil {
		return err
	}
	return k.unsanction(ctx, addrs, actor, reason)
}

// unsanction removes the canonical addresses that are listed and emits one
// event for them. It applies no limits.
func (k Keeper) unsanction(ctx context.Context, addrs []string, actor, reason string) error {
	role, err := k.actorRole(ctx, actor)
	if err != nil {
		return err
	}

	var removed []string
	for _, addr := range addrs {
		has, err := k.BlacklistAccounts.Has(ctx, addr)
		if err != nil {
			return err
		}
		if !has {
			continue
		}
		if err := k.BlacklistAccounts.Remove(ctx, addr); err != nil {
			return err
		}
		removed = append(removed, addr)
	}
	return k.emitListChanged(ctx, removed, false, actor, role, reason)
}

// AppointSanctioners replaces the sanctioner list and unsanctions every
// address that joins it: one sanctioned while its appointment was voting
// couldn't sign anything. A sanctioner already on the list keeps any sanction
// governance gave it.
func (k Keeper) AppointSanctioners(ctx context.Context, sanctioners []types.Sanctioner) error {
	var appointed []string
	for _, s := range sanctioners {
		addr, err := types.NormalizeAddress(s.Address)
		if err != nil {
			return err
		}
		has, err := k.Sanctioners.Has(ctx, addr)
		if err != nil {
			return err
		}
		if !has {
			appointed = append(appointed, addr)
		}
	}
	if err := k.SetSanctioners(ctx, sanctioners); err != nil {
		return err
	}
	return k.unsanction(ctx, appointed, k.authority, types.ReasonAppointed)
}

// RevokeSanctioner removes target from the sanctioner list on behalf of
// revokedBy. Only add-only sanctioners can be revoked; full sanctioners are
// removed by governance with MsgUpdateSanctioners.
func (k Keeper) RevokeSanctioner(ctx context.Context, target, revokedBy, reason string) error {
	if err := k.checkLimits(ctx, 1, reason); err != nil {
		return err
	}
	addr, err := types.NormalizeAddress(target)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", target)
	}
	s, err := k.getSanctioner(ctx, addr)
	if err != nil {
		return err
	}
	if s.Role != types.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY {
		return errorsmod.Wrapf(types.ErrRoleNotPermitted, "%s has role %s; only add-only sanctioners can be revoked", addr, s.Role)
	}

	if err := k.Sanctioners.Remove(ctx, addr); err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.EventManager().EmitTypedEvent(&types.EventSanctionerRevoked{
		Address:   addr,
		RevokedBy: revokedBy,
		Reason:    reason,
		Height:    sdkCtx.BlockHeight(),
	})
}

// clearCodeDelegation removes an EIP-7702 delegation from addr, so code the
// account authorized before it was sanctioned can no longer run as it. The
// ante handler stops it authorizing new code. Contract code is left alone.
// The EVM clears a delegation the same way, by deleting the code hash, when
// it applies an authorization to the zero address.
func (k Keeper) clearCodeDelegation(ctx context.Context, addr string) error {
	if k.evmKeeper == nil {
		return nil
	}
	acc, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}
	if len(acc) != common.AddressLength {
		// Contract and interchain accounts have 32-byte addresses and no EVM
		// code; read as one, the last 20 bytes would name another account.
		return nil
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ethAddr := common.BytesToAddress(acc)
	codeHash := k.evmKeeper.GetCodeHash(sdkCtx, ethAddr)
	if evmtypes.IsEmptyCodeHash(codeHash.Bytes()) {
		return nil
	}
	if _, ok := ethtypes.ParseDelegation(k.evmKeeper.GetCode(sdkCtx, codeHash)); ok {
		k.evmKeeper.DeleteCodeHash(sdkCtx, ethAddr)
	}
	return nil
}

// actorRole returns actor's sanctioner role for events, or UNSPECIFIED for
// governance or any other actor that isn't a sanctioner.
func (k Keeper) actorRole(ctx context.Context, actor string) (types.SanctionerRole, error) {
	s, err := k.Sanctioners.Get(ctx, actor)
	if errors.Is(err, collections.ErrNotFound) {
		return types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED, nil
	}
	return s.Role, err
}

// getSanctioner returns the sanctioner entry for addr, with any role.
func (k Keeper) getSanctioner(ctx context.Context, addr string) (types.Sanctioner, error) {
	s, err := k.Sanctioners.Get(ctx, addr)
	if errors.Is(err, collections.ErrNotFound) {
		return types.Sanctioner{}, errorsmod.Wrapf(types.ErrNotSanctioner, "%s", addr)
	}
	return s, err
}

// checkLimits applies the address and reason caps to every actor,
// governance included.
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

// CheckNotModuleAccount fails if addr is a module account. The bank send
// restriction would stop it sending, which could halt the chain (for
// example, the fee collector paying out to distribution in BeginBlock).
func (k Keeper) CheckNotModuleAccount(ctx context.Context, addr string) error {
	isModule, err := k.isModuleAccount(ctx, addr)
	if err != nil {
		return err
	}
	if isModule {
		return errorsmod.Wrapf(types.ErrProtectedAddress, "%s is a module account", addr)
	}
	return nil
}

// isModuleAccount reports whether addr belongs to a module, whether or not
// its account exists yet.
func (k Keeper) isModuleAccount(ctx context.Context, addr string) (bool, error) {
	acc, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return false, err
	}
	for _, perms := range k.accountKeeper.GetModulePermissions() {
		if perms.GetAddress().Equals(acc) {
			return true, nil
		}
	}
	_, isModule := k.accountKeeper.GetAccount(ctx, acc).(sdk.ModuleAccountI)
	return isModule, nil
}

// checkNotProtected fails if addr may only be sanctioned by governance, which
// is anyone on the sanctioner list.
func (k Keeper) checkNotProtected(ctx context.Context, addr string) error {
	isSanctioner, err := k.Sanctioners.Has(ctx, addr)
	if err != nil {
		return err
	}
	if isSanctioner {
		return errorsmod.Wrapf(types.ErrProtectedAddress, "%s is a sanctioner", addr)
	}
	return nil
}

// normalizeAddresses returns the addresses in canonical form, so they match
// what the ante handler looks up.
func normalizeAddresses(addresses []string) ([]string, error) {
	out := make([]string, len(addresses))
	for i, address := range addresses {
		addr, err := types.NormalizeAddress(address)
		if err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", address)
		}
		out[i] = addr
	}
	return out, nil
}
