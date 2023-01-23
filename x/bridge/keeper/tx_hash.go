package keeper

import (
	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTxHash set a specific txHash in the store from its index
func (k Keeper) SetTxHash(ctx sdk.Context, txHash types.TxHash) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TxHashStoreKey(txHash.BridgeIndex))
	b := k.cdc.MustMarshal(&txHash)
	store.Set(txHash.Index, b)
}

// GetTxHash returns a txHash from its index
func (k Keeper) GetTxHash(
	ctx sdk.Context,
	bridgeIndex []byte,
	index []byte,

) (val types.TxHash, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TxHashStoreKey(bridgeIndex))

	if !k.HasTxHash(ctx, bridgeIndex, index) {
		return types.TxHash{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasTxHash(
	ctx sdk.Context,
	bridgeIndex []byte,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TxHashStoreKey(bridgeIndex))
	return store.Has(index)
}

func (k Keeper) GetAllTxHash(ctx sdk.Context) (list []types.TxHash) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TxHashStoreKey(nil))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TxHash
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
