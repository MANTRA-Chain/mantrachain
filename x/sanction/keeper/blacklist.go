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

// AddToBlacklist adds each of the given accounts to the blacklist on behalf
// of governance, returning an error if any of them is already blacklisted.
func (k Keeper) AddToBlacklist(ctx context.Context, accounts []string) error {
	height := sdk.UnwrapSDKContext(ctx).BlockHeight()
	for _, account := range accounts {
		addr, err := types.NormalizeAddress(account)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address: %s", account)
		}
		hasAccount, err := k.BlacklistAccounts.Has(ctx, addr)
		if err != nil {
			return err
		}
		if hasAccount {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s has already been blacklisted", addr)
		}

		record := types.SanctionRecord{
			Actor:     k.authority,
			ActorRole: types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED,
			Height:    height,
		}
		if err := k.addEntry(ctx, addr, record); err != nil {
			return err
		}
	}

	return nil
}

// RemoveFromBlacklist removes each of the given accounts from the blacklist
// on behalf of governance, returning an error if any of them is not currently
// blacklisted. Governance can remove any entry, whoever added it.
func (k Keeper) RemoveFromBlacklist(ctx context.Context, accounts []string) error {
	for _, account := range accounts {
		addr, err := types.NormalizeAddress(account)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address: %s", account)
		}
		hasAccount, err := k.BlacklistAccounts.Has(ctx, addr)
		if err != nil {
			return err
		}
		if !hasAccount {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "blacklist account %s is not blacklisted", addr)
		}

		if err := k.removeEntry(ctx, addr, k.authority, types.SanctionerRole_SANCTIONER_ROLE_UNSPECIFIED, ""); err != nil {
			return err
		}
	}

	return nil
}

// IsSanctioned reports whether addr is on the sanction list, with its record
// if it has one.
func (k Keeper) IsSanctioned(ctx context.Context, addr string) (bool, *types.SanctionRecord, error) {
	has, err := k.BlacklistAccounts.Has(ctx, addr)
	if err != nil || !has {
		return false, nil, err
	}
	record, err := k.SanctionRecords.Get(ctx, addr)
	if errors.Is(err, collections.ErrNotFound) {
		return true, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, &record, nil
}

// addEntry puts addr on the sanction list with its record and emits the
// event. The caller has already checked addr isn't listed.
func (k Keeper) addEntry(ctx context.Context, addr string, record types.SanctionRecord) error {
	if err := k.BlacklistAccounts.Set(ctx, addr); err != nil {
		return err
	}
	if err := k.SanctionRecords.Set(ctx, addr, record); err != nil {
		return err
	}
	return sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSanctionListChanged{
		Address:    addr,
		Sanctioned: true,
		Actor:      record.Actor,
		ActorRole:  record.ActorRole,
		Reason:     record.Reason,
		Height:     record.Height,
	})
}

// removeEntry takes addr off the sanction list, drops its record and emits
// the event. The caller has already checked addr is listed.
func (k Keeper) removeEntry(ctx context.Context, addr, actor string, role types.SanctionerRole, reason string) error {
	if err := k.BlacklistAccounts.Remove(ctx, addr); err != nil {
		return err
	}
	if err := k.SanctionRecords.Remove(ctx, addr); err != nil {
		return err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.EventManager().EmitTypedEvent(&types.EventSanctionListChanged{
		Address:    addr,
		Sanctioned: false,
		Actor:      actor,
		ActorRole:  role,
		Reason:     reason,
		Height:     sdkCtx.BlockHeight(),
	})
}
