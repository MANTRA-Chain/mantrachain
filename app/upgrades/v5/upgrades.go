package v5

import (
	"context"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
)

// skip-mev feemarket store prefixes to delete from feemarket store
// const (
// 	skipFeeMarketPrefixParams = iota + 1
// 	skipFeeMarketPrefixState
// 	skipFeeMarketPrefixEnableHeight = 3
// )

// var (
// 	skipKeyParams        = []byte{skipFeeMarketPrefixParams}
// 	skipKeyState         = []byte{skipFeeMarketPrefixState}
// 	skipKeyEnabledHeight = []byte{skipFeeMarketPrefixEnableHeight}
// )

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		// TODO: Do we need to delete the old feemarket state first? Code below is not working
		// delete skip-mev feemarket state
		// feemarketstore := runtime.NewKVStoreService(storetypes.NewKVStoreKey(feemarkettypes.StoreKey)).OpenKVStore(c)
		// if err := feemarketstore.Delete(skipKeyParams); err != nil {
		// 	return vm, err
		// }
		// if err := feemarketstore.Delete(skipKeyState); err != nil {
		// 	return vm, err
		// }
		// if err := feemarketstore.Delete(skipKeyEnabledHeight); err != nil {
		// 	return vm, err
		// }

		// delete feemarket from the version map so it can be reinitialized to evm feemarket
		delete(vm, feemarkettypes.StoreKey)

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Reset the feemarket module with new params
		feemarketKeeper := keepers.FeeMarketKeeper
		params := feemarkettypes.DefaultParams()
		params.MinGasPrice = sdkmath.LegacyMustNewDecFromStr("0.01")
		if err := feemarketKeeper.SetParams(ctx, params); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v5 complete")
		return vm, nil
	}
}
