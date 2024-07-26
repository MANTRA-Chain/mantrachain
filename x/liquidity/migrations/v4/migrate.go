package v4

import (
	"cosmossdk.io/core/store"

	"github.com/MANTRA-Finance/mantrachain/x/liquidity/exported"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(
	ctx sdk.Context,
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	legacySubspace exported.Subspace,
	guardKeeper types.GuardKeeper,
) error {
	guardKeeper.AddTransferAccAddressesWhitelist(ctx, []sdk.AccAddress{
		types.DefaultFeeCollectorAddress,
		types.DefaultDustCollectorAddress,
		types.GlobalEscrowAddress,
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
