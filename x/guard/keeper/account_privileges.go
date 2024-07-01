package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetAccountPrivileges(ctx sdk.Context, account sdk.AccAddress, accountPrivileges []byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.AccountPrivilegesStoreKey())
	store.Set(account, accountPrivileges)
}

func (k Keeper) HasAccountPrivileges(
	ctx sdk.Context,
	account sdk.AccAddress,
) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.AccountPrivilegesStoreKey())
	return store.Has(account)
}

func (k Keeper) GetAccountPrivileges(
	ctx sdk.Context,
	account sdk.AccAddress,
	defaultAccountPrivileges []byte,
) (val []byte, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.AccountPrivilegesStoreKey())

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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.AccountPrivilegesStoreKey())

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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.AccountPrivilegesStoreKey())
	store.Delete(account)
}

func (k Keeper) GetAllAccountPrivileges(ctx sdk.Context) (list []*types.AccountPrivileges) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.AccountPrivilegesStoreKey())
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, &types.AccountPrivileges{
			Account:    iterator.Key(),
			Privileges: iterator.Value(),
		})
	}

	return
}
