package keeper

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDomainName set a specific domainName in the store from its index
func (k Keeper) SetDomainName(ctx sdk.Context, domainName types.DomainName) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainNameKeyPrefix))
	b := k.cdc.MustMarshal(&domainName)
	store.Set(types.DomainNameKey(
		domainName.Index,
	), b)
}

// GetDomainName returns a domainName from its index
func (k Keeper) GetDomainName(
	ctx sdk.Context,
	domain string,
	domainName string,

) (val types.DomainName, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainNameKeyPrefix))

	b := store.Get(types.DomainNameKey(
		domainName + "@" + domain,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDomainName removes a domainName from the store
func (k Keeper) RemoveDomainName(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainNameKeyPrefix))
	store.Delete(types.DomainNameKey(
		index,
	))
}

// GetAllDomainName returns all domainName
func (k Keeper) GetAllDomainName(ctx sdk.Context) (list []types.DomainName) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainNameKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DomainName
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
