package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetClaimed(ctx sdk.Context, claimed types.Claimed) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClaimedStoreKey())
	b := k.cdc.MustMarshal(&claimed)
	store.Set(claimed.Index, b)
}

// GetClaimed returns a claimed from its index
func (k Keeper) GetClaimed(
	ctx sdk.Context,
	index []byte,

) (val types.Claimed, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClaimedStoreKey())

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllClaimed returns all campaign
func (k Keeper) GetAllClaimed(ctx sdk.Context) (list []types.Claimed) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ClaimedStoreKey())
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Claimed
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetLastCampaignId(ctx sdk.Context) (id uint64, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefix(types.LastCampaignIdKey))
	if bz == nil {
		return
	}
	return sdk.BigEndianToUint64(bz), true
}

func (k Keeper) SetLastCampaignId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.LastCampaignIdKey), sdk.Uint64ToBigEndian(id))
}

// SetCampaign set a specific campaign in the store from its index
func (k Keeper) SetCampaign(ctx sdk.Context, campaign types.Campaign) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignStoreKey())
	b := k.cdc.MustMarshal(&campaign)
	store.Set(campaign.Index, b)
}

// GetCampaign returns a campaign from its index
func (k Keeper) GetCampaign(
	ctx sdk.Context,
	index []byte,

) (val types.Campaign, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignStoreKey())

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCampaign removes a campaign from the store
func (k Keeper) RemoveCampaign(
	ctx sdk.Context,
	index []byte,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignStoreKey())
	store.Delete(index)
}

// GetAllCampaign returns all campaign
func (k Keeper) GetAllCampaign(ctx sdk.Context) (list []types.Campaign) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignStoreKey())
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Campaign
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) IterateAllCampaigns(ctx sdk.Context, cb func(plan types.Campaign) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignStoreKey())
	iter := sdk.KVStorePrefixIterator(store, []byte{})

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var campaign types.Campaign
		k.cdc.MustUnmarshal(iter.Value(), &campaign)
		if cb(campaign) {
			break
		}
	}
}
