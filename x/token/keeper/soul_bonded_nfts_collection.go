package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/LimeChain/mantrachain/x/token/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
)

// SetSoulBondedNftsCollection set a specific soulBondedNftsCollection in the store from its index
func (k Keeper) SetSoulBondedNftsCollection(ctx sdk.Context, soulBondedNftsCollection types.SoulBondedNftsCollection) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	b := k.cdc.MustMarshal(&soulBondedNftsCollection)
	store.Set(types.SoulBondedNftsCollectionKey(
        soulBondedNftsCollection.Index,
    ), b)
}

// GetSoulBondedNftsCollection returns a soulBondedNftsCollection from its index
func (k Keeper) GetSoulBondedNftsCollection(
    ctx sdk.Context,
    index string,
    
) (val types.SoulBondedNftsCollection, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))

	b := store.Get(types.SoulBondedNftsCollectionKey(
        index,
    ))
    if b == nil {
        return val, false
    }

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSoulBondedNftsCollection removes a soulBondedNftsCollection from the store
func (k Keeper) RemoveSoulBondedNftsCollection(
    ctx sdk.Context,
    index string,
    
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	store.Delete(types.SoulBondedNftsCollectionKey(
	    index,
    ))
}

// GetAllSoulBondedNftsCollection returns all soulBondedNftsCollection
func (k Keeper) GetAllSoulBondedNftsCollection(ctx sdk.Context) (list []types.SoulBondedNftsCollection) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SoulBondedNftsCollectionKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SoulBondedNftsCollection
		k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
	}

    return
}
