package keeper

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDomain set a specific domain in the store from its index
func (k Keeper) SetDomain(ctx sdk.Context, domain types.Domain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainKeyPrefix))
	b := k.cdc.MustMarshal(&domain)
	store.Set(types.DomainKey(
		utils.GetDomainIndex(domain.Domain),
	), b)
}

// GetDomain returns a domain from its domain
func (k Keeper) GetDomain(
	ctx sdk.Context,
	domain string,

) (val types.Domain, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainKeyPrefix))

	if !k.HasDomain(ctx, domain) {
		return types.Domain{}, false
	}

	index := utils.GetDomainIndex(domain)

	b := store.Get(types.DomainKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// HasDomain checks if the domain exists in the store
func (k Keeper) HasDomain(ctx sdk.Context, domain string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainKeyPrefix))
	return store.Has(types.DomainKey(
		utils.GetDomainIndex(domain),
	))
}

// GetAllDomain returns all domain
func (k Keeper) GetAllDomain(ctx sdk.Context) (list []types.Domain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Domain
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
