package keeper

import (
	"cosmossdk.io/store/prefix"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetGuardTransferCoins(ctx sdk.Context) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.GuardTransferCoinsKey))
	store.Set([]byte{0}, types.Placeholder)
}

func (k Keeper) RemoveGuardTransferCoins(ctx sdk.Context) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.GuardTransferCoinsKey))
	store.Delete([]byte{0})
}

func (k Keeper) HasGuardTransferCoins(
	ctx sdk.Context,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.GuardTransferCoinsKey))
	return store.Has([]byte{0})
}
