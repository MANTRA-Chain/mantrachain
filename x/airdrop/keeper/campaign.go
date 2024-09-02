package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetClaimed(ctx sdk.Context, claimed types.Claimed) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ClaimedStoreKey())
	b := k.cdc.MustMarshal(&claimed)
	store.Set(claimed.Index, b)
}

// GetClaimed returns a claimed from its index
func (k Keeper) GetClaimed(
	ctx sdk.Context,
	index []byte,
) (val types.Claimed, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ClaimedStoreKey())

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllClaimed returns all campaign
func (k Keeper) GetAllClaimed(ctx sdk.Context) (list []types.Claimed) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ClaimedStoreKey())
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Claimed
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetLastCampaignId(ctx sdk.Context) (id uint64, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.KeyPrefix(types.LastCampaignIdKey))
	if bz == nil {
		return
	}
	return sdk.BigEndianToUint64(bz), true
}

func (k Keeper) SetLastCampaignId(ctx sdk.Context, id uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.KeyPrefix(types.LastCampaignIdKey), sdk.Uint64ToBigEndian(id))
}

// SetCampaign set a specific campaign in the store from its index
func (k Keeper) SetCampaign(ctx sdk.Context, campaign types.Campaign) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CampaignStoreKey())
	b := k.cdc.MustMarshal(&campaign)
	store.Set(campaign.Index, b)
}

// GetCampaign returns a campaign from its index
func (k Keeper) GetCampaign(
	ctx sdk.Context,
	index []byte,
) (val types.Campaign, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CampaignStoreKey())

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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CampaignStoreKey())
	store.Delete(index)
}

// GetAllCampaign returns all campaign
func (k Keeper) GetAllCampaign(ctx sdk.Context) (list []types.Campaign) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CampaignStoreKey())
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Campaign
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) IterateAllCampaigns(ctx sdk.Context, cb func(plan types.Campaign) (stop bool)) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CampaignStoreKey())
	iter := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var campaign types.Campaign
		k.cdc.MustUnmarshal(iter.Value(), &campaign)
		if cb(campaign) {
			break
		}
	}
}
