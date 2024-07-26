package v2

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"

	"github.com/MANTRA-Finance/mantrachain/x/guard/exported"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func MigrateStore(
	ctx sdk.Context,
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	legacySubspace exported.Subspace,
	whitelistedTransferAccAddrs []string,
) error {
	store := storeService.OpenKVStore(ctx)
	storeAdapter := runtime.KVStoreAdapter(store)
	whitelistTransfersAccAddrsStore := prefix.NewStore(storeAdapter, types.WhitelistTransfersAccAddrsStoreKey())

	for _, moduleName := range whitelistedTransferAccAddrs {
		account := authtypes.NewModuleAddress(moduleName)
		whitelistTransfersAccAddrs := types.WhitelistTransfersAccAddrs{
			Index:         types.GetWhitelistTransfersAccAddrsIndex(account),
			Account:       account,
			IsWhitelisted: true,
		}
		b := cdc.MustMarshal(&whitelistTransfersAccAddrs)
		whitelistTransfersAccAddrsStore.Set(whitelistTransfersAccAddrs.Index, b)
	}

	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	return store.Set(types.ParamsKey, bz)
}
