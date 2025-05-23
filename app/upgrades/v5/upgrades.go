package v5

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
)

const (
	skipFeemarketPrefixParams = iota + 1
	skipFeemarketPrefixState
	skipFeemarketPrefixEnableHeight = 3
)

var (
	// KeyParams is the store key for the feemarket module's parameters.
	skipKeyParams = []byte{skipFeemarketPrefixParams}

	// KeyState is the store key for the feemarket module's data.
	skipKeyState = []byte{skipFeemarketPrefixState}

	// KeyEnabledHeight is the store key for the feemarket module's enabled height.
	skipKeyEnabledHeight = []byte{skipFeemarketPrefixEnableHeight}
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		// delete skip-mev feemarket
		feemarketstore := runtime.NewKVStoreService(storetypes.NewKVStoreKey(feemarkettypes.StoreKey)).OpenKVStore(c)
		if err := feemarketstore.Delete(skipKeyParams); err != nil {
			return vm, err
		}
		if err := feemarketstore.Delete(skipKeyState); err != nil {
			return vm, err
		}
		if err := feemarketstore.Delete(skipKeyEnabledHeight); err != nil {
			return vm, err
		}

		// initialize new cosmos evm feemarket
		keepers.FeeMarketKeeper.SetParams(ctx, feemarkettypes.DefaultParams())

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v5 complete")
		return vm, nil
	}
}
