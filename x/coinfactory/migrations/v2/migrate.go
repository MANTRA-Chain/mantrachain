package v2

import (
	"cosmossdk.io/core/store"

	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/exported"
	"github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func MigrateStore(
	ctx sdk.Context,
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	legacySubspace exported.Subspace,
	guardKeeper types.GuardKeeper,
) error {
	guardKeeper.AddTransferAccAddressesWhitelist(ctx, []sdk.AccAddress{
		address.Module(types.ModuleName),
	})

	store := storeService.OpenKVStore(ctx)
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	return store.Set(types.ParamsKey, bz)
}
