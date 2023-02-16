package keeper

import (
	"github.com/LimeChain/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAccPerm set a specific accPerm in the store from its index
func (k Keeper) SetAccPerm(ctx sdk.Context, accPerm types.AccPerm) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AccPermKeyPrefix))
	b := k.cdc.MustMarshal(&accPerm)
	store.Set(types.AccPermKey(
		accPerm.Id,
	), b)
}

// GetAccPerm returns a accPerm from its index
func (k Keeper) GetAccPerm(
	ctx sdk.Context,
	id string,

) (val types.AccPerm, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AccPermKeyPrefix))

	b := store.Get(types.AccPermKey(
		id,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAccPerm removes a accPerm from the store
func (k Keeper) RemoveAccPerm(
	ctx sdk.Context,
	id string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AccPermKeyPrefix))
	store.Delete(types.AccPermKey(
		id,
	))
}

// GetAllAccPerm returns all accPerm
func (k Keeper) GetAllAccPerm(ctx sdk.Context) (list []types.AccPerm) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AccPermKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AccPerm
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
