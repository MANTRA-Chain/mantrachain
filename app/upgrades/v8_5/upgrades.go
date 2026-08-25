package v8_5

import (
	"context"
	"fmt"

	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ratelimittypes "github.com/cosmos/ibc-go/v11/modules/apps/rate-limiting/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
	storekeys map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting v8.5.0 upgrade...")

		// Virtual fee collection makes evm's DeductFees panic if the fee denom's
		// bank metadata carries the display denom at an exponent other than 18;
		// fail the upgrade here instead of panicking on the first EVM tx.
		evmDenom := evmtypes.GetEVMCoinDenom()
		displayDenom := evmtypes.GetEVMCoinDisplayDenom()
		if metadata, found := keepers.BankKeeper.GetDenomMetaData(ctx, evmDenom); found {
			for _, unit := range metadata.DenomUnits {
				if unit.Denom == displayDenom && unit.Exponent != evmtypes.EighteenDecimals.Uint32() {
					return vm, fmt.Errorf(
						"bank metadata for %s has denom unit %s with exponent %d, virtual fee collection requires 18",
						evmDenom, unit.Denom, unit.Exponent,
					)
				}
			}
		}

		// The replaced rate-limiting fork already reported consensus version 2,
		// so ibc-go's Migrate1to2 (which clears its legacy pending-packet keys)
		// would be skipped; force it to run.
		if _, ok := vm[ratelimittypes.ModuleName]; ok {
			vm[ratelimittypes.ModuleName] = 1
		}

		ctx.Logger().Info("Running module migrations...")
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v8.5.0 complete")
		return vm, nil
	}
}
