package keeper

import (
	"github.com/MANTRA-Finance/aumega/x/guard/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetGuardTransferCoins(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GuardTransferCoinsKey))
	store.Set([]byte{0}, types.Placeholder)
}

func (k Keeper) RemoveGuardTransferCoins(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GuardTransferCoinsKey))
	store.Delete([]byte{0})
}

func (k Keeper) HasGuardTransferCoins(
	ctx sdk.Context,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GuardTransferCoinsKey))
	return store.Has([]byte{0})
}
