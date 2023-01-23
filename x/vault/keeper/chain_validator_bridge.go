package keeper

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetChainValidatorBridge(ctx sdk.Context, chain string, validator string, chainValidatorBridge types.ChainValidatorBridge) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainValidatorBridgeKey(chain))
	b := k.cdc.MustMarshal(&chainValidatorBridge)
	store.Set(types.GetChainValidatorBridgeIndex(validator), b)
}

// RemoveChainValidatorBridge removes a chainValidatorBridge from the store
func (k Keeper) RemoveChainValidatorBridge(
	ctx sdk.Context,
	chain string,
	validator string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainValidatorBridgeKey(chain))
	store.Delete(types.GetChainValidatorBridgeIndex(
		validator,
	))
}

func (k Keeper) GetChainValidatorBridge(ctx sdk.Context, chain string, validator string) (val types.ChainValidatorBridge, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainValidatorBridgeKey(chain))

	if !k.HasChainValidatorBridge(ctx, chain, validator) {
		return val, false
	}

	b := store.Get(types.GetChainValidatorBridgeIndex(validator))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasChainValidatorBridge(ctx sdk.Context, chain string, validator string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainValidatorBridgeKey(chain))
	return store.Has(types.GetChainValidatorBridgeIndex(validator))
}

func (k Keeper) GetAllChainValidatorBridge(ctx sdk.Context, chain *string) (list []types.ChainValidatorBridge) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ChainValidatorBridgeStoreKey(chain))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ChainValidatorBridge
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
