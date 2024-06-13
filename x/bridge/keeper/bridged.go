package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetBridged set a specific bridged in the store from its index
func (k Keeper) SetBridged(ctx sdk.Context, bridged types.Bridged) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BridgedKeyPrefix))
	b := k.cdc.MustMarshal(&bridged)
	store.Set(bridged.Index, b)
}

func (k Keeper) HasBridged(
	ctx sdk.Context,
	ethTxHash string,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BridgedKeyPrefix))
	return store.Has(types.BridgedKey(
		ethTxHash,
	))
}

// GetBridged returns a bridged from its index
func (k Keeper) GetBridged(
	ctx sdk.Context,
	ethTxHash string,

) (val types.Bridged, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BridgedKeyPrefix))

	b := store.Get(types.BridgedKey(
		ethTxHash,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllBridged returns all bridged
func (k Keeper) GetAllBridged(ctx sdk.Context) (list []types.Bridged) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BridgedKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Bridged
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
