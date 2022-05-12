package keeper

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/LimeChain/mantrachain/x/mns/utils"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDomainName set a specific domainName in the store from its index
func (k Keeper) SetDomainName(ctx sdk.Context, domainName types.DomainName) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainNameKeyPrefix))
	b := k.cdc.MustMarshal(&domainName)
	store.Set(types.DomainNameKey(
		utils.GetDomainNameIndex(domainName.Domain, domainName.DomainName),
	), b)
}

// HasDomainName checks if the domain name exists in the store
func (k Keeper) HasDomainName(ctx sdk.Context, domain string, domainName string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainNameKeyPrefix))
	return store.Has(types.DomainNameKey(
		utils.GetDomainNameIndex(domain, domainName),
	))
}

// GetDomainName returns a domainName from its index
func (k Keeper) GetDomainName(
	ctx sdk.Context,
	domain string,
	domainName string,

) (val types.DomainName, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DomainNameKeyPrefix))

	if !k.HasDomainName(ctx, domain, domainName) {
		return types.DomainName{}, false
	}

	index := utils.GetDomainNameIndex(domain, domainName)

	b := store.Get(types.DomainNameKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
