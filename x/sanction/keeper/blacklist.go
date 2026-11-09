package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
