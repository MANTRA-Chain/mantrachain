package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
)

// SetCw20Contract set cw20Contract in the store
func (k Keeper) SetCw20Contract(ctx sdk.Context, cw20Contract types.Cw20Contract) {
	store :=  prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cw20ContractKey))
	b := k.cdc.MustMarshal(&cw20Contract)
	store.Set([]byte{0}, b)
}

// GetCw20Contract returns cw20Contract
func (k Keeper) GetCw20Contract(ctx sdk.Context) (val types.Cw20Contract, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cw20ContractKey))

	b := store.Get([]byte{0})
    if b == nil {
        return val, false
    }

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCw20Contract removes cw20Contract from the store
func (k Keeper) RemoveCw20Contract(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.Cw20ContractKey))
	store.Delete([]byte{0})
}
