package keeper

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDomainName set a specific domainName in the store from its index
func (k Keeper) SetDomainName(ctx sdk.Context, domainName types.DomainName) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DomainNameStoreKey(&domainName.Domain))
	b := k.cdc.MustMarshal(&domainName)
	store.Set(types.GetDomainNameIndex(domainName.Domain, domainName.DomainName), b)
}

// HasDomainName checks if the domain name exists in the store
func (k Keeper) HasDomainName(ctx sdk.Context, domain string, domainName string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DomainNameStoreKey(&domain))
	return store.Has(types.GetDomainNameIndex(domain, domainName))
}

// GetDomainName returns a domainName from its index
func (k Keeper) GetDomainName(
	ctx sdk.Context,
	domain string,
	domainName string,

) (val types.DomainName, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DomainNameStoreKey(&domain))

	if !k.HasDomainName(ctx, domain, domainName) {
		return types.DomainName{}, false
	}

	index := types.GetDomainNameIndex(domain, domainName)

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllDomainName returns all domain
func (k Keeper) GetAllDomainName(ctx sdk.Context, domain *string) (list []types.DomainName) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DomainNameStoreKey(domain))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DomainName
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
