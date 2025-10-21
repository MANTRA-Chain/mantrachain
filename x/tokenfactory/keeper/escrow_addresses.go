package keeper

import (
	"cosmossdk.io/store/prefix"
	"github.com/MANTRA-Chain/mantrachain/v7/x/tokenfactory/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StoreEscrowAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.EscrowAddressKey)
	prefixStore.Set(address.Bytes(), []byte{0})
}

func (k Keeper) IsEscrowAddress(ctx sdk.Context, address sdk.AccAddress) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.EscrowAddressKey)
	bz := prefixStore.Get(address.Bytes())

	return len(bz) != 0
}

func (k Keeper) RemoveEscrowAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.EscrowAddressKey)
	prefixStore.Delete(address)
}
