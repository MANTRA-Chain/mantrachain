package keeper

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
	utils "github.com/LimeChain/mantrachain/x/mns/utils"
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

	domainNameIndex := utils.GetDomainNameIndex(domain, domainName)

	b := store.Get(types.DomainNameKey(domainNameIndex))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
