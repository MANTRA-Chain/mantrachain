package v5rc5

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
	storekeys map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		store := runtime.NewKVStoreService(storekeys[erc20types.StoreKey]).OpenKVStore(c)
		const addressLength = 42
		for key, fn := range map[string]func(ctx sdk.Context, address common.Address){
			"DynamicPrecompiles": keepers.Erc20Keeper.SetDynamicPrecompile,
			"NativePrecompiles":  keepers.Erc20Keeper.SetNativePrecompile,
		} {
			bz, err := store.Get([]byte(key))
			if err != nil {
				return vm, err
			}
			for i := 0; i < len(bz); i += addressLength {
				address := common.HexToAddress(string(bz[i : i+addressLength]))
				ctx.Logger().Info("Set"+key, "address", address.String())
				fn(ctx, address)
			}
		}
		ctx.Logger().Info("Upgrade v5.0.0-rc5 complete")
		return vm, nil
	}
}
