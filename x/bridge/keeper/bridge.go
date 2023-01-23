package keeper

import (
	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetBridge(ctx sdk.Context, bridge types.Bridge) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgeStoreKey(bridge.Creator))
	b := k.cdc.MustMarshal(&bridge)
	store.Set(bridge.Index, b)
}

func (k Keeper) GetBridge(
	ctx sdk.Context,
	creator sdk.AccAddress,
	index []byte,
) (val types.Bridge, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgeStoreKey(creator))

	if !k.HasBridge(ctx, creator, index) {
		return types.Bridge{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasBridge(
	ctx sdk.Context,
	creator sdk.AccAddress,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgeStoreKey(creator))
	return store.Has(index)
}

func (k Keeper) GetAllBridge(ctx sdk.Context) (list []types.Bridge) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgeStoreKey(nil))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Bridge
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
