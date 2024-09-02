package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetWhitelistTransfersAccAddrs(
	ctx sdk.Context,
	whitelistTransfersAccAddrs types.WhitelistTransfersAccAddrs,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.WhitelistTransfersAccAddrsStoreKey())
	b := k.cdc.MustMarshal(&whitelistTransfersAccAddrs)
	store.Set(whitelistTransfersAccAddrs.Index, b)
}

func (k Keeper) GetWhitelistTransfersAccAddrs(
	ctx sdk.Context,
	index []byte,
) (val types.WhitelistTransfersAccAddrs, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.WhitelistTransfersAccAddrsStoreKey())

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllWhitelistTransfersAccAddrs(ctx sdk.Context) (list []types.WhitelistTransfersAccAddrs) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.WhitelistTransfersAccAddrsStoreKey())
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.WhitelistTransfersAccAddrs
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) IsTransfersAccAddrsWhitelisted(
	ctx sdk.Context,
	index []byte,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.WhitelistTransfersAccAddrsStoreKey())

	b := store.Get(index)
	if b == nil {
		return false
	}

	var val types.WhitelistTransfersAccAddrs
	k.cdc.MustUnmarshal(b, &val)
	return val.IsWhitelisted
}
