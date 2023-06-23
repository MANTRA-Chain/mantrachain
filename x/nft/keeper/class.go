package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SaveClass defines a method for creating a new nft class
func (k Keeper) SaveClass(ctx sdk.Context, class types.Class) error {
	if k.HasClass(ctx, class.Id) {
		return sdkerrors.Wrap(types.ErrClassExists, class.Id)
	}
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.Class failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(classStoreKey(class.Id), bz)
	return nil
}

// UpdateClass defines a method for updating an exist nft class
func (k Keeper) UpdateClass(ctx sdk.Context, class types.Class) error {
	if !k.HasClass(ctx, class.Id) {
		return sdkerrors.Wrap(types.ErrClassNotExists, class.Id)
	}
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.Class failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(classStoreKey(class.Id), bz)
	return nil
}

// GetClass defines a method for returning the class information of the specified id
func (k Keeper) GetClass(ctx sdk.Context, classID string) (types.Class, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(classStoreKey(classID))

	var class types.Class
	if len(bz) == 0 {
		return class, false
	}
	k.cdc.MustUnmarshal(bz, &class)
	return class, true
}

// GetClasses defines a method for returning all classes information
func (k Keeper) GetClasses(ctx sdk.Context) (classes []*types.Class) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ClassKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var class types.Class
		k.cdc.MustUnmarshal(iterator.Value(), &class)
		classes = append(classes, &class)
	}
	return
}

// HasClass determines whether the specified classID exist
func (k Keeper) HasClass(ctx sdk.Context, classID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(classStoreKey(classID))
}