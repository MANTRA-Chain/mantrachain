package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"mantrachain/x/guard/types"
)

func (k Keeper) SetAccountPrivileges(ctx sdk.Context, account sdk.AccAddress, accountPrivileges []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountPrivilegesStoreKey())
	store.Set(account, accountPrivileges)
}

func (k Keeper) HasAccountPrivileges(
	ctx sdk.Context,
	account sdk.AccAddress,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountPrivilegesStoreKey())
	return store.Has(account)
}

func (k Keeper) GetAccountPrivileges(
	ctx sdk.Context,
	account sdk.AccAddress,
	defaultAccountPrivileges []byte,
) (val []byte, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountPrivilegesStoreKey())

	if !k.HasAccountPrivileges(ctx, account) {
		if defaultAccountPrivileges != nil {
			return defaultAccountPrivileges, true
		} else {
			return []byte{}, false
		}
	}

	b := store.Get(account)

	if b == nil {
		if defaultAccountPrivileges != nil {
			return defaultAccountPrivileges, true
		} else {
			return val, false
		}
	}

	return b, true
}

func (k Keeper) GetAccountPrivilegesMany(
	ctx sdk.Context,
	accounts []sdk.AccAddress,
	defaultAccountPrivileges []byte,
) (list [][]byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountPrivilegesStoreKey())

	for _, acc := range accounts {
		bz := store.Get(acc)

		if bz == nil && defaultAccountPrivileges != nil {
			bz = defaultAccountPrivileges
		}

		list = append(list, bz)
	}

	return
}

func (k Keeper) RemoveAccountPrivileges(
	ctx sdk.Context,
	account []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountPrivilegesStoreKey())
	store.Delete(account)
}

func (k Keeper) GetAllAccountPrivileges(ctx sdk.Context) (list []*types.AccountPrivileges) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AccountPrivilegesStoreKey())
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, &types.AccountPrivileges{
			Account:    iterator.Key(),
			Privileges: iterator.Value(),
		})
	}

	return
}
