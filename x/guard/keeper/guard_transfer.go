package keeper

import (
	"github.com/LimeChain/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetGuardTransfer set guardTransfer in the store
func (k Keeper) SetGuardTransfer(ctx sdk.Context, guardTransfer types.GuardTransfer) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GuardTransferKey))
	b := k.cdc.MustMarshal(&guardTransfer)
	store.Set([]byte{0}, b)
}

// GetGuardTransfer returns guardTransfer
func (k Keeper) GetGuardTransfer(ctx sdk.Context) (val types.GuardTransfer, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GuardTransferKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
