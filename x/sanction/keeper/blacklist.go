package keeper

import (
	"context"

	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IsSanctioned reports whether addr is on the sanction list.
func (k Keeper) IsSanctioned(ctx context.Context, addr string) (bool, error) {
	return k.BlacklistAccounts.Has(ctx, addr)
}

// emitListChanged emits one event for a message that changed the list. It
// emits nothing if no address changed.
func (k Keeper) emitListChanged(ctx context.Context, addrs []string, sanctioned bool, actor string, role types.SanctionerRole, reason string) error {
	if len(addrs) == 0 {
		return nil
	}
	return sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSanctionListChanged{
		Addresses:  addrs,
		Sanctioned: sanctioned,
		Actor:      actor,
		ActorRole:  role,
		Reason:     reason,
	})
}
