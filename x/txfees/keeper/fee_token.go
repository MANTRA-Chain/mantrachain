package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/txfees/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetFeeToken set a specific feeToken in the store from its denom
func (k Keeper) SetFeeToken(ctx sdk.Context, feeToken types.FeeToken) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeeTokenKeyPrefix))
	b := k.cdc.MustMarshal(&feeToken)
	store.Set(types.FeeTokenKey(
		feeToken.Denom,
	), b)
}

// GetFeeToken returns a feeToken from its denom
func (k Keeper) GetFeeToken(
	ctx sdk.Context,
	denom string,

) (val types.FeeToken, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeeTokenKeyPrefix))

	b := store.Get(types.FeeTokenKey(
		denom,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveFeeToken removes a feeToken from the store
func (k Keeper) RemoveFeeToken(
	ctx sdk.Context,
	denom string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeeTokenKeyPrefix))
	store.Delete(types.FeeTokenKey(
		denom,
	))
}

// GetAllFeeToken returns all feeToken
func (k Keeper) GetAllFeeToken(ctx sdk.Context) (list []types.FeeToken) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeeTokenKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.FeeToken
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
